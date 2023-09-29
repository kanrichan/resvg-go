package internal

import (
	"context"
	_ "embed"
	"errors"

	"github.com/tetratelabs/wazero/api"
)

// wasmzip was compiled using `cargo build --release --target wasm32-wasi`
// and packed into gzip
//
//go:embed resvg.wasm.gz
var WasmGZ []byte

var (
	ErrWasmFunctionNotFound = errors.New("wasm error: function not found")
	ErrWasmReturnInvaild    = errors.New("wasm error: return invalid")
	ErrWasmMemoryOutOfRange = errors.New("wasm error: memory out of range")
)

const (
	ExportNameFontdbDatabaseDefault            = "fontdb_database_default"
	ExportNameFontdbDatabaseDelete             = "fontdb_database_delete"
	ExportNameFontdbDatabaseLoadFontData       = "fontdb_database_load_font_data"
	ExportNameFontdbDatabaseLoadFontFile       = "fontdb_database_load_font_file"
	ExportNameFontdbDatabaseLoadFontsDir       = "fontdb_database_load_fonts_dir"
	ExportNameFontdbDatabaseLen                = "fontdb_database_len"
	ExportNameFontdbDatabaseSetSerifFamily     = "fontdb_database_set_serif_family"
	ExportNameFontdbDatabaseSetSansSerifFamily = "fontdb_database_set_sans_serif_family"
	ExportNameFontdbDatabaseSetCursiveFamily   = "fontdb_database_set_cursive_family"
	ExportNameFontdbDatabaseSetFantasyFamily   = "fontdb_database_set_fantasy_family"
	ExportNameFontdbDatabaseSetMonospaceFamily = "fontdb_database_set_monospace_family"
	ExportNameUsvgOptionsDefault               = "usvg_options_default"
	ExportNameUsvgOptionsDelete                = "usvg_options_delete"
	ExportNameUsvgOptionsSetResourcesDir       = "usvg_options_set_resources_dir"
	ExportNameUsvgOptionsSetDpi                = "usvg_options_set_dpi"
	ExportNameUsvgOptionsSetFontFamily         = "usvg_options_set_font_family"
	ExportNameUsvgOptionsSetFontSize           = "usvg_options_set_font_size"
	ExportNameUsvgOptionsSetLanguages          = "usvg_options_set_languages"
	ExportNameUsvgOptionsSetShapeRenderingMode = "usvg_options_set_shape_rendering_mode"
	ExportNameUsvgOptionsSetTextRenderingMode  = "usvg_options_set_text_rendering_mode"
	ExportNameUsvgOptionsSetImageRenderingMode = "usvg_options_set_image_rendering_mode"
	ExportNameUsvgOptionsSetDefaultSize        = "usvg_options_set_default_size"
	ExportNameTinySkiaPixmapNew                = "tiny_skia_pixmap_new"
	ExportNameTinySkiaPixmapDelete             = "tiny_skia_pixmap_delete"
	ExportNameTinySkiaPixmapEncodePng          = "tiny_skia_pixmap_encode_png"
	ExportNameTinySkiaPixmapGetWidth           = "tiny_skia_pixmap_get_width"
	ExportNameTinySkiaPixmapGetHeight          = "tiny_skia_pixmap_get_height"
	ExportNameTinySkiaTransformIdentity        = "tiny_skia_transform_identity"
	ExportNameTinySkiaTransformFromRow         = "tiny_skia_transform_from_row"
	ExportNameTinySkiaTransformFromTranslate   = "tiny_skia_transform_from_translate"
	ExportNameTinySkiaTransformFromScale       = "tiny_skia_transform_from_scale"
	ExportNameTinySkiaTransformFromSkew        = "tiny_skia_transform_from_skew"
	ExportNameTinySkiaTransformFromRotate      = "tiny_skia_transform_from_rotate"
	ExportNameTinySkiaTransformFromRotateAt    = "tiny_skia_transform_from_rotate_at"
	ExportNameTinySkiaTransformDelete          = "tiny_skia_transform_delete"
	ExportNameUsvgTreeFromData                 = "usvg_tree_from_data"
	ExportNameUsvgTreeDelete                   = "usvg_tree_delete"
	ExportNameUsvgTreeConvertText              = "usvg_tree_convert_text"
	ExportNameUsvgTreeGetWidth                 = "usvg_tree_get_size_width"
	ExportNameUsvgTreeGetHeight                = "usvg_tree_get_size_height"
	ExportNameResvgTreeFromUsvg                = "resvg_tree_from_usvg"
	ExportNameResvgTreeDelete                  = "resvg_tree_delete"
	ExportNameResvgTreeRender                  = "resvg_tree_render"
	ExportNameMemoryMalloc                     = "memory_malloc"
	ExportNameMemoryFree                       = "memory_free"
)

