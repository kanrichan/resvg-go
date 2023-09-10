package resvg

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:generate go run wasm/gen.go

// wasmzip was compiled using `cargo build --release --target wasm32-wasi`
// and packed into gzip
//
//go:embed wasm/resvg.wasm.gz
var wasmgz []byte

var (
	errWasmFunctionNotFound = errors.New("wasm function not found")
	errWasmReturnInvaild    = errors.New("wasm return invalid")
	errWasmMemoryOutOfRange = errors.New("wasm memory out of range")
)

type Context struct {
	context.Context
	r   wazero.Runtime
	mod api.Module
}

type Renderer struct {
	ctx *Context
	ptr int32
}

func NewContext(ctx context.Context) (*Context, error) {
	wasmzr, err := gzip.NewReader(bytes.NewReader(wasmgz))
	if err != nil {
		return nil, err
	}
	defer wasmzr.Close()
	wasm, err := io.ReadAll(wasmzr)
	if err != nil {
		return nil, err
	}

	r := wazero.NewRuntime(ctx)

	config := wazero.NewModuleConfig().
		WithStdout(os.Stdout).WithStderr(os.Stderr).
		WithFS(vfs{})

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.InstantiateWithConfig(ctx, wasm, config)
	if err != nil {
		return nil, err
	}
	return &Context{ctx, r, mod}, nil
}

func (ctx *Context) Close() error {
	return ctx.r.Close(ctx)
}

func (ctx *Context) NewRenderer() (*Renderer, error) {
	fn := ctx.mod.
		ExportedFunction("__renderer_new")
	if fn == nil {
		return nil, errWasmFunctionNotFound
	}
	ret, err := fn.Call(ctx)
	if err != nil {
		return nil, err
	}
	if len(ret) != 1 {
		return nil, errWasmReturnInvaild
	}
	return &Renderer{ctx, api.DecodeI32(ret[0])}, nil
}

func (r *Renderer) Close() error {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_delete")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	_, err := fn.Call(r.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Renderer) Render(svg []byte) ([]byte, error) {
	return r.RenderWithSize(svg, 0, 0)
}

func (r *Renderer) RenderWithSize(svg []byte, width, height uint32) ([]byte, error) {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_render")
	if fn == nil {
		return nil, errWasmFunctionNotFound
	}
	ptr, err := r.ctx.malloc(len(svg))
	if err != nil {
		return nil, err
	}
	if f := r.ctx.mod.Memory().Write(ptr, svg); !f {
		return nil, errWasmMemoryOutOfRange
	}
	ret, err := fn.Call(
		r.ctx,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(svg))),
		api.EncodeU32(width),
		api.EncodeU32(height),
	)
	if err != nil {
		return nil, err
	}
	retptr := uint32(ret[0] >> 32)
	retlen := uint32(ret[0])
	defer r.ctx.free(retptr, int(retlen))
	b, f := r.ctx.mod.Memory().Read(retptr, retlen)
	if !f {
		return nil, errWasmReturnInvaild
	}
	var data = make([]byte, int(retlen), int(retlen))
	copy(data, b)
	return data, nil
}

func (r *Renderer) LoadFontData(data []byte) error {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_fontdb_load_font_data")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	ptr, err := r.ctx.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.ctx.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.ctx,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) LoadFontFile(file string) error {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_fontdb_load_font_file")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	data := []byte(file)
	ptr, err := r.ctx.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.ctx.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.ctx,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) LoadFontDir(dir string) error {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_fontdb_load_fonts_dir")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	data := []byte(dir)
	ptr, err := r.ctx.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.ctx.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.ctx,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) SetResourcesDir(dir string) error {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_options_resources_dir")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	data := []byte(dir)
	ptr, err := r.ctx.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.ctx.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.ctx,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) SetDpi(dpi float32) error {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_options_dpi")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	_, err := fn.Call(
		r.ctx,
		api.EncodeI32(r.ptr),
		api.EncodeF32(dpi),
	)
	return err
}

func (r *Renderer) SetFontFamily(family string) error {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_options_font_family")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	data := []byte(family)
	ptr, err := r.ctx.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.ctx.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.ctx,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) SetFontSize(size float32) error {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_options_font_size")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	_, err := fn.Call(
		r.ctx,
		api.EncodeI32(r.ptr),
		api.EncodeF32(size),
	)
	return err
}

func (r *Renderer) SetLanguages(languages string) error {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_options_languages")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	data := []byte(languages)
	ptr, err := r.ctx.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.ctx.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.ctx,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) SetKeepNamedGroups(iskeep bool) error {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_options_keep_named_groups")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	var flag int32 = 0
	if iskeep {
		flag = 1
	}
	_, err := fn.Call(
		r.ctx,
		api.EncodeI32(r.ptr),
		api.EncodeI32(flag),
	)
	return err
}

func (r *Renderer) SetDefaultSize(width, height float32) error {
	fn := r.ctx.mod.
		ExportedFunction("__renderer_options_default_size")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	_, err := fn.Call(
		r.ctx,
		api.EncodeI32(r.ptr),
		api.EncodeF32(width),
		api.EncodeF32(height),
	)
	return err
}

func (ctx *Context) malloc(size int) (uint32, error) {
	fn := ctx.mod.
		ExportedFunction("__wasm_bytes_malloc")
	if fn == nil {
		return 0, errWasmFunctionNotFound
	}
	ret, err := fn.Call(ctx, api.EncodeI32(int32(size)))
	if err != nil {
		return 0, err
	}
	if len(ret) != 1 || ret[0] == 0 {
		return 0, errWasmReturnInvaild
	}
	return uint32(api.DecodeI32(ret[0])), nil
}

func (ctx *Context) free(ptr uint32, size int) error {
	fn := ctx.mod.
		ExportedFunction("__wasm_bytes_free")
	if fn == nil {
		return errWasmFunctionNotFound
	}
	_, err := fn.Call(
		ctx,
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(size)),
	)
	if err != nil {
		return err
	}
	return nil
}

type vfs struct{}

func (vfs vfs) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err // nil fs.File
	}
	fmt.Println("open", name, f)
	return f, nil
}
