package resvg

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/kanrichan/resvg-go/internal"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:generate go run internal/gen/gen.go

var wasm []byte

type worker struct {
	ctx context.Context
	r   wazero.Runtime
	mod api.Module
}

func NewDefaultWorker(ctx context.Context) (*worker, error) {
	return NewWorker(ctx, wazero.NewRuntimeConfig())
}

func NewWorker(ctx context.Context, config wazero.RuntimeConfig) (*worker, error) {
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
	return &worker{ctx, r, mod}, nil
}

func (wk *worker) Close() error {
	return wk.r.Close(wk.ctx)
}

type tree struct {
	*worker
	ptr int32
}

type pixmap struct {
	*worker
	ptr int32
}

type fontdb struct {
	*worker
	ptr int32
}

type transform func(context.Context, api.Module) (int32, error)

func TransformIdentity() transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformIdentity(ctx, mod)
	}
}

func TransformFromRow(sx float32, ky float32, kx float32, sy float32, tx float32, ty float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromRow(ctx, mod, sx, ky, kx, sy, tx, ty)
	}
}

func TransformFromTranslate(tx float32, ty float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromTranslate(ctx, mod, tx, ty)
	}
}

func TransformFromScale(width float32, height float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromScale(ctx, mod, width, height)
	}
}

func TransformFromSkew(kx float32, ky float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromSkew(ctx, mod, kx, ky)
	}
}

func TransformFromRotate(angle float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromRotate(ctx, mod, angle)
	}
}

func TransformFromRotateAt(angle float32, tx float32, ty float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromRotateAt(ctx, mod, angle, tx, ty)
	}
}

func (wk *worker) NewFontDB() (*fontdb, error) {
	db, err := internal.FontdbDatabaseDefault(wk.ctx, wk.mod)
	if err != nil {
		return nil, err
	}
	return &fontdb{wk, db}, nil
}

func (db *fontdb) Close() error {
	return internal.FontdbDatabaseDelete(db.ctx, db.mod, db.ptr)
}

func (db *fontdb) LoadFontFile(file string) error {
	return internal.FontdbDatabaseLoadFontFile(db.ctx, db.mod, db.ptr, file)
}

func (db *fontdb) LoadFontsDir(dir string) error {
	return internal.FontdbDatabaseLoadFontsDir(db.ctx, db.mod, db.ptr, dir)
}

func (db *fontdb) LoadFromData(data []byte) error {
	return internal.FontdbDatabaseLoadFontData(db.ctx, db.mod, db.ptr, data)
}

func (db *fontdb) SetSerifFamily(family string) error {
	return internal.FontdbDatabaseSetSerifFamily(db.ctx, db.mod, db.ptr, family)
}

func (db *fontdb) SetSansSerifFamily(family string) error {
	return internal.FontdbDatabaseSetSansSerifFamily(db.ctx, db.mod, db.ptr, family)
}

func (db *fontdb) SetCursiveFamily(family string) error {
	return internal.FontdbDatabaseSetCursiveFamily(db.ctx, db.mod, db.ptr, family)
}

func (db *fontdb) SetFantasyFamily(family string) error {
	return internal.FontdbDatabaseSetFantasyFamily(db.ctx, db.mod, db.ptr, family)
}

func (db *fontdb) SetMonospaceFamily(family string) error {
	return internal.FontdbDatabaseSetMonospaceFamily(db.ctx, db.mod, db.ptr, family)
}

type Options struct {
	ResourcesDir       string
	Dpi                float32
	FontFamily         string
	FontSize           float32
	Languages          []string
	ShapeRenderingMode ShapeRenderingMode
	TextRenderingMode  TextRenderingMode
	ImageRenderingMode ImageRenderingMode
	DefaultSizeWidth   float32
	DefaultSizeHeight  float32
}

type ShapeRenderingMode int32
type TextRenderingMode int32
type ImageRenderingMode int32

const (
	ShapeRenderingModeOptimizeSpeed ShapeRenderingMode = iota
	ShapeRenderingModeCrispEdges
	ShapeRenderingModeGeometricPrecision
)

const (
	TextRenderingModeOptimizeSpeed TextRenderingMode = iota
	TextRenderingModeOptimizeLegibility
	TextRenderingModeGeometricPrecision
)

