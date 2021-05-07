/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

/*
THIS FILE IS AUTOGENERATED BY 'rclgo-gen generate'
*/

package sensor_msgs
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"
	std_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/std_msgs/msg"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lsensor_msgs__rosidl_typesupport_c -lsensor_msgs__rosidl_generator_c
#cgo LDFLAGS: -lstd_msgs__rosidl_typesupport_c -lstd_msgs__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <sensor_msgs/msg/multi_echo_laser_scan.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("sensor_msgs/MultiEchoLaserScan", &MultiEchoLaserScan{})
}

// Do not create instances of this type directly. Always use NewMultiEchoLaserScan
// function instead.
type MultiEchoLaserScan struct {
	Header std_msgs.Header `yaml:"header"`// timestamp in the header is the acquisition time of
	AngleMin float32 `yaml:"angle_min"`// start angle of the scan [rad]
	AngleMax float32 `yaml:"angle_max"`// end angle of the scan [rad]
	AngleIncrement float32 `yaml:"angle_increment"`// angular distance between measurements [rad]
	TimeIncrement float32 `yaml:"time_increment"`// time between measurements [seconds] - if your scanner
	ScanTime float32 `yaml:"scan_time"`// time between scans [seconds]. is moving, this will be used in interpolating positionof 3d points
	RangeMin float32 `yaml:"range_min"`// minimum range value [m]
	RangeMax float32 `yaml:"range_max"`// maximum range value [m]
	Ranges []LaserEcho `yaml:"ranges"`// range data [m]
	Intensities []LaserEcho `yaml:"intensities"`// intensity data [device-specific units].  If your. (Note: NaNs, values < range_min or > range_max should be discarded)+Inf measurements are out of range-Inf measurements are too close to determine exact distance.
}

// NewMultiEchoLaserScan creates a new MultiEchoLaserScan with default values.
func NewMultiEchoLaserScan() *MultiEchoLaserScan {
	self := MultiEchoLaserScan{}
	self.SetDefaults(nil)
	return &self
}

func (t *MultiEchoLaserScan) SetDefaults(d interface{}) ros2types.ROS2Msg {
	t.Header.SetDefaults(nil)
	
	return t
}

func (t *MultiEchoLaserScan) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__sensor_msgs__msg__MultiEchoLaserScan())
}
func (t *MultiEchoLaserScan) PrepareMemory() unsafe.Pointer { //returns *C.sensor_msgs__msg__MultiEchoLaserScan
	return (unsafe.Pointer)(C.sensor_msgs__msg__MultiEchoLaserScan__create())
}
func (t *MultiEchoLaserScan) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.sensor_msgs__msg__MultiEchoLaserScan__destroy((*C.sensor_msgs__msg__MultiEchoLaserScan)(pointer_to_free))
}
func (t *MultiEchoLaserScan) AsCStruct() unsafe.Pointer {
	mem := (*C.sensor_msgs__msg__MultiEchoLaserScan)(t.PrepareMemory())
	mem.header = *(*C.std_msgs__msg__Header)(t.Header.AsCStruct())
	mem.angle_min = C.float(t.AngleMin)
	mem.angle_max = C.float(t.AngleMax)
	mem.angle_increment = C.float(t.AngleIncrement)
	mem.time_increment = C.float(t.TimeIncrement)
	mem.scan_time = C.float(t.ScanTime)
	mem.range_min = C.float(t.RangeMin)
	mem.range_max = C.float(t.RangeMax)
	LaserEcho__Sequence_to_C(&mem.ranges, t.Ranges)
	LaserEcho__Sequence_to_C(&mem.intensities, t.Intensities)
	return unsafe.Pointer(mem)
}
func (t *MultiEchoLaserScan) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.sensor_msgs__msg__MultiEchoLaserScan)(ros2_message_buffer)
	t.Header.AsGoStruct(unsafe.Pointer(&mem.header))
	t.AngleMin = float32(mem.angle_min)
	t.AngleMax = float32(mem.angle_max)
	t.AngleIncrement = float32(mem.angle_increment)
	t.TimeIncrement = float32(mem.time_increment)
	t.ScanTime = float32(mem.scan_time)
	t.RangeMin = float32(mem.range_min)
	t.RangeMax = float32(mem.range_max)
	LaserEcho__Sequence_to_Go(&t.Ranges, mem.ranges)
	LaserEcho__Sequence_to_Go(&t.Intensities, mem.intensities)
}
func (t *MultiEchoLaserScan) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CMultiEchoLaserScan = C.sensor_msgs__msg__MultiEchoLaserScan
type CMultiEchoLaserScan__Sequence = C.sensor_msgs__msg__MultiEchoLaserScan__Sequence

func MultiEchoLaserScan__Sequence_to_Go(goSlice *[]MultiEchoLaserScan, cSlice CMultiEchoLaserScan__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]MultiEchoLaserScan, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.sensor_msgs__msg__MultiEchoLaserScan__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_sensor_msgs__msg__MultiEchoLaserScan * uintptr(i)),
		))
		(*goSlice)[i] = MultiEchoLaserScan{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func MultiEchoLaserScan__Sequence_to_C(cSlice *CMultiEchoLaserScan__Sequence, goSlice []MultiEchoLaserScan) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.sensor_msgs__msg__MultiEchoLaserScan)(C.malloc((C.size_t)(C.sizeof_struct_sensor_msgs__msg__MultiEchoLaserScan * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.sensor_msgs__msg__MultiEchoLaserScan)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_sensor_msgs__msg__MultiEchoLaserScan * uintptr(i)),
		))
		*cIdx = *(*C.sensor_msgs__msg__MultiEchoLaserScan)(v.AsCStruct())
	}
}
func MultiEchoLaserScan__Array_to_Go(goSlice []MultiEchoLaserScan, cSlice []CMultiEchoLaserScan) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func MultiEchoLaserScan__Array_to_C(cSlice []CMultiEchoLaserScan, goSlice []MultiEchoLaserScan) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.sensor_msgs__msg__MultiEchoLaserScan)(goSlice[i].AsCStruct())
	}
}

