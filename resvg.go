package main

import (
	"archive/zip"
	"bytes"
	"context"
	_ "embed"
	"encoding/binary"
	"errors"
	"io"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

//go:generate src/gen.go

// wasmzip was compiled using `cargo build --release --target wasm32-unknown-unknown`
// and packed into zip
//
//go:embed target/wasm32-unknown-unknown/release/resvg_go.wasm.zip
var wasmzip []byte
var wasmzr, _ = zip.NewReader(bytes.NewReader(wasmzip), int64(len(wasmzip)))

var (
	// ErrOutOfRange wasm memory out of range
	ErrOutOfRange = errors.New("wasm memory out of range")
)

// DefaultRenderer is the global genderer without mutex
var DefaultRenderer, _ = NewResvgo(nil)

// Resvgo is the resvg_go.wasm runtime
type Resvgo struct {
	ctx  context.Context
	wasi api.Module
}

// NewResvgo init a new resvg_go.wasm runtime
func NewResvgo(customwasm []byte) (rs Resvgo, err error) {
	rs.ctx = context.Background()
	r := wazero.NewRuntime(rs.ctx)
	if customwasm != nil {
		rs.wasi, err = r.Instantiate(rs.ctx, customwasm)
		return
	}
	f, err := wasmzr.Open("resvg_go.wasm")
	if err != nil {
		return
	}
	wasm, err := io.ReadAll(f)
	if err != nil {
		return
	}
	rs.wasi, err = r.Instantiate(rs.ctx, wasm)
	return
}

// RustBytes RustBytes
type RustBytes struct {
	ptr int32
	len int32

	rs *Resvgo
}

// NewRustBytes RustBytes
func (rs *Resvgo) NewRustBytes(data []byte) (o RustBytes, err error) {
	fn := rs.wasi.ExportedFunction("__floattech_bytes_allocate")
	r, err := fn.Call(rs.ctx, api.EncodeI32(int32(len(data))))
	if err != nil {
		return
	}
	if !rs.wasi.Memory().Write(api.DecodeU32(r[0]), []byte(data)) {
		err = ErrOutOfRange
		return
	}
	return RustBytes{ptr: api.DecodeI32(r[0]), len: int32(len(data)), rs: rs}, nil
}

// Free Free
func (o RustBytes) Free() error {
	fn := o.rs.wasi.ExportedFunction("__floattech_bytes_free")
	_, err := fn.Call(o.rs.ctx, uint64(o.ptr), uint64(o.len))
	if err != nil {
		return err
	}
	return nil
}

// RustBytesRet RustBytesRet
type RustBytesRet struct {
	ptr int32

	rs *Resvgo
}

// NewRustBytesRet NewRustBytesRet
func (rs *Resvgo) NewRustBytesRet() (o RustBytesRet, err error) {
	fn := rs.wasi.ExportedFunction("__floattech_bytes_allocate")
	r, err := fn.Call(rs.ctx, 8)
	if err != nil {
		return
	}
	return RustBytesRet{ptr: api.DecodeI32(r[0]), rs: rs}, nil
}

// Read Read
func (o RustBytesRet) Read() ([]byte, error) {
	ob, f := o.rs.wasi.Memory().Read(uint32(o.ptr), 4)
	if !f {
		return nil, ErrOutOfRange
	}
	offset := binary.LittleEndian.Uint32(ob)
	lb, f := o.rs.wasi.Memory().Read(uint32(o.ptr)+4, 4)
	if !f {
		return nil, ErrOutOfRange
	}
	length := binary.LittleEndian.Uint32(lb)
	data, f := o.rs.wasi.Memory().Read(offset, length)
	if !f {
		return nil, ErrOutOfRange
	}
	return data, nil
}

// Free Free
func (o RustBytesRet) Free() error {
	ob, f := o.rs.wasi.Memory().Read(uint32(o.ptr), 4)
	if !f {
		return ErrOutOfRange
	}
	offset := binary.LittleEndian.Uint32(ob)
	lb, f := o.rs.wasi.Memory().Read(uint32(o.ptr)+4, 4)
	if !f {
		return ErrOutOfRange
	}
	length := binary.LittleEndian.Uint32(lb)
	fn := o.rs.wasi.ExportedFunction("__floattech_bytes_free")
	_, err := fn.Call(o.rs.ctx, api.EncodeU32(offset), api.EncodeU32(length))
	if err != nil {
		return err
	}
	_, err = fn.Call(o.rs.ctx, api.EncodeI32(o.ptr), 8)
	return err
}

// FontDatabase FontDatabase
type FontDatabase struct {
	ptr int32

	rs *Resvgo
}

// NewFontDatabase NewFontDatabase
func (rs *Resvgo) NewFontDatabase() (o FontDatabase, err error) {
	fn := rs.wasi.ExportedFunction("__floattech_database_new")
	r, err := fn.Call(rs.ctx)
	if err != nil {
		return
	}
	return FontDatabase{ptr: api.DecodeI32(r[0]), rs: rs}, nil
}

// LoadFontData LoadFontData
func (o FontDatabase) LoadFontData(data []byte) error {
	rb, err := o.rs.NewRustBytes(data)
	if err != nil {
		return err
	}
	defer rb.Free()
	fn := o.rs.wasi.ExportedFunction("__floattech_database_load_font_data")
	_, err = fn.Call(o.rs.ctx, api.EncodeI32(o.ptr), api.EncodeI32(rb.ptr), api.EncodeI32(rb.len))
	return err
}

// Len Len
func (o FontDatabase) Len() (int32, error) {
	fn := o.rs.wasi.ExportedFunction("__floattech_database_len")
	r, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeI32(r[0]), nil
}

// Free Free
func (o FontDatabase) Free() error {
	fn := o.rs.wasi.ExportedFunction("__floattech_database_free")
	_, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	return err
}

// Size Size
type Size struct {
	ptr int32

	rs *Resvgo
}

// NewSize NewSize
func (rs *Resvgo) NewSize(width float64, height float64) (o Size, err error) {
	fn := rs.wasi.ExportedFunction("__floattech_size_new")
	r, err := fn.Call(rs.ctx, api.EncodeF64(width), api.EncodeF64(height))
	if err != nil {
		return
	}
	return Size{ptr: api.DecodeI32(r[0]), rs: rs}, nil
}

// Width Width
func (o Size) Width() (float64, error) {
	fn := o.rs.wasi.ExportedFunction("__floattech_size_width")
	r, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeF64(r[0]), nil
}

// Height Height
func (o Size) Height() (float64, error) {
	fn := o.rs.wasi.ExportedFunction("__floattech_size_height")
	r, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeF64(r[0]), nil
}

// ToScreenSize ToScreenSize
func (o Size) ToScreenSize() (p ScreenSize, err error) {
	fn := o.rs.wasi.ExportedFunction("__floattech_size_to_screen_size")
	r, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return
	}
	return ScreenSize{ptr: api.DecodeI32(r[0]), rs: o.rs}, nil
}

