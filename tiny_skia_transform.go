package resvg

import "github.com/tetratelabs/wazero/api"

// TinySkiaTransform TinySkiaTransform
type TinySkiaTransform struct {
	ptr  int32
	free bool
	inst *instance
}

// TinySkiaTransformDefault TinySkiaTransformDefault
func (inst *instance) TinySkiaTransformDefault() (*TinySkiaTransform, error) {
	fn := inst.mod.ExportedFunction("__tiny_skia_transform_default")
	r, err := fn.Call(inst.ctx)
	if err != nil {
		return nil, err
	}
	return &TinySkiaTransform{api.DecodeI32(r[0]), false, inst}, nil
}

// TinySkiaTransformFromRow TinySkiaTransformFromRow
func (inst *instance) TinySkiaTransformFromRow(sx, ky, kx, sy, tx, ty float32) (*TinySkiaTransform, error) {
	fn := inst.mod.ExportedFunction("__tiny_skia_transform_from_row")
	r, err := fn.Call(
		inst.ctx,
		api.EncodeF32(sx),
		api.EncodeF32(ky),
		api.EncodeF32(kx),
		api.EncodeF32(sy),
		api.EncodeF32(tx),
		api.EncodeF32(ty),
	)
	if err != nil {
		return nil, err
	}
	return &TinySkiaTransform{api.DecodeI32(r[0]), false, inst}, nil
}

// TinySkiaTransformFromTranslate TinySkiaTransformFromTranslate
func (inst *instance) TinySkiaTransformFromTranslate(tx, ty float32) (*TinySkiaTransform, error) {
	fn := inst.mod.ExportedFunction("__tiny_skia_transform_from_translate")
	r, err := fn.Call(
		inst.ctx,
		api.EncodeF32(tx),
		api.EncodeF32(ty),
	)
	if err != nil {
		return nil, err
	}
	return &TinySkiaTransform{api.DecodeI32(r[0]), false, inst}, nil
}

// TinySkiaTransformFromScale TinySkiaTransformFromScale
func (inst *instance) TinySkiaTransformFromScale(sx, sy float32) (*TinySkiaTransform, error) {
	fn := inst.mod.ExportedFunction("__tiny_skia_transform_from_scale")
	r, err := fn.Call(
		inst.ctx,
		api.EncodeF32(sx),
		api.EncodeF32(sy),
	)
	if err != nil {
		return nil, err
	}
	return &TinySkiaTransform{api.DecodeI32(r[0]), false, inst}, nil
}

// TinySkiaTransformFromSkew TinySkiaTransformFromSkew
func (inst *instance) TinySkiaTransformFromSkew(kx, ky float32) (*TinySkiaTransform, error) {
	fn := inst.mod.ExportedFunction("__tiny_skia_transform_from_skew")
	r, err := fn.Call(
		inst.ctx,
		api.EncodeF32(kx),
		api.EncodeF32(ky),
	)
	if err != nil {
		return nil, err
	}
	return &TinySkiaTransform{api.DecodeI32(r[0]), false, inst}, nil
}

// TinySkiaTransformFromRotate TinySkiaTransformFromRotate
func (inst *instance) TinySkiaTransformFromRotate(angle float32) (*TinySkiaTransform, error) {
	fn := inst.mod.ExportedFunction("__tiny_skia_transform_from_rotate")
	r, err := fn.Call(
		inst.ctx,
		api.EncodeF32(angle),
	)
	if err != nil {
		return nil, err
	}
	return &TinySkiaTransform{api.DecodeI32(r[0]), false, inst}, nil
}

// TinySkiaTransformFromRotate TinySkiaTransformFromRotate
func (inst *instance) TinySkiaTransformFromRotateAt(angle, tx, ty float32) (*TinySkiaTransform, error) {
	fn := inst.mod.ExportedFunction("__tiny_skia_transform_from_rotate_at")
	r, err := fn.Call(
		inst.ctx,
		api.EncodeF32(angle),
		api.EncodeF32(tx),
		api.EncodeF32(ty),
	)
	if err != nil {
		return nil, err
	}
	return &TinySkiaTransform{api.DecodeI32(r[0]), false, inst}, nil
}

// Free Free
func (o *TinySkiaTransform) Free() error {
	if o.free {
		return ErrNullWasmPointer
	}
	fn := o.inst.mod.ExportedFunction("__tiny_skia_transform_free")
	_, err := fn.Call(o.inst.ctx, uint64(o.ptr))
	if err != nil {
		return err
	}
	o.free = true
	return nil
}
