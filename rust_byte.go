package resvg

import (
	"encoding/binary"

	"github.com/tetratelabs/wazero/api"
)

// RustBytes RustBytes
type RustBytes struct {
	ptr  int32
	len  int32
	free bool
	inst *instance
}

// NewRustBytes RustBytes
func (inst *instance) NewRustBytes(size int32) (*RustBytes, error) {
	fn := inst.mod.ExportedFunction("__rust_bytes_new")
	r, err := fn.Call(inst.ctx, api.EncodeI32(size))
	if err != nil {
		return nil, err
	}
	return &RustBytes{api.DecodeI32(r[0]), size, false, inst}, nil
}

func (o *RustBytes) Read() ([]byte, error) {
	data, f := o.inst.mod.Memory().Read(uint32(o.ptr), uint32(o.len))
	if !f {
		return nil, ErrOutOfRange
	}
	return data, nil
}

// ReadString ReadString
func (o *RustBytes) ReadString() (string, error) {
	data, err := o.Read()
	return string(data), err
}

func (o *RustBytes) Write(data []byte) error {
	if !o.inst.mod.Memory().Write(uint32(o.ptr), []byte(data)) {
		return ErrOutOfRange
	}
	return nil
}

// WriteString WriteString
func (o *RustBytes) WriteString(data string) error {
	return o.Write([]byte(data))
}

// Free Free
func (o *RustBytes) Free() error {
	fn := o.inst.mod.ExportedFunction("__rust_bytes_free")
	_, err := fn.Call(o.inst.ctx, api.EncodeI32(o.ptr), uint64(o.len))
	return err
}

// RustBytesPointer RustBytesPointer
type RustBytesPointer struct {
	s *RustBytes
	o *RustBytes
}

// NewRustBytesPointer NewRustBytesPointer
func (inst *instance) NewRustBytesPointer() (*RustBytesPointer, error) {
	s, err := inst.NewRustBytes(8)
	return &RustBytesPointer{s, nil}, err
}

// Free Free
func (o *RustBytesPointer) Free() error {
	return o.s.Free()
}

func (o *RustBytesPointer) Read() ([]byte, error) {
	err := o.binding()
	if err != nil {
		return nil, err
	}
	return o.o.Read()
}

func (o *RustBytesPointer) Write(data []byte) error {
	err := o.binding()
	if err != nil {
		return err
	}
	return o.o.Write(data)
}

func (o *RustBytesPointer) binding() error {
	if o.o != nil {
		return nil
	}
	ob, f := o.s.inst.mod.Memory().Read(uint32(o.s.ptr), 4)
	if !f {
		return ErrOutOfRange
	}
	offset := binary.LittleEndian.Uint32(ob)
	lb, f := o.s.inst.mod.Memory().Read(uint32(o.s.ptr)+4, 4)
	if !f {
		return ErrOutOfRange
	}
	length := binary.LittleEndian.Uint32(lb)
	o.o = &RustBytes{
		ptr:  int32(offset),
		len:  int32(length),
		free: false,
		inst: o.s.inst,
	}
	return nil
}
