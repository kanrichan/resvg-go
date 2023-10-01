package internal

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

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
	ExportNameTinySkiaPixmapDecodePNG          = "tiny_skia_pixmap_decode_png"
	ExportNameTinySkiaPixmapDelete             = "tiny_skia_pixmap_delete"
	ExportNameTinySkiaPixmapEncodePNG          = "tiny_skia_pixmap_encode_png"
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
	m, err := MemoryMalloc(ctx, module, len(file)+1)
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, m, len(file)+1)
	if err := CStrWrite(ctx, module, m, file); err != nil {
		return err
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
	)
	if err != nil {
		return err
	}
	if len(resp) != 1 {
		return ErrWasmReturnInvaild
	}
	if resp[0] == 0 {
		return nil
	}
	error, err := CStrRead(ctx, module, int32(resp[0]))
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, int32(resp[0]), 4)
	return errors.New(error)
}

func FontdbDatabaseLoadFontsDir(ctx context.Context, module api.Module, database int32, dir string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseLoadFontsDir)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	m, err := MemoryMalloc(ctx, module, len(dir)+1)
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, m, len(dir)+1)
	if err := CStrWrite(ctx, module, m, dir); err != nil {
		return err
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
	)
	if err != nil {
		return err
	}
	if len(resp) != 1 {
		return ErrWasmReturnInvaild
	}
	if resp[0] == 0 {
		return nil
	}
	error, err := CStrRead(ctx, module, int32(resp[0]))
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, int32(resp[0]), 4)
	return errors.New(error)
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
	m, err := MemoryMalloc(ctx, module, len(family)+1)
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, m, len(family)+1)
	if err := CStrWrite(ctx, module, m, family); err != nil {
		return err
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
	)
	if err != nil {
		return err
	}
	if len(resp) != 1 {
		return ErrWasmReturnInvaild
	}
	if resp[0] == 0 {
		return nil
	}
	error, err := CStrRead(ctx, module, int32(resp[0]))
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, int32(resp[0]), 4)
	return errors.New(error)
}

func FontdbDatabaseSetSansSerifFamily(ctx context.Context, module api.Module, database int32, family string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseSetSansSerifFamily)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	m, err := MemoryMalloc(ctx, module, len(family)+1)
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, m, len(family)+1)
	if err := CStrWrite(ctx, module, m, family); err != nil {
		return err
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
	)
	if err != nil {
		return err
	}
	if len(resp) != 1 {
		return ErrWasmReturnInvaild
	}
	if resp[0] == 0 {
		return nil
	}
	error, err := CStrRead(ctx, module, int32(resp[0]))
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, int32(resp[0]), 4)
	return errors.New(error)
}

func FontdbDatabaseSetCursiveFamily(ctx context.Context, module api.Module, database int32, family string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseSetCursiveFamily)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	m, err := MemoryMalloc(ctx, module, len(family)+1)
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, m, len(family)+1)
	if err := CStrWrite(ctx, module, m, family); err != nil {
		return err
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
	)
	if err != nil {
		return err
	}
	if len(resp) != 1 {
		return ErrWasmReturnInvaild
	}
	if resp[0] == 0 {
		return nil
	}
	error, err := CStrRead(ctx, module, int32(resp[0]))
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, int32(resp[0]), 4)
	return errors.New(error)
}

func FontdbDatabaseSetFantasyFamily(ctx context.Context, module api.Module, database int32, family string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseSetFantasyFamily)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	m, err := MemoryMalloc(ctx, module, len(family)+1)
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, m, len(family)+1)
	if err := CStrWrite(ctx, module, m, family); err != nil {
		return err
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
	)
	if err != nil {
		return err
	}
	if len(resp) != 1 {
		return ErrWasmReturnInvaild
	}
	if resp[0] == 0 {
		return nil
	}
	error, err := CStrRead(ctx, module, int32(resp[0]))
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, int32(resp[0]), 4)
	return errors.New(error)
}

