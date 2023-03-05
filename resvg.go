package main

import (
	"context"
	_ "embed"
	"encoding/binary"
	"errors"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// greetWasm was compiled using `cargo build --release --target wasm32-unknown-unknown`
//
//go:embed target/wasm32-unknown-unknown/release/resvg_go.wasm
var wasm []byte

var (
	ctx  = context.Background()
	wasi api.Module

	errOutOfRange = errors.New("wasm memory out of range")
)

func init() {
	r := wazero.NewRuntime(ctx)
	var err error
	wasi, err = r.Instantiate(ctx, wasm)
	if err != nil {
		panic(err)
	}
}

// RustBytes RustBytes
type RustBytes struct {
	ptr int32
	len int32
}

// NewRustBytes RustBytes
func NewRustBytes(data []byte) (*RustBytes, error) {
	fn := wasi.ExportedFunction("__floattech_bytes_allocate")
	r, err := fn.Call(ctx, api.EncodeI32(int32(len(data))))
	if err != nil {
		return nil, err
	}
	if !wasi.Memory().Write(api.DecodeU32(r[0]), []byte(data)) {
		return nil, errOutOfRange
	}
	return &RustBytes{api.DecodeI32(r[0]), int32(len(data))}, nil
}

// Free Free
func (o *RustBytes) Free() error {
	fn := wasi.ExportedFunction("__floattech_bytes_free")
	_, err := fn.Call(ctx, uint64(o.ptr), uint64(o.len))
	if err != nil {
		return err
	}
	o = nil
	return nil
}

// RustBytesRet RustBytesRet
type RustBytesRet struct {
	ptr int32
}

// NewRustBytesRet NewRustBytesRet
func NewRustBytesRet() (*RustBytesRet, error) {
	fn := wasi.ExportedFunction("__floattech_bytes_allocate")
	r, err := fn.Call(ctx, 8)
	if err != nil {
		return nil, err
	}
	return &RustBytesRet{api.DecodeI32(r[0])}, nil
}

// Read Read
func (o *RustBytesRet) Read() ([]byte, error) {
	ob, f := wasi.Memory().Read(uint32(o.ptr), 4)
	if !f {
		return nil, errOutOfRange
	}
	offset := binary.LittleEndian.Uint32(ob)
	lb, f := wasi.Memory().Read(uint32(o.ptr)+4, 4)
	if !f {
		return nil, errOutOfRange
	}
	length := binary.LittleEndian.Uint32(lb)
	data, f := wasi.Memory().Read(offset, length)
	if !f {
		return nil, errOutOfRange
	}
	return data, nil
}

// Free Free
func (o *RustBytesRet) Free() error {
	ob, f := wasi.Memory().Read(uint32(o.ptr), 4)
	if !f {
		return errOutOfRange
	}
	offset := binary.LittleEndian.Uint32(ob)
	lb, f := wasi.Memory().Read(uint32(o.ptr)+4, 4)
	if !f {
		return errOutOfRange
	}
	length := binary.LittleEndian.Uint32(lb)
	fn := wasi.ExportedFunction("__floattech_bytes_free")
	_, err := fn.Call(ctx, api.EncodeU32(offset), api.EncodeU32(length))
	if err != nil {
		return err
	}
	_, err = fn.Call(ctx, api.EncodeI32(o.ptr), 8)
	if err != nil {
		return err
	}
	o = nil
	return nil
}

// FontDatabase FontDatabase
type FontDatabase struct {
	ptr int32
}

// NewFontDatabase NewFontDatabase
func NewFontDatabase() (*FontDatabase, error) {
	fn := wasi.ExportedFunction("__floattech_database_new")
	r, err := fn.Call(ctx)
	if err != nil {
		return nil, err
	}
	return &FontDatabase{ptr: api.DecodeI32(r[0])}, nil
}

// LoadFontData LoadFontData
func (o *FontDatabase) LoadFontData(data []byte) error {
	rb, err := NewRustBytes(data)
	if err != nil {
		return err
	}
	defer rb.Free()
	fn := wasi.ExportedFunction("__floattech_database_load_font_data")
	_, err = fn.Call(ctx, api.EncodeI32(o.ptr), api.EncodeI32(rb.ptr), api.EncodeI32(rb.len))
	return err
}

// Len Len
func (o *FontDatabase) Len() (int32, error) {
	fn := wasi.ExportedFunction("__floattech_database_len")
	r, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeI32(r[0]), nil
}

// Free Free
func (o *FontDatabase) Free() error {
	fn := wasi.ExportedFunction("__floattech_database_free")
	_, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o = nil
	return nil
}

// Size Size
type Size struct {
	ptr int32
}

// NewSize NewSize
func NewSize(width float64, height float64) (*Size, error) {
	fn := wasi.ExportedFunction("__floattech_size_new")
	r, err := fn.Call(ctx, api.EncodeF64(width), api.EncodeF64(height))
	if err != nil {
		return nil, err
	}
	return &Size{ptr: api.DecodeI32(r[0])}, nil
}

// Width Width
func (o *Size) Width() (float64, error) {
	fn := wasi.ExportedFunction("__floattech_size_width")
	r, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeF64(r[0]), nil
}

// Height Height
func (o *Size) Height() (float64, error) {
	fn := wasi.ExportedFunction("__floattech_size_height")
	r, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeF64(r[0]), nil
}

// ToScreenSize ToScreenSize
func (o *Size) ToScreenSize() (*ScreenSize, error) {
	fn := wasi.ExportedFunction("__floattech_size_to_screen_size")
	r, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return nil, err
	}
	return &ScreenSize{ptr: api.DecodeI32(r[0])}, nil
}

