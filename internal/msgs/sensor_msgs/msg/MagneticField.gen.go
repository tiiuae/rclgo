/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

/*
THIS FILE IS AUTOGENERATED BY 'rclgo-gen generate'
*/

package sensor_msgs_msg
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo/types"
	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	geometry_msgs_msg "github.com/tiiuae/rclgo/internal/msgs/geometry_msgs/msg"
	std_msgs_msg "github.com/tiiuae/rclgo/internal/msgs/std_msgs/msg"
	primitives "github.com/tiiuae/rclgo/pkg/rclgo/primitives"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/galactic/lib -Wl,-rpath=/opt/ros/galactic/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lsensor_msgs__rosidl_typesupport_c -lsensor_msgs__rosidl_generator_c
#cgo LDFLAGS: -lgeometry_msgs__rosidl_typesupport_c -lgeometry_msgs__rosidl_generator_c
#cgo LDFLAGS: -lstd_msgs__rosidl_typesupport_c -lstd_msgs__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/galactic/include

#include <rosidl_runtime_c/message_type_support_struct.h>

#include <sensor_msgs/msg/magnetic_field.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("sensor_msgs/MagneticField", MagneticFieldTypeSupport)
}

// Do not create instances of this type directly. Always use NewMagneticField
// function instead.
type MagneticField struct {
	Header std_msgs_msg.Header `yaml:"header"`// timestamp is the time the
	MagneticField geometry_msgs_msg.Vector3 `yaml:"magnetic_field"`// x, y, and z components of the
	MagneticFieldCovariance [9]float64 `yaml:"magnetic_field_covariance"`// Row major about x, y, z axes
}

// NewMagneticField creates a new MagneticField with default values.
func NewMagneticField() *MagneticField {
	self := MagneticField{}
	self.SetDefaults()
	return &self
}

func (t *MagneticField) Clone() *MagneticField {
	c := &MagneticField{}
	c.Header = *t.Header.Clone()
	c.MagneticField = *t.MagneticField.Clone()
	c.MagneticFieldCovariance = t.MagneticFieldCovariance
	return c
}

func (t *MagneticField) CloneMsg() types.Message {
	return t.Clone()
}

func (t *MagneticField) SetDefaults() {
	t.Header.SetDefaults()
	t.MagneticField.SetDefaults()
	t.MagneticFieldCovariance = [9]float64{}
}

// CloneMagneticFieldSlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func CloneMagneticFieldSlice(dst, src []MagneticField) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var MagneticFieldTypeSupport types.MessageTypeSupport = _MagneticFieldTypeSupport{}

type _MagneticFieldTypeSupport struct{}

func (t _MagneticFieldTypeSupport) New() types.Message {
	return NewMagneticField()
}

func (t _MagneticFieldTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.sensor_msgs__msg__MagneticField
	return (unsafe.Pointer)(C.sensor_msgs__msg__MagneticField__create())
}

func (t _MagneticFieldTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.sensor_msgs__msg__MagneticField__destroy((*C.sensor_msgs__msg__MagneticField)(pointer_to_free))
}

func (t _MagneticFieldTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	m := msg.(*MagneticField)
	mem := (*C.sensor_msgs__msg__MagneticField)(dst)
	std_msgs_msg.HeaderTypeSupport.AsCStruct(unsafe.Pointer(&mem.header), &m.Header)
	geometry_msgs_msg.Vector3TypeSupport.AsCStruct(unsafe.Pointer(&mem.magnetic_field), &m.MagneticField)
	cSlice_magnetic_field_covariance := mem.magnetic_field_covariance[:]
	primitives.Float64__Array_to_C(*(*[]primitives.CFloat64)(unsafe.Pointer(&cSlice_magnetic_field_covariance)), m.MagneticFieldCovariance[:])
}

func (t _MagneticFieldTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	m := msg.(*MagneticField)
	mem := (*C.sensor_msgs__msg__MagneticField)(ros2_message_buffer)
	std_msgs_msg.HeaderTypeSupport.AsGoStruct(&m.Header, unsafe.Pointer(&mem.header))
	geometry_msgs_msg.Vector3TypeSupport.AsGoStruct(&m.MagneticField, unsafe.Pointer(&mem.magnetic_field))
	cSlice_magnetic_field_covariance := mem.magnetic_field_covariance[:]
	primitives.Float64__Array_to_Go(m.MagneticFieldCovariance[:], *(*[]primitives.CFloat64)(unsafe.Pointer(&cSlice_magnetic_field_covariance)))
}

func (t _MagneticFieldTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__sensor_msgs__msg__MagneticField())
}

type CMagneticField = C.sensor_msgs__msg__MagneticField
type CMagneticField__Sequence = C.sensor_msgs__msg__MagneticField__Sequence

func MagneticField__Sequence_to_Go(goSlice *[]MagneticField, cSlice CMagneticField__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]MagneticField, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.sensor_msgs__msg__MagneticField__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_sensor_msgs__msg__MagneticField * uintptr(i)),
		))
		MagneticFieldTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(cIdx))
	}
}
func MagneticField__Sequence_to_C(cSlice *CMagneticField__Sequence, goSlice []MagneticField) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.sensor_msgs__msg__MagneticField)(C.malloc((C.size_t)(C.sizeof_struct_sensor_msgs__msg__MagneticField * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.sensor_msgs__msg__MagneticField)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_sensor_msgs__msg__MagneticField * uintptr(i)),
		))
		MagneticFieldTypeSupport.AsCStruct(unsafe.Pointer(cIdx), &v)
	}
}
func MagneticField__Array_to_Go(goSlice []MagneticField, cSlice []CMagneticField) {
	for i := 0; i < len(cSlice); i++ {
		MagneticFieldTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func MagneticField__Array_to_C(cSlice []CMagneticField, goSlice []MagneticField) {
	for i := 0; i < len(goSlice); i++ {
		MagneticFieldTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}
