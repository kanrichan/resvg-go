package resvg

import "github.com/tetratelabs/wazero/api"

// UsvgSize UsvgSize
type UsvgSize struct {
	ptr  int32
	free bool
	inst *Resvg
}

// NewUsvgSize NewUsvgSize
func (inst *Resvg) NewUsvgSize(width, height float64) (*UsvgSize, error) {
	fn := inst.mod.ExportedFunction("__usvg_size_new")
	r, err := fn.Call(
		inst.ctx,
		api.EncodeF64(width),
		api.EncodeF64(height),
	)
	if err != nil {
		return nil, err
	}
	return &UsvgSize{api.DecodeI32(r[0]), false, inst}, nil
}

// Free Free
func (o *UsvgSize) Free() error {
	if o.free {
		return ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_size_free")
	_, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o.free = true
	return nil
}

// Width Width
func (o *UsvgSize) Width() (float64, error) {
	if o.free {
		return 0, ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_size_width")
	r, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeF64(r[0]), nil
}

// Height Height
func (o *UsvgSize) Height() (float64, error) {
	if o.free {
		return 0, ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_size_height")
	r, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeF64(r[0]), nil
}

// ToScreenSize ToScreenSize
func (o *UsvgSize) ToScreenSize() (*UsvgScreenSize, error) {
	if o.free {
		return nil, ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_size_to_screen_size")
	r, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return nil, err
	}
	return &UsvgScreenSize{api.DecodeI32(r[0]), false, o.inst}, nil
}
