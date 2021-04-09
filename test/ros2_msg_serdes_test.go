package test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	sensor_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/sensor_msgs/msg"
	std_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/std_msgs/msg"
	test_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/test_msgs/msg"
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

/*
ROS2 test_msgs -package has test messages for all the ways messages can be defined.
*/
func TestSerDesROS2Messages_test_msgs(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)

	Convey("test_msgs.Arrays", t, func() {
		goObj := &test_msgs.Arrays{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__Arrays()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Arrays())
		So((*_Ctype_struct_test_msgs__msg__Arrays)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__Arrays())
	})
	Convey("test_msgs.BasicTypes", t, func() {
		goObj := &test_msgs.BasicTypes{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__BasicTypes()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__BasicTypes())
		So((*_Ctype_struct_test_msgs__msg__BasicTypes)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__BasicTypes())
	})
	Convey("test_msgs.Builtins", t, func() {
		goObj := &test_msgs.Builtins{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__Builtins()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Builtins())
		So((*_Ctype_struct_test_msgs__msg__Builtins)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__Builtins())
	})
	Convey("test_msgs.BoundedSequences", t, func() {
		goObj := &test_msgs.BoundedSequences{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__BoundedSequences()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__BoundedSequences())
		So((*_Ctype_struct_test_msgs__msg__BoundedSequences)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__BoundedSequences())
	})
	Convey("test_msgs.Constants", t, func() {
		goObj := &test_msgs.Constants{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__Constants()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Constants())
		So((*_Ctype_struct_test_msgs__msg__Constants)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__Constants())
	})
	Convey("test_msgs.Defaults", t, func() {
		goObj := &test_msgs.Defaults{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__Defaults()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Defaults())
		So((*_Ctype_struct_test_msgs__msg__Defaults)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__Defaults())
	})
	Convey("test_msgs.Empty", t, func() {
		goObj := &test_msgs.Empty{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__Empty()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Empty())
		So((*_Ctype_struct_test_msgs__msg__Empty)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__Empty())
	})
	Convey("test_msgs.MultiNested", t, func() {
		goObj := &test_msgs.MultiNested{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__MultiNested()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__MultiNested())
		So((*_Ctype_struct_test_msgs__msg__MultiNested)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__MultiNested())
	})
	Convey("test_msgs.Nested", t, func() {
		goObj := &test_msgs.Nested{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__Nested()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Nested())
		So((*_Ctype_struct_test_msgs__msg__Nested)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__Nested())
	})
	Convey("test_msgs.UnboundedSequences do not allocate memory for empty slices", t, func() {
		goObj := &test_msgs.UnboundedSequences{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__UnboundedSequences_no_allocate_memory_on_empty_slice()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__UnboundedSequences_no_allocate_memory_on_empty_slice())
		So((*_Ctype_struct_test_msgs__msg__UnboundedSequences)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__UnboundedSequences_no_allocate_memory_on_empty_slice())
	})
	Convey("test_msgs.UnboundedSequences", t, func() {
		goObj := &test_msgs.UnboundedSequences{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__UnboundedSequences()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__UnboundedSequences())
		So((*_Ctype_struct_test_msgs__msg__UnboundedSequences)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__UnboundedSequences())
	})
	Convey("test_msgs.WStrings", t, func() {
		goObj := &test_msgs.WStrings{}
		goObj.AsGoStruct(unsafe.Pointer(Fixture_C_test_msgs__WStrings()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__WStrings())
		So((*_Ctype_struct_test_msgs__msg__WStrings)(goObj.AsCStruct()), ShouldResemble, Fixture_C_test_msgs__WStrings())
	})
}