// Free Free
func (o *Size) Free() error {
	fn := wasi.ExportedFunction("__floattech_size_free")
	_, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o = nil
	return nil
}

// ScreenSize ScreenSize
type ScreenSize struct {
	ptr int32
}

// NewScreenSize NewScreenSize
func NewScreenSize(width uint32, height uint32) (*ScreenSize, error) {
	fn := wasi.ExportedFunction("__floattech_screen_size_new")
	r, err := fn.Call(ctx, api.EncodeU32(width), api.EncodeU32(height))
	if err != nil {
		return nil, err
	}
	return &ScreenSize{ptr: api.DecodeI32(r[0])}, nil
}

// Width Width
func (o *ScreenSize) Width() (uint32, error) {
	fn := wasi.ExportedFunction("__floattech_screen_size_width")
	r, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeU32(r[0]), nil
}

// Height Height
func (o *ScreenSize) Height() (uint32, error) {
	fn := wasi.ExportedFunction("__floattech_screen_size_height")
	r, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeU32(r[0]), nil
}

// Free Free
func (o *ScreenSize) Free() error {
	fn := wasi.ExportedFunction("__floattech_screen_size_free")
	_, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o = nil
	return nil
}

// Options Options
type Options struct {
	ptr int32
}

// DefaultOptions DefaultOptions
func DefaultOptions() (*Options, error) {
	fn := wasi.ExportedFunction("__floattech_option_default")
	r, err := fn.Call(ctx)
	if err != nil {
		return nil, err
	}
	return &Options{ptr: api.DecodeI32(r[0])}, nil
}

// SetFontFamily SetFontFamily
func (o *Options) SetFontFamily(fontFamily string) error {
	rb, err := NewRustBytes([]byte(fontFamily))
	if err != nil {
		return err
	}
	defer rb.Free()
	fn := wasi.ExportedFunction("__floattech_option_set_font_family")
	_, err = fn.Call(ctx, api.EncodeI32(o.ptr), api.EncodeI32(rb.ptr), api.EncodeI32(rb.len))
	if err != nil {
		return err
	}
	return nil
}

// SetDefaultSize SetDefaultSize
func (o *Options) SetDefaultSize(size *Size) error {
	fn := wasi.ExportedFunction("__floattech_option_set_default_size")
	_, err := fn.Call(ctx, api.EncodeI32(o.ptr), api.EncodeI32(size.ptr))
	if err != nil {
		return err
	}
	return nil
}

// Free Free
func (o *Options) Free() error {
	fn := wasi.ExportedFunction("__floattech_option_free")
	_, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o = nil
	return nil
}

// Tree Tree
type Tree struct {
	ptr int32
}

// TreeFromData TreeFromData
func TreeFromData(data []byte, opt *Options) (*Tree, error) {
	rb, err := NewRustBytes(data)
	if err != nil {
		return nil, err
	}
	defer rb.Free()
	fn := wasi.ExportedFunction("__floattech_tree_from_data")
	r, err := fn.Call(ctx, api.EncodeI32(rb.ptr), api.EncodeI32(rb.len), api.EncodeI32(opt.ptr))
	if err != nil {
		return nil, err
	}
	return &Tree{ptr: api.DecodeI32(r[0])}, nil
}

// ConvertText ConvertText
func (o *Tree) ConvertText(db *FontDatabase, keepNamedGroups bool) error {
	fn := wasi.ExportedFunction("__floattech_tree_convert_text")
	var f int32
	if keepNamedGroups {
		f = 1
	}
	_, err := fn.Call(ctx, api.EncodeI32(o.ptr), api.EncodeI32(db.ptr), api.EncodeI32(f))
	if err != nil {
		return err
	}
	return nil
}

// GetSize GetSize
func (o *Tree) GetSize() (*Size, error) {
	fn := wasi.ExportedFunction("__floattech_tree_get_size")
	r, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return nil, err
	}
	return &Size{ptr: api.DecodeI32(r[0])}, nil
}

// Free Free
func (o *Tree) Free() error {
	fn := wasi.ExportedFunction("__floattech_tree_free")
	_, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o = nil
	return nil
}

// Pixmap Pixmap
type Pixmap struct {
	ptr int32
}

// NewPixmap NewPixmap
func NewPixmap(width uint32, height uint32) (*Pixmap, error) {
	fn := wasi.ExportedFunction("__floattech_pixmap_new")
	r, err := fn.Call(ctx, api.EncodeU32(width), api.EncodeU32(height))
	if err != nil {
		return nil, err
	}
	return &Pixmap{ptr: api.DecodeI32(r[0])}, nil
}

// EncodePNG EncodePNG
func (o *Pixmap) EncodePNG() ([]byte, error) {
	rb, err := NewRustBytesRet()
	if err != nil {
		return nil, err
	}
	defer rb.Free()
	fn := wasi.ExportedFunction("__floattech_pixmap_encode_png")
	_, err = fn.Call(ctx, api.EncodeI32(rb.ptr), api.EncodeI32(o.ptr))
	if err != nil {
		return nil, err
	}
	return rb.Read()
}

// Free Free
func (o *Pixmap) Free() error {
	fn := wasi.ExportedFunction("__floattech_pixmap_free")
	_, err := fn.Call(ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o = nil
	return nil
}

// Render Render
func Render(tree *Tree, pixmap *Pixmap) error {
	fn := wasi.ExportedFunction("__floattech_render")
	_, err := fn.Call(ctx, api.EncodeI32(tree.ptr), api.EncodeI32(pixmap.ptr))
	return err
}
