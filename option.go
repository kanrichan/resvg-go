package main

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
type Options struct {
	ptr *int32
}

// DefaultOptions DefaultOptions
func DefaultOptions() (*Options, error) {
	r, err := funcOptionDefault.Call(ctx)
	if err != nil {
		return nil, err
	}
	var o = &Options{ptr: new(int32)}
	*o.ptr = api.DecodeI32(r[0])
	runtime.SetFinalizer(o, func(o *Options) {
		o.Free()
	})
	return o, nil
}

// SetFontFamily SetFontFamily
func (o *Options) SetFontFamily(fontFamily string) error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
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
		ctx, api.EncodeI32(*o.ptr), api.EncodeI32(*rb.ptr), api.EncodeI32(rb.len))
	return err
}

// SetDefaultSize SetDefaultSize
func (o *Options) SetDefaultSize(size *Size) error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	if size.ptr == nil {
		return ErrNullWasmPointer
	}
	_, err := funcOptionSetDefaultSize.Call(
		ctx, api.EncodeI32(*o.ptr), api.EncodeI32(*size.ptr))
	if err != nil {
		return err
	}
	size.ptr = nil
	return nil
}

// Free Free
func (o *Options) Free() error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	if _, err := funcOptionFree.Call(
		ctx, uint64(*o.ptr)); err != nil {
		return err
	}
	o.ptr = nil
	return nil
}
