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

var (
	ctx  = context.Background()
	wasi api.Module

	// ErrOutOfRange wasm memory out of range
	ErrOutOfRange = errors.New("wasm memory out of range")
	// ErrNullWasmPointer null wasm pointer
	ErrNullWasmPointer = errors.New("null wasm pointer")
)

func init() {
	r := wazero.NewRuntime(ctx)
	f, err := wasmzr.Open("resvg_go.wasm")
	if err != nil {
		panic(err)
	}
	wasm, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	wasi, err = r.Instantiate(ctx, wasm)
	if err != nil {
		panic(err)
	}
}

var funcTreeRender = wasi.ExportedFunction("__floattech_render")

// Render Render
func Render(tree *Tree, pixmap *Pixmap) error {
	_, err := funcTreeRender.Call(
		ctx,
		api.EncodeI32(int32(*tree)),
		api.EncodeI32(int32(*pixmap)),
	)
	return err
}
