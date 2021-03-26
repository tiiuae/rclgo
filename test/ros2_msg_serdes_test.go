package test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	sensor_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/sensor_msgs/msg"
	std_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/std_msgs/msg"
)

func TestSerDesROS2Messages(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)

	Convey("std_msgs.ColorRGBA", t, func() {
		goObj := &std_msgs.ColorRGBA{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_std_msgs__ColorRGBA()))
		So(goObj, ShouldResemble, Fixture_Go_std_msgs__ColorRGBA())
		So((*_Ctype_struct_std_msgs__msg__ColorRGBA)(goObj.AsCStruct()), ShouldResemble, Fixture_C_std_msgs__ColorRGBA())
	})
	Convey("std_msgs.String", t, func() {
		goObj := &std_msgs.String{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_std_msgs__String()))
		So(goObj, ShouldResemble, Fixture_Go_std_msgs__String())
		So((*_Ctype_struct_std_msgs__msg__String)(goObj.AsCStruct()), ShouldResemble, Fixture_C_std_msgs__String())
	})
	Convey("sensor_msgs.ChannelFloat32", t, func() {
		goObj := &sensor_msgs.ChannelFloat32{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_sensor_msgs__ChannelFloat32()))
		So(goObj, ShouldResemble, Fixture_Go_sensor_msgs__ChannelFloat32())
		So((*_Ctype_struct_sensor_msgs__msg__ChannelFloat32)(goObj.AsCStruct()), ShouldResemble, Fixture_C_sensor_msgs__ChannelFloat32())
	})
	Convey("sensor_msgs.Illuminance", t, func() {
		goObj := &sensor_msgs.Illuminance{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_sensor_msgs__Illuminance()))
		So(goObj, ShouldResemble, Fixture_Go_sensor_msgs__Illuminance())
		So((*_Ctype_struct_sensor_msgs__msg__Illuminance)(goObj.AsCStruct()), ShouldResemble, Fixture_C_sensor_msgs__Illuminance())
	})
	Convey("std_msgs.Int64MultiArray", t, func() {
		goObj := &std_msgs.Int64MultiArray{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_std_msgs__Int64MultiArray()))
		So(goObj, ShouldResemble, Fixture_Go_std_msgs__Int64MultiArray())
		So((*_Ctype_struct_std_msgs__msg__Int64MultiArray)(goObj.AsCStruct()), ShouldResemble, Fixture_C_std_msgs__Int64MultiArray())
	})
}
