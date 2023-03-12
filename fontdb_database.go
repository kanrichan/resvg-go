package resvg

import (
	"github.com/tetratelabs/wazero/api"
)

// FontdbDatabase FontdbDatabase
type FontdbDatabase struct {
	ptr  int32
	free bool
	inst *Resvg
}

// NewFontdbDatabase NewFontdbDatabase
func (inst *Resvg) NewFontdbDatabase() (*FontdbDatabase, error) {
	fn := inst.mod.ExportedFunction("__fontdb_database_new")
	r, err := fn.Call(inst.ctx)
	if err != nil {
		return nil, err
	}
	return &FontdbDatabase{api.DecodeI32(r[0]), false, inst}, nil
}

// Free Free
func (o *FontdbDatabase) Free() error {
	if o.free {
		return ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__fontdb_database_free")
	_, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o.free = true
	return nil
}

// LoadFontData LoadFontData
func (o *FontdbDatabase) LoadFontData(data []byte) error {
	if o.free {
		return ErrNullWasmPointer
	}
	rb, err := o.inst.NewRustBytes(int32(len(data)))
	if err != nil {
		return err
	}
	err = rb.Write(data)
	if err != nil {
		return err
	}
	fn := o.inst.mod.ExportedFunction("__fontdb_database_load_font_data")
	_, err = fn.Call(
		o.inst.ctx, api.EncodeI32(o.ptr),
		api.EncodeI32(rb.ptr),
		api.EncodeI32(rb.len),
	)
	if err != nil {
		return err
	}
	rb.free = true // free by rust wasm
	return nil
}

// Len Len
func (o *FontdbDatabase) Len() (int32, error) {
	if o.free {
		return 0, ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__fontdb_database_len")
	r, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeI32(r[0]), nil
}

// SetSerifFamily SetSerifFamily
func (o *FontdbDatabase) SetSerifFamily(family string) error {
	if o.free {
		return ErrNullWasmPointer
	}
	rb, err := o.inst.NewRustBytes(int32(len(family)))
	if err != nil {
		return err
	}
	err = rb.WriteString(family)
	if err != nil {
		return err
	}
	fn := o.inst.mod.ExportedFunction("__fontdb_database_set_serif_family")
	_, err = fn.Call(
		o.inst.ctx, api.EncodeI32(o.ptr),
		api.EncodeI32(rb.ptr),
		api.EncodeI32(rb.len),
	)
	if err != nil {
		return err
	}
	rb.free = true // free by rust wasm
	return nil
}

// SetSansSerifFamily SetSansSerifFamily
func (o *FontdbDatabase) SetSansSerifFamily(family string) error {
	if o.free {
		return ErrNullWasmPointer
	}
	rb, err := o.inst.NewRustBytes(int32(len(family)))
	if err != nil {
		return err
	}
	err = rb.WriteString(family)
	if err != nil {
		return err
	}
	fn := o.inst.mod.ExportedFunction("__fontdb_database_set_sans_serif_family")
	_, err = fn.Call(
		o.inst.ctx, api.EncodeI32(o.ptr),
		api.EncodeI32(rb.ptr),
		api.EncodeI32(rb.len),
	)
	if err != nil {
		return err
	}
	rb.free = true // free by rust wasm
	return nil
}

// SetCursiveFamily SetCursiveFamily
func (o *FontdbDatabase) SetCursiveFamily(family string) error {
	if o.free {
		return ErrNullWasmPointer
	}
	rb, err := o.inst.NewRustBytes(int32(len(family)))
	if err != nil {
		return err
	}
	err = rb.WriteString(family)
	if err != nil {
		return err
	}
	fn := o.inst.mod.ExportedFunction("__fontdb_database_set_cursive_family")
	_, err = fn.Call(
		o.inst.ctx, api.EncodeI32(o.ptr),
		api.EncodeI32(rb.ptr),
		api.EncodeI32(rb.len),
	)
	if err != nil {
		return err
	}
	rb.free = true // free by rust wasm
	return nil
}

// SetFantasyFamily SetFantasyFamily
func (o *FontdbDatabase) SetFantasyFamily(family string) error {
	if o.free {
		return ErrNullWasmPointer
	}
	rb, err := o.inst.NewRustBytes(int32(len(family)))
	if err != nil {
		return err
	}
	err = rb.WriteString(family)
	if err != nil {
		return err
	}
	fn := o.inst.mod.ExportedFunction("__fontdb_database_set_fantasy_family")
	_, err = fn.Call(
		o.inst.ctx, api.EncodeI32(o.ptr),
		api.EncodeI32(rb.ptr),
		api.EncodeI32(rb.len),
	)
	if err != nil {
		return err
	}
	rb.free = true // free by rust wasm
	return nil
}

// SetMonospaceFamily SetMonospaceFamily
func (o *FontdbDatabase) SetMonospaceFamily(family string) error {
	if o.free {
		return ErrNullWasmPointer
	}
	rb, err := o.inst.NewRustBytes(int32(len(family)))
	if err != nil {
		return err
	}
	defer rb.Free()
	err = rb.WriteString(family)
	if err != nil {
		return err
	}
	fn := o.inst.mod.ExportedFunction("__fontdb_database_set_monospace_family")
	_, err = fn.Call(
		o.inst.ctx, api.EncodeI32(o.ptr),
		api.EncodeI32(rb.ptr),
		api.EncodeI32(rb.len),
	)
	if err != nil {
		return err
	}
	rb.free = true // free by rust wasm
	return nil
}