func FontdbDatabaseDefault(ctx context.Context, module api.Module) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseDefault)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func FontdbDatabaseDelete(ctx context.Context, module api.Module, database int32) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseDelete)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func FontdbDatabaseLoadFontData(ctx context.Context, module api.Module, database int32, data []byte) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseLoadFontData)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return err
	}
	if !module.Memory().Write(uint32(m), data) {
		return ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func FontdbDatabaseLoadFontFile(ctx context.Context, module api.Module, database int32, file string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseLoadFontFile)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	data := []byte(file)
	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return err
	}
	if !module.Memory().Write(uint32(m), data) {
		return ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func FontdbDatabaseLoadFontsDir(ctx context.Context, module api.Module, database int32, dir string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseLoadFontsDir)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	data := []byte(dir)
	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return err
	}
	if !module.Memory().Write(uint32(m), data) {
		return ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func FontdbDatabaseLen(ctx context.Context, module api.Module, database int32) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseLen)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func FontdbDatabaseSetSerifFamily(ctx context.Context, module api.Module, database int32, family string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseSetSerifFamily)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	data := []byte(family)
	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return err
	}
	if !module.Memory().Write(uint32(m), data) {
		return ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func FontdbDatabaseSetSansSerifFamily(ctx context.Context, module api.Module, database int32, family string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseSetSansSerifFamily)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	data := []byte(family)
	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return err
	}
	if !module.Memory().Write(uint32(m), data) {
		return ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func FontdbDatabaseSetCursiveFamily(ctx context.Context, module api.Module, database int32, family string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseSetCursiveFamily)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	data := []byte(family)
	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return err
	}
	if !module.Memory().Write(uint32(m), data) {
		return ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func FontdbDatabaseSetFantasyFamily(ctx context.Context, module api.Module, database int32, family string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseSetFantasyFamily)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	data := []byte(family)
	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return err
	}
	if !module.Memory().Write(uint32(m), data) {
		return ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func FontdbDatabaseSetMonospaceFamily(ctx context.Context, module api.Module, database int32, family string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseSetMonospaceFamily)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	data := []byte(family)
	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return err
	}
	if !module.Memory().Write(uint32(m), data) {
		return ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgOptionsDefault(ctx context.Context, module api.Module) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameUsvgOptionsDefault)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func UsvgOptionsDelete(ctx context.Context, module api.Module, options int32) error {
	fn := module.
		ExportedFunction(ExportNameUsvgOptionsDelete)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgOptionsSetResourcesDir(ctx context.Context, module api.Module, options int32, dir string) error {
	fn := module.
		ExportedFunction(ExportNameUsvgOptionsSetResourcesDir)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	data := []byte(dir)
	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return err
	}
	if !module.Memory().Write(uint32(m), data) {
		return ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgOptionsSetDpi(ctx context.Context, module api.Module, options int32, dpi float32) error {
	fn := module.
		ExportedFunction(ExportNameUsvgOptionsSetDpi)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeF32(dpi),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgOptionsSetFontFamily(ctx context.Context, module api.Module, options int32, family string) error {
	fn := module.
		ExportedFunction(ExportNameUsvgOptionsSetFontFamily)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	data := []byte(family)
	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return err
	}
	if !module.Memory().Write(uint32(m), data) {
		return ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgOptionsSetFontSize(ctx context.Context, module api.Module, options int32, size float32) error {
	fn := module.
		ExportedFunction(ExportNameUsvgOptionsSetFontSize)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeF32(size),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgOptionsSetLanguages(ctx context.Context, module api.Module, options int32, languages string) error {
	fn := module.
		ExportedFunction(ExportNameUsvgOptionsSetLanguages)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	data := []byte(languages)
	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return err
	}
	if !module.Memory().Write(uint32(m), data) {
		return ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgOptionsSetShapeRenderingMode(ctx context.Context, module api.Module, options int32, mode int32) error {
	fn := module.
		ExportedFunction(ExportNameUsvgOptionsSetShapeRenderingMode)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeI32(mode),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgOptionsSetTextRenderingMode(ctx context.Context, module api.Module, options int32, mode int32) error {
	fn := module.
		ExportedFunction(ExportNameUsvgOptionsSetTextRenderingMode)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeI32(mode),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgOptionsSetImageRenderingMode(ctx context.Context, module api.Module, options int32, mode int32) error {
	fn := module.
		ExportedFunction(ExportNameUsvgOptionsSetImageRenderingMode)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeI32(mode),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgOptionsSetDefaultSize(ctx context.Context, module api.Module, options int32, width float32, height float32) error {
	fn := module.
		ExportedFunction(ExportNameUsvgOptionsSetDefaultSize)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeF32(width),
		api.EncodeF32(height),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func TinySkiaPixmapNew(ctx context.Context, module api.Module, width uint32, height uint32) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaPixmapNew)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeU32(width),
		api.EncodeU32(height),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func TinySkiaPixmapDelete(ctx context.Context, module api.Module, pixmap int32) error {
	fn := module.
		ExportedFunction(ExportNameTinySkiaPixmapDelete)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(pixmap),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func TinySkiaPixmapEncodePng(ctx context.Context, module api.Module, pixmap int32) ([]byte, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaPixmapEncodePng)
	if fn == nil {
		return nil, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(pixmap),
	)
	if err != nil {
		return nil, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return nil, ErrWasmReturnInvaild
	}
	respptr := int32(resp[0] >> 32)
	resplen := uint32(resp[0])
	defer MemoryFree(ctx, module, respptr, int(resplen))
	b, f := module.Memory().Read(uint32(respptr), resplen)
	if !f {
		return nil, ErrWasmReturnInvaild
	}
	var data = make([]byte, int(resplen), int(resplen))
	copy(data, b)
	return data, nil
}

func TinySkiaPixmapGetWidth(ctx context.Context, module api.Module, pixmap int32) (uint32, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaPixmapGetWidth)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(pixmap),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeU32(resp[0]), nil
}

func TinySkiaPixmapGetHeight(ctx context.Context, module api.Module, pixmap int32) (uint32, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaPixmapGetHeight)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(pixmap),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeU32(resp[0]), nil
}

func TinySkiaTransformIdentity(ctx context.Context, module api.Module) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaTransformIdentity)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func TinySkiaTransformFromRow(ctx context.Context, module api.Module, sx float32, ky float32, kx float32, sy float32, tx float32, ty float32) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaTransformFromRow)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeF32(sx),
		api.EncodeF32(ky),
		api.EncodeF32(kx),
		api.EncodeF32(sy),
		api.EncodeF32(tx),
		api.EncodeF32(ty),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func TinySkiaTransformFromTranslate(ctx context.Context, module api.Module, tx float32, ty float32) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaTransformFromTranslate)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeF32(tx),
		api.EncodeF32(ty),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func TinySkiaTransformFromScale(ctx context.Context, module api.Module, width float32, height float32) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaTransformFromScale)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeF32(width),
		api.EncodeF32(height),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func TinySkiaTransformFromSkew(ctx context.Context, module api.Module, kx float32, ky float32) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaTransformFromSkew)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeF32(kx),
		api.EncodeF32(ky),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func TinySkiaTransformFromRotate(ctx context.Context, module api.Module, angle float32) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaTransformFromRotate)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeF32(angle),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func TinySkiaTransformFromRotateAt(ctx context.Context, module api.Module, angle float32, tx float32, ty float32) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaTransformFromRotateAt)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeF32(angle),
		api.EncodeF32(tx),
		api.EncodeF32(ty),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func TinySkiaTransformDelete(ctx context.Context, module api.Module, transform int32) error {
	fn := module.
		ExportedFunction(ExportNameTinySkiaTransformDelete)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(transform),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgTreeFromData(ctx context.Context, module api.Module, data []byte, options int32) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameUsvgTreeFromData)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}

	m, err := MemoryMalloc(ctx, module, len(data))
	if err != nil {
		return 0, err
	}
	if !module.Memory().Write(uint32(m), data) {
		return 0, ErrWasmMemoryOutOfRange
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
		api.EncodeI32(options),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func UsvgTreeDelete(ctx context.Context, module api.Module, tree int32) error {
	fn := module.
		ExportedFunction(ExportNameUsvgTreeDelete)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(tree),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgTreeConvertText(ctx context.Context, module api.Module, tree int32, database int32) error {
	fn := module.
		ExportedFunction(ExportNameUsvgTreeConvertText)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(tree),
		api.EncodeI32(database),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func UsvgTreeGetWidth(ctx context.Context, module api.Module, tree int32) (float32, error) {
	fn := module.
		ExportedFunction(ExportNameUsvgTreeGetWidth)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(tree),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeF32(resp[0]), nil
}

func UsvgTreeGetHeight(ctx context.Context, module api.Module, tree int32) (float32, error) {
	fn := module.
		ExportedFunction(ExportNameUsvgTreeGetHeight)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(tree),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeF32(resp[0]), nil
}

func ResvgTreeFromUsvg(ctx context.Context, module api.Module, tree int32) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameResvgTreeFromUsvg)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(tree),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func ResvgTreeDelete(ctx context.Context, module api.Module, rtree int32) error {
	fn := module.
		ExportedFunction(ExportNameResvgTreeDelete)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(rtree),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func ResvgTreeRender(ctx context.Context, module api.Module, rtree int32, transform int32, pixmap int32) error {
	fn := module.
		ExportedFunction(ExportNameResvgTreeRender)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(rtree),
		api.EncodeI32(transform),
		api.EncodeI32(pixmap),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}

func MemoryMalloc(ctx context.Context, module api.Module, size int) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameMemoryMalloc)
	if fn == nil {
		return 0, ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(int32(size)),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 1 || resp[0] == 0 {
		return 0, ErrWasmReturnInvaild
	}
	return api.DecodeI32(resp[0]), nil
}

func MemoryFree(ctx context.Context, module api.Module, ptr int32, size int) error {
	fn := module.
		ExportedFunction(ExportNameMemoryFree)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(ptr),
		api.EncodeI32(int32(size)),
	)
	if err != nil {
		return err
	}
	if len(resp) != 0 {
		return ErrWasmReturnInvaild
	}
	return nil
}
