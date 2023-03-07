package main

import (
	"runtime"

	"github.com/tetratelabs/wazero/api"
)

var (
	funcSizeNew          = wasi.ExportedFunction("__floattech_size_new")
	funcSizeWidth        = wasi.ExportedFunction("__floattech_size_width")
	funcSizeHeight       = wasi.ExportedFunction("__floattech_size_height")
	funcSizeToScreenSize = wasi.ExportedFunction("__floattech_size_to_screen_size")
	funcSizeFree         = wasi.ExportedFunction("__floattech_size_free")

	funcScreenSizeNew    = wasi.ExportedFunction("__floattech_screen_size_new")
	funcScreenSizeWidth  = wasi.ExportedFunction("__floattech_screen_size_width")
	funcScreenSizeHeight = wasi.ExportedFunction("__floattech_screen_size_height")
	funcScreenSizeFree   = wasi.ExportedFunction("__floattech_screen_size_free")
)

// Size Size
type Size struct {
	ptr *int32
}

// NewSize NewSize
func NewSize(width float64, height float64) (*Size, error) {
	r, err := funcSizeNew.Call(
		ctx, api.EncodeF64(width), api.EncodeF64(height))
	if err != nil {
		return nil, err
	}
	var o = &Size{ptr: new(int32)}
	*o.ptr = api.DecodeI32(r[0])
	runtime.SetFinalizer(o, func(o *Size) {
		o.Free()
	})
	return o, nil
}

// Width Width
func (o Size) Width() (float64, error) {
	if o.ptr == nil {
		return 0, ErrNullWasmPointer
	}
	r, err := funcSizeWidth.Call(
		ctx, api.EncodeI32(*o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeF64(r[0]), nil
}

// Height Height
func (o Size) Height() (float64, error) {
	if o.ptr == nil {
		return 0, ErrNullWasmPointer
	}
	r, err := funcSizeHeight.Call(
		ctx, api.EncodeI32(*o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeF64(r[0]), nil
}

// ToScreenSize ToScreenSize
func (o Size) ToScreenSize() (*ScreenSize, error) {
	if o.ptr == nil {
		return nil, ErrNullWasmPointer
	}
	r, err := funcSizeToScreenSize.Call(
		ctx, api.EncodeI32(*o.ptr))
	if err != nil {
		return nil, err
	}
	var oo = &ScreenSize{ptr: new(int32)}
	*oo.ptr = api.DecodeI32(r[0])
	runtime.SetFinalizer(oo, func(oo *ScreenSize) {
		oo.Free()
	})
	return oo, nil
}

// Free Free
func (o Size) Free() error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	if _, err := funcSizeFree.Call(
		ctx, uint64(*o.ptr)); err != nil {
		return err
	}
	o.ptr = nil
	return nil
}

// ScreenSize ScreenSize
type ScreenSize struct {
	ptr *int32
}

// NewScreenSize NewScreenSize
func NewScreenSize(width uint32, height uint32) (*ScreenSize, error) {
	r, err := funcScreenSizeNew.Call(ctx, api.EncodeU32(width), api.EncodeU32(height))
	if err != nil {
		return nil, err
	}
	var o = &ScreenSize{ptr: new(int32)}
	*o.ptr = api.DecodeI32(r[0])
	runtime.SetFinalizer(o, func(o *ScreenSize) {
		o.Free()
	})
	return o, nil
}

// Width Width
func (o ScreenSize) Width() (uint32, error) {
	if o.ptr == nil {
		return 0, ErrNullWasmPointer
	}
	r, err := funcScreenSizeWidth.Call(
		ctx, api.EncodeI32(*o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeU32(r[0]), nil
}

// Height Height
func (o ScreenSize) Height() (uint32, error) {
	if o.ptr == nil {
		return 0, ErrNullWasmPointer
	}
	r, err := funcScreenSizeHeight.Call(
		ctx, api.EncodeI32(*o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeU32(r[0]), nil
}

// Free Free
func (o ScreenSize) Free() error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	if _, err := funcScreenSizeFree.Call(
		ctx, uint64(*o.ptr)); err != nil {
		return err
	}
	o.ptr = nil
	return nil
}