// Free Free
func (o Size) Free() error {
	fn := o.rs.wasi.ExportedFunction("__floattech_size_free")
	_, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	return err
}

// ScreenSize ScreenSize
type ScreenSize struct {
	ptr int32

	rs *Resvgo
}

// NewScreenSize NewScreenSize
func (rs *Resvgo) NewScreenSize(width uint32, height uint32) (o ScreenSize, err error) {
	fn := rs.wasi.ExportedFunction("__floattech_screen_size_new")
	r, err := fn.Call(rs.ctx, api.EncodeU32(width), api.EncodeU32(height))
	if err != nil {
		return
	}
	return ScreenSize{ptr: api.DecodeI32(r[0]), rs: rs}, nil
}

// Width Width
func (o ScreenSize) Width() (uint32, error) {
	fn := o.rs.wasi.ExportedFunction("__floattech_screen_size_width")
	r, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeU32(r[0]), nil
}

// Height Height
func (o ScreenSize) Height() (uint32, error) {
	fn := o.rs.wasi.ExportedFunction("__floattech_screen_size_height")
	r, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeU32(r[0]), nil
}

// Free Free
func (o ScreenSize) Free() error {
	fn := o.rs.wasi.ExportedFunction("__floattech_screen_size_free")
	_, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	return err
}

// DefaultOptions DefaultOptions
type DefaultOptions struct {
	ptr int32

	rs *Resvgo
}

// NewDefaultOptions NewDefaultOptions
func (rs *Resvgo) NewDefaultOptions() (o DefaultOptions, err error) {
	fn := rs.wasi.ExportedFunction("__floattech_option_default")
	r, err := fn.Call(rs.ctx)
	if err != nil {
		return
	}
	return DefaultOptions{ptr: api.DecodeI32(r[0]), rs: rs}, nil
}

