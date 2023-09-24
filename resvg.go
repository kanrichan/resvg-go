package resvg

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"io"
	"io/fs"
	"os"

	"github.com/kanrichan/resvg-go/internal"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:generate go run internal/gen/gen.go

var wasm []byte

type Worker struct {
	context.Context
	r   wazero.Runtime
	mod api.Module
}

func NewDefaultWorker(ctx context.Context) (*Worker, error) {
	return NewWorker(ctx, wazero.NewRuntimeConfig())
}

func NewWorker(ctx context.Context, config wazero.RuntimeConfig) (*Worker, error) {
	if wasm == nil {
		wasmgzr, err := gzip.NewReader(bytes.NewReader(internal.WasmGZ))
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

func (wk *Worker) Close() error {
	return wk.r.Close(wk)
}

func (wk *Worker) Render(svg []byte) ([]byte, error) {
	options, err := internal.UsvgOptionsDefault(wk, wk.mod)
	if err != nil {
		return nil, err
	}
	defer internal.UsvgOptionsDelete(wk, wk.mod, options)
	tree, err := internal.UsvgTreeFromData(wk, wk.mod, svg, options)
	if err != nil {
		return nil, err
	}
	defer internal.UsvgTreeDelete(wk, wk.mod, tree)
	rtree, err := internal.ResvgTreeFromUsvg(wk, wk.mod, tree)
	if err != nil {
		return nil, err
	}
	defer internal.ResvgTreeDelete(wk, wk.mod, rtree)
	transform, err := internal.TinySkiaTransformIdentity(wk, wk.mod)
	if err != nil {
		return nil, err
	}
	defer internal.TinySkiaTransformDelete(wk, wk.mod, transform)
	width, err := internal.UsvgTreeGetWidth(wk, wk.mod, tree)
	if err != nil {
		return nil, err
	}
	height, err := internal.UsvgTreeGetHeight(wk, wk.mod, tree)
	if err != nil {
		return nil, err
	}
	pixmap, err := internal.TinySkiaPixmapNew(wk, wk.mod, uint32(width), uint32(height))
	if err != nil {
		return nil, err
	}
	defer internal.TinySkiaPixmapDelete(wk, wk.mod, pixmap)
	err = internal.ResvgTreeRender(wk, wk.mod, rtree, transform, pixmap)
	if err != nil {
		return nil, err
	}
	data, err := internal.TinySkiaPixmapEncodePng(wk, wk.mod, pixmap)
	if err != nil {
		return nil, err
	}
	return data, err
}

type renderer struct {
	fontdb  int32
	options options
}

type options int32

func NewRenderer() *renderer {
	return &renderer{}
}

func (opt options) SetDpi(dpi float32) {

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