func FontdbDatabaseSetMonospaceFamily(ctx context.Context, module api.Module, database int32, family string) error {
	fn := module.
		ExportedFunction(ExportNameFontdbDatabaseSetMonospaceFamily)
	if fn == nil {
		return ErrWasmFunctionNotFound
	}
	m, err := MemoryMalloc(ctx, module, len(family)+1)
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, m, len(family)+1)
	if err := CStrWrite(ctx, module, m, family); err != nil {
		return err
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(database),
		api.EncodeI32(m),
	)
	if err != nil {
		return err
	}
	if len(resp) != 1 {
		return ErrWasmReturnInvaild
	}
	if resp[0] == 0 {
		return nil
	}
	error, err := CStrRead(ctx, module, int32(resp[0]))
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, int32(resp[0]), 4)
	return errors.New(error)
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
	m, err := MemoryMalloc(ctx, module, len(dir)+1)
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, m, len(dir)+1)
	if err := CStrWrite(ctx, module, m, dir); err != nil {
		return err
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeI32(m),
	)
	if err != nil {
		return err
	}
	if len(resp) != 1 {
		return ErrWasmReturnInvaild
	}
	if resp[0] == 0 {
		return nil
	}
	error, err := CStrRead(ctx, module, int32(resp[0]))
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, int32(resp[0]), 4)
	return errors.New(error)
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
	m, err := MemoryMalloc(ctx, module, len(family)+1)
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, m, len(family)+1)
	if err := CStrWrite(ctx, module, m, family); err != nil {
		return err
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeI32(m),
	)
	if err != nil {
		return err
	}
	if len(resp) != 1 {
		return ErrWasmReturnInvaild
	}
	if resp[0] == 0 {
		return nil
	}
	error, err := CStrRead(ctx, module, int32(resp[0]))
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, int32(resp[0]), 4)
	return errors.New(error)
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
	m, err := MemoryMalloc(ctx, module, len(languages)+1)
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, m, len(languages)+1)
	if err := CStrWrite(ctx, module, m, languages); err != nil {
		return err
	}
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(options),
		api.EncodeI32(m),
	)
	if err != nil {
		return err
	}
	if len(resp) != 1 {
		return ErrWasmReturnInvaild
	}
	if resp[0] == 0 {
		return nil
	}
	error, err := CStrRead(ctx, module, int32(resp[0]))
	if err != nil {
		return err
	}
	defer MemoryFree(ctx, module, int32(resp[0]), 4)
	return errors.New(error)
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
	r, err := MemoryMalloc(ctx, module, 8)
	if err != nil {
		return 0, err
	}
	defer MemoryFree(ctx, module, r, 8)
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(r),
		api.EncodeU32(width),
		api.EncodeU32(height),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 0 {
		return 0, ErrWasmReturnInvaild
	}
	result, err := Result32Read(ctx, module, r)
	if err != nil {
		return 0, err
	}
	if result.ok {
		return result.data, nil
	}
	error, err := CStrRead(ctx, module, result.data)
	if err != nil {
		return 0, err
	}
	defer MemoryFree(ctx, module, result.data, len(error)+1)
	return 0, errors.New(error)
}

func TinySkiaPixmapDecodePNG(ctx context.Context, module api.Module, data []byte) (int32, error) {
	fn := module.
		ExportedFunction(ExportNameTinySkiaPixmapDecodePNG)
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
	r, err := MemoryMalloc(ctx, module, 8)
	if err != nil {
		return 0, err
	}
	defer MemoryFree(ctx, module, r, 8)
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(r),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 0 {
		return 0, ErrWasmReturnInvaild
	}
	result, err := Result32Read(ctx, module, r)
	if err != nil {
		return 0, err
	}
	if result.ok {
		return result.data, nil
	}
	error, err := CStrRead(ctx, module, result.data)
	if err != nil {
		return 0, err
	}
	defer MemoryFree(ctx, module, result.data, len(error)+1)
	return 0, errors.New(error)
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
		ExportedFunction(ExportNameTinySkiaPixmapEncodePNG)
	if fn == nil {
		return nil, ErrWasmFunctionNotFound
	}
	r, err := MemoryMalloc(ctx, module, 8)
	if err != nil {
		return nil, err
	}
	defer MemoryFree(ctx, module, r, 8)
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(r),
		api.EncodeI32(pixmap),
	)
	if err != nil {
		return nil, err
	}
	if len(resp) != 0 {
		return nil, ErrWasmReturnInvaild
	}
	result, err := Result64Read(ctx, module, r)
	if err != nil {
		return nil, err
	}
	if result.ok {
		respptr := uint32(result.data >> 32)
		resplen := uint32(result.data)
		fmt.Println(respptr, resplen)
		defer MemoryFree(ctx, module, int32(respptr), int(resplen))
		b, f := module.Memory().Read(respptr, resplen)
		if !f {
			return nil, ErrWasmReturnInvaild
		}
		var data = make([]byte, int(resplen), int(resplen))
		copy(data, b)
		return data, nil
	}
	error, err := CStrRead(ctx, module, int32(result.data))
	if err != nil {
		return nil, err
	}
	defer MemoryFree(ctx, module, int32(result.data), len(error)+1)
	return nil, errors.New(error)
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
	r, err := MemoryMalloc(ctx, module, 8)
	if err != nil {
		return 0, err
	}
	defer MemoryFree(ctx, module, r, 8)
	resp, err := fn.Call(
		ctx,
		api.EncodeI32(r),
		api.EncodeI32(m),
		api.EncodeI32(int32(len(data))),
		api.EncodeI32(options),
	)
	if err != nil {
		return 0, err
	}
	if len(resp) != 0 {
		return 0, ErrWasmReturnInvaild
	}
	result, err := Result32Read(ctx, module, r)
	if err != nil {
		return 0, err
	}
	if result.ok {
		return result.data, nil
	}
	error, err := CStrRead(ctx, module, result.data)
	if err != nil {
		return 0, err
	}
	defer MemoryFree(ctx, module, result.data, len(error)+1)
	return 0, errors.New(error)
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

