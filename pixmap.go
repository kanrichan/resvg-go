package main

import (
	"runtime"

	"github.com/tetratelabs/wazero/api"
)

var (
	funcPixmapNew       = wasi.ExportedFunction("__floattech_pixmap_new")
	funcPixmapEncodePNG = wasi.ExportedFunction("__floattech_pixmap_encode_png")
	funcPixmapFree      = wasi.ExportedFunction("__floattech_pixmap_free")
)

// Pixmap Pixmap
type Pixmap struct {
	ptr *int32
}

// NewPixmap NewPixmap
func NewPixmap(width uint32, height uint32) (*Pixmap, error) {
	r, err := funcPixmapNew.Call(
		ctx, api.EncodeU32(width), api.EncodeU32(height))
	if err != nil {
		return nil, err
	}
	var o = &Pixmap{ptr: new(int32)}
	*o.ptr = api.DecodeI32(r[0])
	runtime.SetFinalizer(o, func(o *Size) {
		o.Free()
	})
	return o, nil
}

// EncodePNG EncodePNG
func (o Pixmap) EncodePNG() ([]byte, error) {
	rb, err := NewRustBytesPointer()
	if err != nil {
		return nil, err
	}
	defer rb.Free()
	_, err = funcPixmapEncodePNG.Call(
		ctx, api.EncodeI32(*rb.ptr), api.EncodeI32(*o.ptr))
	if err != nil {
		return nil, err
	}
	return rb.Read()
}

// Free Free
func (o Pixmap) Free() error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	if _, err := funcPixmapFree.Call(
		ctx, uint64(*o.ptr)); err != nil {
		return err
	}
	o.ptr = nil
	return nil
}
