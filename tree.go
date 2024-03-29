package resvg

import (
	_ "embed"
	"path/filepath"
	"strings"

	"github.com/kanrichan/resvg-go/internal"
)

// Tree SVG tree
type Tree struct {
	wk  *Worker
	ptr int32
}

// NewTreeFromData parses `Tree` from an SVG data.
// Can contain a gzip compressed data.
func (wk *Worker) NewTreeFromData(data []byte, options *Options) (*Tree, error) {
	if !wk.used.CompareAndSwap(false, true) {
		return nil, ErrWorkerIsBeingUsed
	}
	defer wk.used.Store(false)
	o, err := internal.UsvgOptionsDefault(wk.ctx, wk.mod)
	if err != nil {
		return nil, err
	}
	defer internal.UsvgOptionsDelete(wk.ctx, wk.mod, o)
	if options != nil {
		if options.ResourcesDir != "" {
			p, err := filepath.Abs(options.ResourcesDir)
			if err != nil {
				return nil, err
			}
			internal.UsvgOptionsSetResourcesDir(
				wk.ctx, wk.mod, o,
				p,
			)
		}
		if options.Dpi != 0 {
			internal.UsvgOptionsSetDpi(
				wk.ctx, wk.mod, o,
				options.Dpi,
			)
		}
		if options.FontFamily != "" {
			internal.UsvgOptionsSetFontFamily(
				wk.ctx, wk.mod, o,
				options.FontFamily,
			)
		}
		if options.FontSize != 0 {
			internal.UsvgOptionsSetFontSize(
				wk.ctx, wk.mod, o,
				options.FontSize,
			)
		}
		if options.Languages != nil && len(options.Languages) != 0 {
			internal.UsvgOptionsSetLanguages(
				wk.ctx, wk.mod, o,
				strings.Join(options.Languages, " "),
			)
		}
		if options.ShapeRenderingMode != 0 {
			internal.UsvgOptionsSetShapeRenderingMode(
				wk.ctx, wk.mod, o,
				int32(options.ShapeRenderingMode),
			)
		}
		if options.TextRenderingMode != 0 {
			internal.UsvgOptionsSetTextRenderingMode(
				wk.ctx, wk.mod, o,
				int32(options.TextRenderingMode),
			)
		}
		if options.ImageRenderingMode != 0 {
			internal.UsvgOptionsSetImageRenderingMode(
				wk.ctx, wk.mod, o,
				int32(options.ImageRenderingMode),
			)
		}
		if options.DefaultSizeWidth != 0 && options.DefaultSizeHeight != 0 {
			internal.UsvgOptionsSetDefaultSize(
				wk.ctx, wk.mod, o, options.DefaultSizeWidth,
				options.DefaultSizeHeight,
			)
		}
	}
	t, err := internal.UsvgTreeFromData(wk.ctx, wk.mod, data, o)
	if err != nil {
		return nil, err
	}
	return &Tree{wk, t}, nil
}

// Close cloes the `Tree` and recovers memory.
func (t *Tree) Close() error {
	if !t.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer t.wk.used.Store(false)
	if t.ptr == 0 {
		return ErrPointerIsNil
	}
	err := internal.UsvgTreeDelete(t.wk.ctx, t.wk.mod, t.ptr)
	if err != nil {
		return err
	}
	t.ptr = 0
	return nil
}

// ConvertText converts text nodes into `Tree`.
func (t *Tree) ConvertText(fontdb *FontDB) error {
	if t.wk != fontdb.wk {
		return ErrPointerIsNil
	}
	if !t.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer t.wk.used.Store(false)
	if t.ptr == 0 {
		return ErrPointerIsNil
	}
	if fontdb.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.UsvgTreeConvertText(t.wk.ctx, t.wk.mod, t.ptr, fontdb.ptr)
}

// GetSize returns Tree's width and height.
func (t *Tree) GetSize() (float32, float32, error) {
	if !t.wk.used.CompareAndSwap(false, true) {
		return 0, 0, ErrWorkerIsBeingUsed
	}
	defer t.wk.used.Store(false)
	if t.ptr == 0 {
		return 0, 0, ErrPointerIsNil
	}
	width, err := internal.UsvgTreeGetWidth(t.wk.ctx, t.wk.mod, t.ptr)
	if err != nil {
		return 0, 0, err
	}
	height, err := internal.UsvgTreeGetHeight(t.wk.ctx, t.wk.mod, t.ptr)
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

// Render renders the tree onto the pixmap.
func (t *Tree) Render(transform transform, pixmap *Pixmap) error {
	if t.wk != pixmap.wk {
		return ErrPointerIsNil
	}
	if !t.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer t.wk.used.Store(false)
	if t.ptr == 0 {
		return ErrPointerIsNil
	}
	if pixmap.ptr == 0 {
		return ErrPointerIsNil
	}
	rt, err := internal.ResvgTreeFromUsvg(t.wk.ctx, t.wk.mod, t.ptr)
	if err != nil {
		return err
	}
	defer internal.ResvgTreeDelete(t.wk.ctx, t.wk.mod, rt)
	tf, err := transform(t.wk.ctx, t.wk.mod)
	if err != nil {
		return err
	}
	defer internal.TinySkiaTransformDelete(t.wk.ctx, t.wk.mod, tf)
	return internal.ResvgTreeRender(t.wk.ctx, t.wk.mod, rt, tf, pixmap.ptr)
}
