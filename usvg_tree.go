package resvg

import (
	"github.com/tetratelabs/wazero/api"
)

// UsvgTree UsvgTree
type UsvgTree struct {
	ptr  int32
	free bool
	inst *instance
}

// NewUsvgTree NewUsvgTree
func (inst *instance) UsvgTreeFromData(data []byte, opt *UsvgOptions) (*UsvgTree, error) {
	rb, err := inst.NewRustBytes(int32(len(data)))
	if err != nil {
		return nil, err
	}
	err = rb.Write(data)
	if err != nil {
		return nil, err
	}
	fn := inst.mod.ExportedFunction("__usvg_tree_from_data")
	r, err := fn.Call(
		inst.ctx,
		api.EncodeI32(rb.ptr),
		api.EncodeI32(rb.len),
		api.EncodeI32(opt.ptr),
	)
	if err != nil {
		return nil, err
	}
	rb.free = true
	return &UsvgTree{api.DecodeI32(r[0]), false, inst}, nil
}

// Free Free
func (o *UsvgTree) Free() error {
	if o.free {
		return ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_tree_free")
	_, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr))
	if err != nil {
		return err
	}
	o.free = true
	return nil
}

// ConvertText ConvertText
func (o *UsvgTree) ConvertText(db *FontdbDatabase, keep bool) error {
	if o.free {
		return ErrNullWasmPointer
	}
	var k int32
	if keep {
		k = 1
	}
	fn := o.inst.mod.ExportedFunction("__usvg_tree_convert_text")
	_, err := fn.Call(
		o.inst.ctx,
		api.EncodeI32(o.ptr),
		api.EncodeI32(db.ptr),
		api.EncodeI32(k),
	)
	return err
}

// GetSizeClone GetSizeClone
func (o *UsvgTree) GetSizeClone() (*UsvgSize, error) {
	if o.free {
		return nil, ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__usvg_tree_get_size_clone")
	r, err := fn.Call(
		o.inst.ctx,
		api.EncodeI32(o.ptr),
	)
	if err != nil {
		return nil, err
	}
	return &UsvgSize{api.DecodeI32(r[0]), false, o.inst}, nil
}
