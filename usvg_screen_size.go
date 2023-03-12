package resvg

import "github.com/tetratelabs/wazero/api"

// UsvgScreenSize UsvgScreenSize
type UsvgScreenSize struct {
	ptr  int32
	free bool
	inst *instance
}

// NewUsvgScreenSize NewUsvgScreenSize
func (inst *instance) NewUsvgScreenSize(width, height float64) (*UsvgScreenSize, error) {
	fn := inst.mod.ExportedFunction("__usvg_screen_size_new")
	r, err := fn.Call(
		inst.ctx,
		api.EncodeF64(width),
		api.EncodeF64(height),
	)
	if err != nil {
		return nil, err
	}
	return &UsvgScreenSize{api.DecodeI32(r[0]), false, inst}, nil
}

// Free Free
func (o *UsvgScreenSize) Free() error {
	if o.free {
		return ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_screen_size_free")
	_, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o.free = true
	return nil
}

// Width Width
func (o *UsvgScreenSize) Width() (uint32, error) {
	if o.free {
		return 0, ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_screen_size_width")
	r, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeU32(r[0]), nil
}

// Height Height
func (o *UsvgScreenSize) Height() (uint32, error) {
	if o.free {
		return 0, ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_screen_size_height")
	r, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return 0, err
	}
	return api.DecodeU32(r[0]), nil
}
