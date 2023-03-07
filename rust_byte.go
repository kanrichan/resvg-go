package main

import (
	"encoding/binary"
	"runtime"

	"github.com/tetratelabs/wazero/api"
)

// RustBytes RustBytes
type RustBytes struct {
	ptr *int32
	len int32
}

var funcRustByteAllocate = wasi.ExportedFunction("__floattech_bytes_allocate")
var funcRustByteFree = wasi.ExportedFunction("__floattech_bytes_free")

// NewRustBytes RustBytes
func NewRustBytes(size int32) (*RustBytes, error) {
	r, err := funcRustByteAllocate.Call(
		ctx, api.EncodeI32(size))
	if err != nil {
		return nil, err
	}
	var o = &RustBytes{ptr: new(int32), len: size}
	*o.ptr = api.DecodeI32(r[0])
	runtime.SetFinalizer(o, func(o *RustBytes) {
		o.Free()
	})
	return o, nil
}

func (o *RustBytes) Read() ([]byte, error) {
	if o.ptr == nil {
		return nil, ErrNullWasmPointer
	}
	data, f := wasi.Memory().Read(uint32(*o.ptr), uint32(o.len))
	if !f {
		return nil, ErrOutOfRange
	}
	return data, nil
}

func (o *RustBytes) Write(data []byte) error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	if !wasi.Memory().Write(uint32(*o.ptr), []byte(data)) {
		return ErrOutOfRange
	}
	return nil
}

// Free Free
func (o *RustBytes) Free() error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	if _, err := funcRustByteFree.Call(
		ctx, uint64(*o.ptr), uint64(o.len)); err != nil {
		return err
	}
	o.ptr = nil
	return nil
}

// RustBytesPointer RustBytesPointer
type RustBytesPointer struct {
	ptr *int32
	o   *RustBytes
}

// NewRustBytesPointer NewRustBytesPointer
func NewRustBytesPointer() (*RustBytesPointer, error) {
	r, err := funcRustByteAllocate.Call(ctx, 8)
	if err != nil {
		return nil, err
	}
	var o = &RustBytesPointer{ptr: new(int32)}
	*o.ptr = api.DecodeI32(r[0])
	runtime.SetFinalizer(o, func(o *RustBytes) {
		o.Free()
	})
	return o, nil
}

func (o *RustBytesPointer) Read() ([]byte, error) {
	if o.ptr == nil {
		return nil, ErrNullWasmPointer
	}
	err := o.binding()
	if err != nil {
		return nil, err
	}
	return o.o.Read()
}

func (o *RustBytesPointer) Write(data []byte) error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	err := o.binding()
	if err != nil {
		return err
	}
	return o.o.Write(data)
}

// Free Free
func (o *RustBytesPointer) Free() error {
	if o.ptr == nil {
		return ErrNullWasmPointer
	}
	err := o.binding()
	if err != nil {
		return err
	}
	err = o.o.Free()
	if err == nil {
		return nil
	}
	_, err = funcRustByteFree.Call(
		ctx, api.EncodeI32(*o.ptr), 8)
	if err != nil {
		return err
	}
	o.ptr = nil
	return nil
}

func (o *RustBytesPointer) binding() error {
	if o.o != nil {
		return nil
	}
	ob, f := wasi.Memory().Read(uint32(*o.ptr), 4)
	if !f {
		return ErrOutOfRange
	}
	offset := binary.LittleEndian.Uint32(ob)
	lb, f := wasi.Memory().Read(uint32(*o.ptr)+4, 4)
	if !f {
		return ErrOutOfRange
	}
	length := binary.LittleEndian.Uint32(lb)
	o.o = &RustBytes{
		ptr: new(int32),
		len: int32(length),
	}
	*o.o.ptr = int32(offset)
	return nil
}
