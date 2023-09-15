package resvg

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"runtime"
	"strings"

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

var wasm []byte

var (
	errWasmFunctionNotFound = errors.New("wasm function not found")
	errWasmReturnInvaild    = errors.New("wasm return invalid")
	errWasmMemoryOutOfRange = errors.New("wasm memory out of range")
)

const (
	fnWasmBytesMalloc             = "__wasm_bytes_malloc"
	fnWasmBytesFree               = "__wasm_bytes_free"
	fnRendererNew                 = "__renderer_new"
	fnRendererDelete              = "__renderer_delete"
	fnRendererRender              = "__renderer_render"
	fnRendererFontdbLoadFontData  = "__renderer_fontdb_load_font_data"
	fnRendererFontdbLoadFontFile  = "__renderer_fontdb_load_font_file"
	fnRendererFontdbLoadFontDir   = "__renderer_fontdb_load_fonts_dir"
	fnRendererOptionsResourcesDir = "__renderer_options_resources_dir"
	fnRendererOptionsDpi          = "__renderer_options_dpi"
	fnRendererOptionsFontFamily   = "__renderer_options_font_family"
	fnRendererOptionsFontSize     = "__renderer_options_font_size"
	fnRendererOptionsLanguages    = "__renderer_options_languages"
	fnRendererOptionsDefaultSize  = "__renderer_options_default_size"
)

type Worker struct {
	context.Context
	r   wazero.Runtime
	mod api.Module
}

type Renderer struct {
	worker *Worker
	ptr    int32
}

func NewDefaultWorker(ctx context.Context) (*Worker, error) {
	return NewWorker(ctx, wazero.NewRuntimeConfig())
}

func NewWorker(ctx context.Context, config wazero.RuntimeConfig) (*Worker, error) {
	if wasm == nil {
		wasmgzr, err := gzip.NewReader(bytes.NewReader(wasmgz))
		if err != nil {
			return nil, err
		}
		defer wasmgzr.Close()
		b, err := io.ReadAll(wasmgzr)
		if err != nil {
			return nil, err
		}
		wasm = b
	}

	r := wazero.NewRuntimeWithConfig(ctx, config)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	moduleConfig := wazero.NewModuleConfig().
		WithStdout(os.Stdout).WithStderr(os.Stderr).
		WithFS(vfs{})

	mod, err := r.InstantiateWithConfig(ctx, wasm, moduleConfig)
	if err != nil {
		return nil, err
	}
	return &Worker{ctx, r, mod}, nil
}

func (ctx *Worker) Close() error {
	return ctx.r.Close(ctx)
}

func (ctx *Worker) NewRenderer() (*Renderer, error) {
	fn := ctx.mod.ExportedFunction(fnRendererNew)
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
	fn := r.worker.mod.ExportedFunction(fnRendererDelete)
	if fn == nil {
		return errWasmFunctionNotFound
	}
	_, err := fn.Call(r.worker)
	if err != nil {
		return err
	}
	return nil
}

func (r *Renderer) GetWorker() *Worker {
	return r.worker
}

func (r *Renderer) Render(svg []byte) ([]byte, error) {
	return r.RenderWithSize(svg, 0, 0)
}

