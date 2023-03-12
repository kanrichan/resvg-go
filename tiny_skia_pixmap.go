package resvg

import (
	"github.com/tetratelabs/wazero/api"
)

// TinySkiaPixmap TinySkiaPixmap
type TinySkiaPixmap struct {
	ptr  int32
	free bool
	inst *instance
}

// NewTinySkiaPixmap NewTinySkiaPixmap
func (inst *instance) NewTinySkiaPixmap(width, height uint32) (*TinySkiaPixmap, error) {
	fn := inst.mod.ExportedFunction("__tiny_skia_pixmap_new")
	r, err := fn.Call(
		inst.ctx,
		api.EncodeU32(width),
		api.EncodeU32(height),
	)
	if err != nil {
		return nil, err
	}
	return &TinySkiaPixmap{api.DecodeI32(r[0]), false, inst}, nil
}

// Free Free
func (o *TinySkiaPixmap) Free() error {
	if o.free {
		return ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__tiny_skia_pixmap_free")
	_, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o.free = true
	return nil
}

// EncodePNG EncodePNG
func (o *TinySkiaPixmap) EncodePNG() ([]byte, error) {
	if o.free {
		return nil, ErrNullWasmPointer
	}
	rb, err := o.inst.NewRustBytesPointer()
	if err != nil {
		return nil, err
	}
	defer rb.Free()
	fn := o.inst.mod.ExportedFunction("__tiny_skia_pixmap_encode_png")
	_, err = fn.Call(
		o.inst.ctx,
		api.EncodeI32(rb.s.ptr),
		api.EncodeI32(o.ptr),
	)
	if err != nil {
		return nil, err
	}
	return rb.Read()
}
