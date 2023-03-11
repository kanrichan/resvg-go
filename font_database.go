package resvg

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
type FontDatabase int32

// NewFontDatabase NewFontDatabase
func NewFontDatabase() (*FontDatabase, error) {
	r, err := funcFontDatabaseNew.Call(ctx)
	if err != nil {
		return nil, err
	}
	o := FontDatabase(api.DecodeI32(r[0]))
	runtime.SetFinalizer(&o, func(o *FontDatabase) {
		o.Free()
	})
	return &o, nil
}

// LoadFontData LoadFontData
func (o *FontDatabase) LoadFontData(data []byte) error {
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
		ctx, api.EncodeI32(int32(*o)),
		api.EncodeI32(rb.ptr),
		api.EncodeI32(rb.len),
	)
	return err
}

// Len Len
func (o *FontDatabase) Len() (int32, error) {
	r, err := funcFontDatabaseLen.Call(ctx, api.EncodeI32(int32(*o)))
	if err != nil {
		return 0, err
	}
	return api.DecodeI32(r[0]), nil
}

// Free Free
func (o *FontDatabase) Free() error {
	_, err := funcFontDatabaseFree.Call(ctx, uint64(*o))
	return err
}
