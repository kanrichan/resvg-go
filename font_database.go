package main

import (
	"runtime"

	"github.com/tetratelabs/wazero/api"
)

var (
	funcFontDatabaseNew          = wasi.ExportedFunction("__floattech_database_new")
	funcFontDatabaseFree         = wasi.ExportedFunction("__floattech_database_free")
	funcFontDatabaseLoadFontData = wasi.ExportedFunction("__floattech_database_load_font_data")
	funcFontDatabaseLen          = wasi.ExportedFunction("__floattech_database_len")
)

// FontDatabase FontDatabase
type FontDatabase struct {
	ptr *int32
}

// NewFontDatabase NewFontDatabase
func NewFontDatabase() (*FontDatabase, error) {
	r, err := funcFontDatabaseNew.Call(ctx)
	if err != nil {
		return nil, err
	}
	var o = &FontDatabase{ptr: new(int32)}
	*o.ptr = api.DecodeI32(r[0])
	runtime.SetFinalizer(o, func(o *FontDatabase) {
		o.Free()
	})
	return o, nil
}

// LoadFontData LoadFontData
func (o *FontDatabase) LoadFontData(data []byte) error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	rb, err := NewRustBytes(int32(len(data)))
	if err != nil {
		return err
	}
	defer rb.Free()
	err = rb.Write(data)
	if err != nil {
		return err
	}
	_, err = funcFontDatabaseLoadFontData.Call(
		ctx, api.EncodeI32(*o.ptr), api.EncodeI32(*rb.ptr), api.EncodeI32(rb.len))
	return err
}

// Len Len
func (o *FontDatabase) Len() (int32, error) {
	if o.ptr == nil {
		return 0, ErrNullWasmPointer
	}
	r, err := funcFontDatabaseLen.Call(ctx, api.EncodeI32(*o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeI32(r[0]), nil
}

// Free Free
func (o *FontDatabase) Free() error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	if _, err := funcFontDatabaseFree.Call(
		ctx, uint64(*o.ptr)); err != nil {
		return err
	}
	o.ptr = nil
	return nil
}
