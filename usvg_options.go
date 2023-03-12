package resvg

import "github.com/tetratelabs/wazero/api"

// UsvgOptions UsvgOptions
type UsvgOptions struct {
	ptr  int32
	free bool
	inst *Resvg
}

// UsvgOptionsDefault UsvgOptionsDefault
func (inst *Resvg) UsvgOptionsDefault() (*UsvgOptions, error) {
	fn := inst.mod.ExportedFunction("__usvg_options_default")
	r, err := fn.Call(inst.ctx)
	if err != nil {
		return nil, err
	}
	return &UsvgOptions{api.DecodeI32(r[0]), false, inst}, nil
}

// Free Free
func (o *UsvgOptions) Free() error {
	if o.free {
		return ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_options_free")
	_, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o.free = true
	return nil
}

// SetDPI SetDPI
func (o *UsvgOptions) SetDPI(dpi float64) error {
	if o.free {
		return ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_options_set_dpi")
	_, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr), api.EncodeF64(dpi))
	return err
}

// SetFontFamily SetFontFamily
func (o *UsvgOptions) SetFontFamily(family string) error {
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
	fn := o.inst.mod.ExportedFunction("__usvg_options_set_font_family")
	_, err = fn.Call(o.inst.ctx, api.EncodeI32(o.ptr), api.EncodeI32(rb.ptr), api.EncodeI32(rb.len))
	if err != nil {
		return err
	}
	rb.free = true // free by rust wasm
	return nil
}

// SetFontSize SetFontSize
func (o *UsvgOptions) SetFontSize(size float64) error {
	if o.free {
		return ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_options_set_font_size")
	_, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr), api.EncodeF64(size))
	return err
}

// SetKeepNamedGroup SetKeepNamedGroup
func (o *UsvgOptions) SetKeepNamedGroup(keep bool) error {
	if o.free {
		return ErrNullWasmPointer
	}
	var k int32
	if keep {
		k = 1
	}
	fn := o.inst.mod.ExportedFunction("__usvg_options_set_keep_named_groups")
	_, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr), api.EncodeI32(k))
	return err
}

// SetDefaultSize SetDefaultSize
func (o *UsvgOptions) SetDefaultSize(size *UsvgSize) error {
	if o.free {
		return ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_options_set_default_size")
	_, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr), api.EncodeI32(size.ptr))
	if err != nil {
		return err
	}
	size.free = true // free by rust wasm
	return nil
}
