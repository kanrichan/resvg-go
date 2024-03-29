package resvg

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"io/fs"
	"os"
	"sync/atomic"

	"github.com/kanrichan/resvg-go/internal"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Worker resvg wasm worker
// `Worker` are not goroutine-safe!
type Worker struct {
	ctx  context.Context
	r    wazero.Runtime
	mod  api.Module
	used *atomic.Bool
}

// NewDefaultWorker initialize a resvg wasm worker by default
// `Worker` are not goroutine-safe!
func NewDefaultWorker(ctx context.Context) (*Worker, error) {
	return NewWorker(ctx, wazero.NewRuntimeConfig())
}

// NewWorker initialize a resvg wasm worker with wazero.RuntimeConfig
// `Worker` are not goroutine-safe!
func NewWorker(ctx context.Context, config wazero.RuntimeConfig) (*Worker, error) {
	wasmgzr, err := gzip.NewReader(bytes.NewReader(internal.WasmGZ))
	if err != nil {
		return nil, err
	}
	defer wasmgzr.Close()
	wasm, err := io.ReadAll(wasmgzr)
	if err != nil {
		return nil, err
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
	return &Worker{ctx, r, mod, &atomic.Bool{}}, nil
}

// Close cloes the `Worker`
func (wk *Worker) Close() error {
	return wk.r.Close(wk.ctx)
}

// vfs wasm mount directory
type vfs struct{}

// Open opens file in the wasm mount directory
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
