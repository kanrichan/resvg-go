package main

import (
	"runtime"

	"github.com/tetratelabs/wazero/api"
)

var (
	funcTreeFromData    = wasi.ExportedFunction("__floattech_tree_from_data")
	funcTreeConvertText = wasi.ExportedFunction("__floattech_tree_convert_text")
	funcTreeGetSize     = wasi.ExportedFunction("__floattech_tree_get_size")
	funcTreeFree        = wasi.ExportedFunction("__floattech_tree_free")
)

// Tree Tree
type Tree struct {
	ptr *int32
}

// TreeFromData TreeFromData
func TreeFromData(data []byte, opt *Options) (*Tree, error) {
	rb, err := NewRustBytes(int32(len(data)))
	if err != nil {
		return nil, err
	}
	defer rb.Free()
	err = rb.Write(data)
	if err != nil {
		return nil, err
	}
	r, err := funcTreeFromData.Call(
		ctx, api.EncodeI32(*rb.ptr), api.EncodeI32(rb.len), api.EncodeI32(*opt.ptr))
	if err != nil {
		return nil, err
	}
	var o = &Tree{ptr: new(int32)}
	*o.ptr = api.DecodeI32(r[0])
	runtime.SetFinalizer(o, func(o *Tree) {
		o.Free()
	})
	return o, nil
}

// ConvertText ConvertText
func (o *Tree) ConvertText(db *FontDatabase, keepNamedGroups bool) error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	f := int32(0)
	if keepNamedGroups {
		f = 1
	}
	_, err := funcTreeConvertText.Call(
		ctx, api.EncodeI32(*o.ptr), api.EncodeI32(*db.ptr), api.EncodeI32(f))
	return err
}

// GetSize GetSize
func (o *Tree) GetSize() (*Size, error) {
	if o.ptr == nil {
		return nil, ErrNullWasmPointer
	}
	r, err := funcTreeGetSize.Call(ctx, api.EncodeI32(*o.ptr))
	if err != nil {
		return nil, err
	}
	var oo = &Size{ptr: new(int32)}
	*oo.ptr = api.DecodeI32(r[0])
	runtime.SetFinalizer(oo, func(oo *Size) {
		oo.Free()
	})
	return oo, nil
}

// Free Free
func (o *Tree) Free() error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	if _, err := funcTreeFree.Call(
		ctx, uint64(*o.ptr)); err != nil {
		return err
	}
	o.ptr = nil
	return nil
}
