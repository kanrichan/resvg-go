package resvg

import "github.com/tetratelabs/wazero/api"

// UsvgFitTo UsvgFitTo
type UsvgFitTo struct {
	ptr  int32
	free bool
	inst *Resvg
}

// UsvgFitToOriginal UsvgFitToOriginal
func (inst *Resvg) UsvgFitToOriginal() (*UsvgFitTo, error) {
	fn := inst.mod.ExportedFunction("__usvg_fit_to_original")
	r, err := fn.Call(inst.ctx)
	if err != nil {
		return nil, err
	}
	return &UsvgFitTo{api.DecodeI32(r[0]), false, inst}, nil
}

// UsvgFitToWidth UsvgFitToWidth
func (inst *Resvg) UsvgFitToWidth(width uint32) (*UsvgFitTo, error) {
	fn := inst.mod.ExportedFunction("__usvg_fit_to_width")
	r, err := fn.Call(inst.ctx, api.EncodeU32(width))
	if err != nil {
		return nil, err
	}
	return &UsvgFitTo{api.DecodeI32(r[0]), false, inst}, nil
}

// UsvgFitToHeight UsvgFitToHeight
func (inst *Resvg) UsvgFitToHeight(height uint32) (*UsvgFitTo, error) {
	fn := inst.mod.ExportedFunction("__usvg_fit_to_height")
	r, err := fn.Call(inst.ctx, api.EncodeU32(height))
	if err != nil {
		return nil, err
	}
	return &UsvgFitTo{api.DecodeI32(r[0]), false, inst}, nil
}

// UsvgFitToSize UsvgFitToSize
func (inst *Resvg) UsvgFitToSize(width, height uint32) (*UsvgFitTo, error) {
	fn := inst.mod.ExportedFunction("__usvg_fit_to_size")
	r, err := fn.Call(inst.ctx, api.EncodeU32(width), api.EncodeU32(height))
	if err != nil {
		return nil, err
	}
	return &UsvgFitTo{api.DecodeI32(r[0]), false, inst}, nil
}

// UsvgFitToZoom UsvgFitToZoom
func (inst *Resvg) UsvgFitToZoom(zoom float32) (*UsvgFitTo, error) {
	fn := inst.mod.ExportedFunction("__usvg_fit_to_zoom")
	r, err := fn.Call(inst.ctx, api.EncodeF32(zoom))
	if err != nil {
		return nil, err
	}
	return &UsvgFitTo{api.DecodeI32(r[0]), false, inst}, nil
}

// Free Free
func (o *UsvgFitTo) Free() error {
	if o.free {
		return ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_fit_to_free")
	_, err := fn.Call(o.inst.ctx, uint64(o.ptr))
	if err != nil {
		return err
	}
	o.free = true
	return nil
}