// SetFontFamily SetFontFamily
func (o DefaultOptions) SetFontFamily(fontFamily string) error {
	rb, err := o.rs.NewRustBytes([]byte(fontFamily))
	if err != nil {
		return err
	}
	defer rb.Free()
	fn := o.rs.wasi.ExportedFunction("__floattech_option_set_font_family")
	_, err = fn.Call(o.rs.ctx, api.EncodeI32(o.ptr), api.EncodeI32(rb.ptr), api.EncodeI32(rb.len))
	return err
}

// SetDefaultSize SetDefaultSize
func (o DefaultOptions) SetDefaultSize(size *Size) error {
	fn := o.rs.wasi.ExportedFunction("__floattech_option_set_default_size")
	_, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr), api.EncodeI32(size.ptr))
	return err
}

// Free Free
func (o DefaultOptions) Free() error {
	fn := o.rs.wasi.ExportedFunction("__floattech_option_free")
	_, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	return err
}

// Tree Tree
type Tree struct {
	ptr int32

	rs *Resvgo
}

// NewTreeFromData NewTreeFromData
func (rs *Resvgo) NewTreeFromData(data []byte, opt *DefaultOptions) (o Tree, err error) {
	rb, err := rs.NewRustBytes(data)
	if err != nil {
		return
	}
	defer rb.Free()
	fn := rs.wasi.ExportedFunction("__floattech_tree_from_data")
	r, err := fn.Call(rs.ctx, api.EncodeI32(rb.ptr), api.EncodeI32(rb.len), api.EncodeI32(opt.ptr))
	if err != nil {
		return
	}
	return Tree{ptr: api.DecodeI32(r[0]), rs: rs}, nil
}

// ConvertText ConvertText
func (o Tree) ConvertText(db *FontDatabase, keepNamedGroups bool) error {
	fn := o.rs.wasi.ExportedFunction("__floattech_tree_convert_text")
	f := int32(0)
	if keepNamedGroups {
		f = 1
	}
	_, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr), api.EncodeI32(db.ptr), api.EncodeI32(f))
	return err
}

// GetSize GetSize
func (o Tree) GetSize() (p Size, err error) {
	fn := o.rs.wasi.ExportedFunction("__floattech_tree_get_size")
	r, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return
	}
	return Size{ptr: api.DecodeI32(r[0]), rs: o.rs}, nil
}

// Free Free
func (o Tree) Free() error {
	fn := o.rs.wasi.ExportedFunction("__floattech_tree_free")
	_, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	return err
}

// Pixmap Pixmap
type Pixmap struct {
	ptr int32

	rs *Resvgo
}

// NewPixmap NewPixmap
func (rs *Resvgo) NewPixmap(width uint32, height uint32) (o Pixmap, err error) {
	fn := rs.wasi.ExportedFunction("__floattech_pixmap_new")
	r, err := fn.Call(rs.ctx, api.EncodeU32(width), api.EncodeU32(height))
	if err != nil {
		return
	}
	return Pixmap{ptr: api.DecodeI32(r[0]), rs: rs}, nil
}

// EncodePNG EncodePNG
func (o Pixmap) EncodePNG() ([]byte, error) {
	rb, err := o.rs.NewRustBytesRet()
	if err != nil {
		return nil, err
	}
	defer rb.Free()
	fn := o.rs.wasi.ExportedFunction("__floattech_pixmap_encode_png")
	_, err = fn.Call(o.rs.ctx, api.EncodeI32(rb.ptr), api.EncodeI32(o.ptr))
	if err != nil {
		return nil, err
	}
	return rb.Read()
}

// Free Free
func (o Pixmap) Free() error {
	fn := o.rs.wasi.ExportedFunction("__floattech_pixmap_free")
	_, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr))
	return err
}

// Render Render
func (o Tree) Render(pixmap Pixmap) error {
	fn := o.rs.wasi.ExportedFunction("__floattech_render")
	_, err := fn.Call(o.rs.ctx, api.EncodeI32(o.ptr), api.EncodeI32(pixmap.ptr))
	return err
}

// Render Render
func (o Pixmap) Render(tree Tree) error {
	fn := o.rs.wasi.ExportedFunction("__floattech_render")
	_, err := fn.Call(o.rs.ctx, api.EncodeI32(tree.ptr), api.EncodeI32(o.ptr))
	return err
}
