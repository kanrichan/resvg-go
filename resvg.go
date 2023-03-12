package resvg

import (
	"archive/zip"
	"bytes"
	"context"
	_ "embed"
	"errors"
	"io"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

//go:generate go run src/gen.go

// wasmzip was compiled using `cargo build --release --target wasm32-unknown-unknown`
// and packed into zip
//
//go:embed target/wasm32-unknown-unknown/release/resvg_go.wasm.zip
var wasmzip []byte
var wasmzr, _ = zip.NewReader(bytes.NewReader(wasmzip), int64(len(wasmzip)))

// instance instance
type instance struct {
	ctx context.Context
	mod api.Module
}

var (
	// ErrOutOfRange wasm memory out of range
	ErrOutOfRange = errors.New("wasm memory out of range")
	// ErrNullWasmPointer null wasm pointer
	ErrNullWasmPointer = errors.New("null wasm pointer")
)

// DefaultResvg DefaultResvg
// instance
func DefaultResvg() (*instance, error) {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	f, err := wasmzr.Open("resvg_go.wasm")
	if err != nil {
		return nil, err
	}
	wasm, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	inst, err := r.Instantiate(ctx, wasm)
	if err != nil {
		return nil, err
	}
	return &instance{ctx, inst}, nil
}

func (inst *instance) ResvgRender(tree *UsvgTree, ft *UsvgFitTo, tf *TinySkiaTransform, pixmap *TinySkiaPixmap) error {
	if tree.free || ft.free || tf.free || pixmap.free {
		return ErrNullWasmPointer
	}
	fn := inst.mod.ExportedFunction("__resvg_render")
	_, err := fn.Call(
		inst.ctx,
		api.EncodeI32(tree.ptr),
		api.EncodeI32(ft.ptr),
		api.EncodeI32(tf.ptr),
		api.EncodeI32(pixmap.ptr),
	)
	if err != nil {
		return err
	}
	ft.free = true
	tf.free = true
	return nil
}
