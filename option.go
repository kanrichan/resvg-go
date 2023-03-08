package resvg

import (
	"runtime"

	"github.com/tetratelabs/wazero/api"
)

var (
	funcOptionDefault        = wasi.ExportedFunction("__floattech_option_default")
	funcOptionSetFontFamily  = wasi.ExportedFunction("__floattech_option_set_font_family")
	funcOptionSetDefaultSize = wasi.ExportedFunction("__floattech_option_set_default_size")
	funcOptionFree           = wasi.ExportedFunction("__floattech_option_free")
)

// Options Options
type Options int32

// NewDefaultOptions NewDefaultOptions
func NewDefaultOptions() (*Options, error) {
	r, err := funcOptionDefault.Call(ctx)
	if err != nil {
		return nil, err
	}
	o := Options(api.DecodeI32(r[0]))
	runtime.SetFinalizer(&o, func(o *Options) {
		o.Free()
	})
	return &o, nil
}

// SetFontFamily SetFontFamily
func (o *Options) SetFontFamily(fontFamily string) error {
	data := []byte(fontFamily)
	rb, err := NewRustBytes(int32(len(data)))
	if err != nil {
		return err
	}
	defer rb.Free()
	err = rb.Write(data)
	if err != nil {
		return err
	}
	_, err = funcOptionSetFontFamily.Call(
		ctx, api.EncodeI32(int32(*o)),
		api.EncodeI32(rb.ptr),
		api.EncodeI32(rb.len),
	)
	return err
}

// SetDefaultSize SetDefaultSize
func (o *Options) SetDefaultSize(size *Size) error {
	_, err := funcOptionSetDefaultSize.Call(
		ctx, api.EncodeI32(int32(*o)),
		api.EncodeI32(int32(*size)),
	)
	return err
}

// Free Free
func (o *Options) Free() error {
	_, err := funcOptionFree.Call(ctx, uint64(*o))
	return err
}