func (r *Renderer) RenderWithSize(svg []byte, width, height uint32) ([]byte, error) {
	fn := r.worker.mod.ExportedFunction(fnRendererRender)
	if fn == nil {
		return nil, errWasmFunctionNotFound
	}
	ptr, err := r.worker.malloc(len(svg))
	if err != nil {
		return nil, err
	}
	if f := r.worker.mod.Memory().Write(ptr, svg); !f {
		return nil, errWasmMemoryOutOfRange
	}
	ret, err := fn.Call(
		r.worker,
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
	defer r.worker.free(retptr, int(retlen))
	b, f := r.worker.mod.Memory().Read(retptr, retlen)
	if !f {
		return nil, errWasmReturnInvaild
	}
	var data = make([]byte, int(retlen), int(retlen))
	copy(data, b)
	return data, nil
}

func (r *Renderer) LoadSystemFonts() (err error) {
	switch runtime.GOOS {
	case "windows":
		err = r.LoadFontDir("C:\\Windows\\Fonts")
		if err != nil {
			return
		}
		home := os.Getenv("USERPROFILE")
		if home == "" {
			return nil
		}
		err = r.LoadFontDir(path.Join(home, "AppData\\Local\\Microsoft\\Windows\\Fonts"))
		if err != nil {
			return
		}
		err = r.LoadFontDir(path.Join(home, "AppData\\Roaming\\Microsoft\\Windows\\Fonts"))
		if err != nil {
			return
		}
	case "darwin":
		err = r.LoadFontDir("/Library/Fonts")
		if err != nil {
			return
		}
		err = r.LoadFontDir("/System/Library/Fonts")
		if err != nil {
			return
		}
		var dir []fs.DirEntry
		dir, err = os.ReadDir("/System/Library/AssetsV2")
		if err != nil {
			return
		}
		for _, entry := range dir {
			if entry.IsDir() && strings.HasPrefix(entry.Name(), "com_apple_MobileAsset_Font") {
				err = r.LoadFontDir(path.Join("/System/Library/AssetsV2", entry.Name()))
				if err != nil {
					return
				}
			}
		}
	case "linux":
		err = r.LoadFontDir("/usr/share/fonts")
		if err != nil {
			return
		}
		err = r.LoadFontDir("/usr/local/share/fonts")
		if err != nil {
			return
		}
		home := os.Getenv("HOME")
		if home == "" {
			return nil
		}
		err = r.LoadFontDir(path.Join(home, ".fonts"))
		if err != nil {
			return
		}
		err = r.LoadFontDir(path.Join(home, ".local/share/fonts"))
		if err != nil {
			return
		}
	}
	return
}

func (r *Renderer) LoadFontData(data []byte) error {
	fn := r.worker.mod.ExportedFunction(fnRendererFontdbLoadFontData)
	if fn == nil {
		return errWasmFunctionNotFound
	}
	ptr, err := r.worker.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.worker.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.worker,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) LoadFontFile(file string) error {
	fn := r.worker.mod.ExportedFunction(fnRendererFontdbLoadFontFile)
	if fn == nil {
		return errWasmFunctionNotFound
	}
	data := []byte(file)
	ptr, err := r.worker.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.worker.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.worker,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) LoadFontDir(dir string) error {
	fn := r.worker.mod.ExportedFunction(fnRendererFontdbLoadFontDir)
	if fn == nil {
		return errWasmFunctionNotFound
	}
	data := []byte(dir)
	ptr, err := r.worker.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.worker.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.worker,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) SetResourcesDir(dir string) error {
	fn := r.worker.mod.ExportedFunction(fnRendererOptionsResourcesDir)
	if fn == nil {
		return errWasmFunctionNotFound
	}
	data := []byte(dir)
	ptr, err := r.worker.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.worker.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.worker,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) SetDpi(dpi float32) error {
	fn := r.worker.mod.ExportedFunction(fnRendererOptionsDpi)
	if fn == nil {
		return errWasmFunctionNotFound
	}
	_, err := fn.Call(
		r.worker,
		api.EncodeI32(r.ptr),
		api.EncodeF32(dpi),
	)
	return err
}

func (r *Renderer) SetFontFamily(family string) error {
	fn := r.worker.mod.ExportedFunction(fnRendererOptionsFontFamily)
	if fn == nil {
		return errWasmFunctionNotFound
	}
	data := []byte(family)
	ptr, err := r.worker.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.worker.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.worker,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) SetFontSize(size float32) error {
	fn := r.worker.mod.ExportedFunction(fnRendererOptionsFontSize)
	if fn == nil {
		return errWasmFunctionNotFound
	}
	_, err := fn.Call(
		r.worker,
		api.EncodeI32(r.ptr),
		api.EncodeF32(size),
	)
	return err
}

func (r *Renderer) SetLanguages(languages string) error {
	fn := r.worker.mod.ExportedFunction(fnRendererOptionsLanguages)
	if fn == nil {
		return errWasmFunctionNotFound
	}
	data := []byte(languages)
	ptr, err := r.worker.malloc(len(data))
	if err != nil {
		return err
	}
	if f := r.worker.mod.Memory().Write(ptr, data); !f {
		return errWasmMemoryOutOfRange
	}
	_, err = fn.Call(
		r.worker,
		api.EncodeI32(r.ptr),
		api.EncodeI32(int32(ptr)),
		api.EncodeI32(int32(len(data))),
	)
	return err
}

func (r *Renderer) SetDefaultSize(width, height float32) error {
	fn := r.worker.mod.ExportedFunction(fnRendererOptionsDefaultSize)
	if fn == nil {
		return errWasmFunctionNotFound
	}
	_, err := fn.Call(
		r.worker,
		api.EncodeI32(r.ptr),
		api.EncodeF32(width),
		api.EncodeF32(height),
	)
	return err
}

func (ctx *Worker) malloc(size int) (uint32, error) {
	fn := ctx.mod.ExportedFunction(fnWasmBytesMalloc)
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

func (ctx *Worker) free(ptr uint32, size int) error {
	fn := ctx.mod.ExportedFunction(fnWasmBytesFree)
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
	return f, nil
}