const (
	ImageRenderingModeOptimizeQuality ImageRenderingMode = iota
	ImageRenderingModeOptimizeSpeed
)

func (wk *worker) NewTree(data []byte, options *Options) (*tree, error) {
	o, err := internal.UsvgOptionsDefault(wk.ctx, wk.mod)
	if err != nil {
		return nil, err
	}
	if options != nil {
		if options.ResourcesDir != "" {
			internal.UsvgOptionsSetResourcesDir(wk.ctx, wk.mod, o, options.ResourcesDir)
		}
		if options.Dpi != 0 {
			internal.UsvgOptionsSetDpi(wk.ctx, wk.mod, o, options.Dpi)
		}
		if options.FontFamily != "" {
			internal.UsvgOptionsSetFontFamily(wk.ctx, wk.mod, o, options.FontFamily)
		}
		if options.FontSize != 0 {
			internal.UsvgOptionsSetFontSize(wk.ctx, wk.mod, o, options.FontSize)
		}
		if options.Languages != nil && len(options.Languages) != 0 {
			internal.UsvgOptionsSetLanguages(wk.ctx, wk.mod, o, strings.Join(options.Languages, " "))
		}
		if options.ShapeRenderingMode != 0 {
			internal.UsvgOptionsSetShapeRenderingMode(wk.ctx, wk.mod, o, int32(options.ShapeRenderingMode))
		}
		if options.TextRenderingMode != 0 {
			internal.UsvgOptionsSetTextRenderingMode(wk.ctx, wk.mod, o, int32(options.TextRenderingMode))
		}
		if options.ImageRenderingMode != 0 {
			internal.UsvgOptionsSetImageRenderingMode(wk.ctx, wk.mod, o, int32(options.ImageRenderingMode))
		}
		if options.DefaultSizeWidth != 0 && options.DefaultSizeHeight != 0 {
			internal.UsvgOptionsSetDefaultSize(wk.ctx, wk.mod, o, options.DefaultSizeWidth, options.DefaultSizeHeight)
		}
	}
	defer internal.UsvgOptionsDelete(wk.ctx, wk.mod, o)
	t, err := internal.UsvgTreeFromData(wk.ctx, wk.mod, data, o)
	if err != nil {
		return nil, err
	}
	return &tree{wk, t}, nil
}

func (t *tree) Close() error {
	return internal.TinySkiaPixmapDelete(t.ctx, t.mod, t.ptr)
}

func (t *tree) ConvertText(fontdb *fontdb) error {
	return internal.UsvgTreeConvertText(t.ctx, t.mod, t.ptr, fontdb.ptr)
}

func (t *tree) GetSize() (float32, float32, error) {
	width, err := internal.UsvgTreeGetHeight(t.ctx, t.mod, t.ptr)
	if err != nil {
		return 0, 0, err
	}
	height, err := internal.UsvgTreeGetHeight(t.ctx, t.mod, t.ptr)
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

func (t *tree) Render(transform transform, pixmap *pixmap) error {
	rt, err := internal.ResvgTreeFromUsvg(t.ctx, t.mod, t.ptr)
	if err != nil {
		return err
	}
	defer internal.ResvgTreeDelete(t.ctx, t.mod, rt)
	tf, err := transform(t.ctx, t.mod)
	if err != nil {
		return err
	}
	defer internal.TinySkiaTransformDelete(t.ctx, t.mod, tf)
	return internal.ResvgTreeRender(t.ctx, t.mod, rt, tf, pixmap.ptr)
}

func (wk *worker) NewPixmap(width uint32, height uint32) (*pixmap, error) {
	pm, err := internal.TinySkiaPixmapNew(wk.ctx, wk.mod, width, height)
	if err != nil {
		return nil, err
	}
	return &pixmap{wk, pm}, nil
}

func (pm *pixmap) Close() error {
	return internal.TinySkiaPixmapDelete(pm.ctx, pm.mod, pm.ptr)
}

func (pm *pixmap) EncodePNG() ([]byte, error) {
	return internal.TinySkiaPixmapEncodePng(pm.ctx, pm.mod, pm.ptr)
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
