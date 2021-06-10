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
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lsensor_msgs__rosidl_typesupport_c -lsensor_msgs__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>

#include <sensor_msgs/msg/region_of_interest.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("sensor_msgs/RegionOfInterest", RegionOfInterestTypeSupport)
}

// Do not create instances of this type directly. Always use NewRegionOfInterest
// function instead.
type RegionOfInterest struct {
	XOffset uint32 `yaml:"x_offset"`// Leftmost pixel of the ROI
	YOffset uint32 `yaml:"y_offset"`// Topmost pixel of the ROI. (0 if the ROI includes the left edge of the image)
	Height uint32 `yaml:"height"`// Height of ROI. (0 if the ROI includes the left edge of the image)(0 if the ROI includes the top edge of the image)
	Width uint32 `yaml:"width"`// Width of ROI. (0 if the ROI includes the left edge of the image)(0 if the ROI includes the top edge of the image)
	DoRectify bool `yaml:"do_rectify"`// True if a distinct rectified ROI should be calculated from the "raw"ROI in this message. Typically this should be False if the full imageis captured (ROI not used), and True if a subwindow is captured (ROIused).
}

// NewRegionOfInterest creates a new RegionOfInterest with default values.
func NewRegionOfInterest() *RegionOfInterest {
	self := RegionOfInterest{}
	self.SetDefaults()
	return &self
}

func (t *RegionOfInterest) Clone() *RegionOfInterest {
	c := &RegionOfInterest{}
	c.XOffset = t.XOffset
	c.YOffset = t.YOffset
	c.Height = t.Height
	c.Width = t.Width
	c.DoRectify = t.DoRectify
	return c
}

func (t *RegionOfInterest) CloneMsg() types.Message {
	return t.Clone()
}

func (t *RegionOfInterest) SetDefaults() {
	
}

// CloneRegionOfInterestSlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func CloneRegionOfInterestSlice(dst, src []RegionOfInterest) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var RegionOfInterestTypeSupport types.MessageTypeSupport = _RegionOfInterestTypeSupport{}

type _RegionOfInterestTypeSupport struct{}

func (t _RegionOfInterestTypeSupport) New() types.Message {
	return NewRegionOfInterest()
}

func (t _RegionOfInterestTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.sensor_msgs__msg__RegionOfInterest
	return (unsafe.Pointer)(C.sensor_msgs__msg__RegionOfInterest__create())
}

func (t _RegionOfInterestTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.sensor_msgs__msg__RegionOfInterest__destroy((*C.sensor_msgs__msg__RegionOfInterest)(pointer_to_free))
}

func (t _RegionOfInterestTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	m := msg.(*RegionOfInterest)
	mem := (*C.sensor_msgs__msg__RegionOfInterest)(dst)
	mem.x_offset = C.uint32_t(m.XOffset)
	mem.y_offset = C.uint32_t(m.YOffset)
	mem.height = C.uint32_t(m.Height)
	mem.width = C.uint32_t(m.Width)
	mem.do_rectify = C.bool(m.DoRectify)
}

func (t _RegionOfInterestTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	m := msg.(*RegionOfInterest)
	mem := (*C.sensor_msgs__msg__RegionOfInterest)(ros2_message_buffer)
	m.XOffset = uint32(mem.x_offset)
	m.YOffset = uint32(mem.y_offset)
	m.Height = uint32(mem.height)
	m.Width = uint32(mem.width)
	m.DoRectify = bool(mem.do_rectify)
}

func (t _RegionOfInterestTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__sensor_msgs__msg__RegionOfInterest())
}

type CRegionOfInterest = C.sensor_msgs__msg__RegionOfInterest
type CRegionOfInterest__Sequence = C.sensor_msgs__msg__RegionOfInterest__Sequence

func RegionOfInterest__Sequence_to_Go(goSlice *[]RegionOfInterest, cSlice CRegionOfInterest__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]RegionOfInterest, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.sensor_msgs__msg__RegionOfInterest__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_sensor_msgs__msg__RegionOfInterest * uintptr(i)),
		))
		RegionOfInterestTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(cIdx))
	}
}
func RegionOfInterest__Sequence_to_C(cSlice *CRegionOfInterest__Sequence, goSlice []RegionOfInterest) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.sensor_msgs__msg__RegionOfInterest)(C.malloc((C.size_t)(C.sizeof_struct_sensor_msgs__msg__RegionOfInterest * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.sensor_msgs__msg__RegionOfInterest)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_sensor_msgs__msg__RegionOfInterest * uintptr(i)),
		))
		RegionOfInterestTypeSupport.AsCStruct(unsafe.Pointer(cIdx), &v)
	}
}
func RegionOfInterest__Array_to_Go(goSlice []RegionOfInterest, cSlice []CRegionOfInterest) {
	for i := 0; i < len(cSlice); i++ {
		RegionOfInterestTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func RegionOfInterest__Array_to_C(cSlice []CRegionOfInterest, goSlice []RegionOfInterest) {
	for i := 0; i < len(goSlice); i++ {
		RegionOfInterestTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}