func CStrMalloc(ctx context.Context, module api.Module, data string) (int32, error) {
	ptr, err := MemoryMalloc(ctx, module, len(data)+1)
	if err != nil {
		return 0, err
	}
	if !module.Memory().WriteString(uint32(ptr), data) {
		return 0, ErrWasmMemoryOutOfRange
	}
	if !module.Memory().WriteByte(uint32(ptr)+uint32(len(data)), 0) {
		return 0, ErrWasmMemoryOutOfRange
	}
	return ptr, nil
}

func CStrFree(ctx context.Context, module api.Module, ptr int32) (string, error) {
	var i int
	for i = 0; i < 1024; i++ {
		b, f := module.Memory().ReadByte(uint32(ptr) + uint32(i))
		if b == 0 || !f {
			break
		}
	}
	if i == 0 {
		return "", ErrWasmMemoryOutOfRange
	}
	defer MemoryFree(ctx, module, int32(ptr), i)
	b, _ := module.Memory().Read(uint32(ptr), uint32(i))
	var e = make([]byte, i, i)
	copy(e, b)
	return string(e), nil
}

func CStrWrite(ctx context.Context, module api.Module, ptr int32, s string) error {
	if !module.Memory().WriteString(uint32(ptr), s) {
		return ErrWasmMemoryOutOfRange
	}
	if !module.Memory().WriteByte(uint32(ptr)+uint32(len(s)), 0) {
		return ErrWasmMemoryOutOfRange
	}
	return nil
}

func CStrRead(ctx context.Context, module api.Module, ptr int32) (string, error) {
	var i int
	for i = 0; i < 1024; i++ {
		b, f := module.Memory().ReadByte(uint32(ptr) + uint32(i))
		if b == 0 || !f {
			break
		}
	}
	if i == 0 {
		return "", ErrWasmMemoryOutOfRange
	}
	b, _ := module.Memory().Read(uint32(ptr), uint32(i))
	var e = make([]byte, i, i)
	copy(e, b)
	return string(e), nil
}

type Result32 struct {
	ok   bool
	data int32
}

func Result32Write(ctx context.Context, module api.Module, ptr int32, r Result32) error {
	if r.ok {
		if !module.Memory().WriteUint32Le(uint32(ptr), 0) {
			return ErrWasmMemoryOutOfRange
		}
	} else {
		if module.Memory().WriteUint32Le(uint32(ptr), 1) {
			return ErrWasmMemoryOutOfRange
		}
	}
	if module.Memory().WriteUint32Le(uint32(ptr)+4, uint32(r.data)) {
		return ErrWasmMemoryOutOfRange
	}
	return nil
}

func Result32Read(ctx context.Context, module api.Module, ptr int32) (Result32, error) {
	ok, f := module.Memory().ReadUint32Le(uint32(ptr))
	if !f {
		return Result32{}, ErrWasmMemoryOutOfRange
	}
	o, f := module.Memory().ReadUint32Le(uint32(ptr) + 4)
	if !f {
		return Result32{}, ErrWasmMemoryOutOfRange
	}
	return Result32{ok == 0, int32(o)}, nil
}

type Result64 struct {
	ok   bool
	data int64
}

func Result64Write(ctx context.Context, module api.Module, ptr int32, r Result64) error {
	if r.ok {
		if !module.Memory().WriteUint32Le(uint32(ptr), 0) {
			return ErrWasmMemoryOutOfRange
		}
	} else {
		if module.Memory().WriteUint32Le(uint32(ptr), 1) {
			return ErrWasmMemoryOutOfRange
		}
	}
	if module.Memory().WriteUint64Le(uint32(ptr)+4, uint64(r.data)) {
		return ErrWasmMemoryOutOfRange
	}
	return nil
}

func Result64Read(ctx context.Context, module api.Module, ptr int32) (Result64, error) {
	ok, f := module.Memory().ReadUint32Le(uint32(ptr))
	if !f {
		return Result64{}, ErrWasmMemoryOutOfRange
	}
	o, f := module.Memory().ReadUint64Le(uint32(ptr) + 8)
	if !f {
		return Result64{}, ErrWasmMemoryOutOfRange
	}
	return Result64{ok == 0, int64(o)}, nil
}
