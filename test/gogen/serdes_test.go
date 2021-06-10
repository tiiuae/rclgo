package test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	sensor_msgs "github.com/tiiuae/rclgo/internal/msgs/sensor_msgs/msg"
	std_msgs "github.com/tiiuae/rclgo/internal/msgs/std_msgs/msg"
	test_msgs "github.com/tiiuae/rclgo/internal/msgs/test_msgs/msg"
)

func TestSerDesROS2Messages(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)

	Convey("std_msgs.ColorRGBA", t, func() {
		goObj := std_msgs.NewColorRGBA()
		std_msgs.ColorRGBATypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_std_msgs__ColorRGBA()))
		So(goObj, ShouldResemble, Fixture_Go_std_msgs__ColorRGBA())
		cobj := std_msgs.ColorRGBATypeSupport.PrepareMemory()
		defer std_msgs.ColorRGBATypeSupport.ReleaseMemory(cobj)
		std_msgs.ColorRGBATypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_std_msgs__msg__ColorRGBA)(cobj), ShouldResemble, Fixture_C_std_msgs__ColorRGBA())
	})
	Convey("std_msgs.String", t, func() {
		goObj := std_msgs.NewString()
		std_msgs.StringTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_std_msgs__String()))
		So(goObj, ShouldResemble, Fixture_Go_std_msgs__String())
		cobj := std_msgs.StringTypeSupport.PrepareMemory()
		defer std_msgs.StringTypeSupport.ReleaseMemory(cobj)
		std_msgs.StringTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_std_msgs__msg__String)(cobj), ShouldResemble, Fixture_C_std_msgs__String())
	})
	Convey("sensor_msgs.ChannelFloat32", t, func() {
		goObj := sensor_msgs.NewChannelFloat32()
		sensor_msgs.ChannelFloat32TypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_sensor_msgs__ChannelFloat32()))
		So(goObj, ShouldResemble, Fixture_Go_sensor_msgs__ChannelFloat32())
		cobj := sensor_msgs.ChannelFloat32TypeSupport.PrepareMemory()
		defer sensor_msgs.ChannelFloat32TypeSupport.ReleaseMemory(cobj)
		sensor_msgs.ChannelFloat32TypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_sensor_msgs__msg__ChannelFloat32)(cobj), ShouldResemble, Fixture_C_sensor_msgs__ChannelFloat32())
	})
	Convey("sensor_msgs.Illuminance", t, func() {
		goObj := sensor_msgs.NewIlluminance()
		sensor_msgs.IlluminanceTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_sensor_msgs__Illuminance()))
		So(goObj, ShouldResemble, Fixture_Go_sensor_msgs__Illuminance())
		cobj := sensor_msgs.IlluminanceTypeSupport.PrepareMemory()
		defer sensor_msgs.IlluminanceTypeSupport.ReleaseMemory(cobj)
		sensor_msgs.IlluminanceTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_sensor_msgs__msg__Illuminance)(cobj), ShouldResemble, Fixture_C_sensor_msgs__Illuminance())
	})
	Convey("std_msgs.Int64MultiArray", t, func() {
		goObj := std_msgs.NewInt64MultiArray()
		std_msgs.Int64MultiArrayTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_std_msgs__Int64MultiArray()))
		So(goObj, ShouldResemble, Fixture_Go_std_msgs__Int64MultiArray())
		cobj := std_msgs.Int64MultiArrayTypeSupport.PrepareMemory()
		defer std_msgs.Int64MultiArrayTypeSupport.ReleaseMemory(cobj)
		std_msgs.Int64MultiArrayTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_std_msgs__msg__Int64MultiArray)(cobj), ShouldResemble, Fixture_C_std_msgs__Int64MultiArray())
	})
}

/*
ROS2 test_msgs -package has test messages for all the ways messages can be defined.
*/
func TestSerDesROS2Messages_test_msgs(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)

	Convey("test_msgs.Arrays", t, func() {
		goObj := test_msgs.NewArrays()
		test_msgs.ArraysTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__Arrays()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Arrays())
		cobj := test_msgs.ArraysTypeSupport.PrepareMemory()
		defer test_msgs.ArraysTypeSupport.ReleaseMemory(cobj)
		test_msgs.ArraysTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__Arrays)(cobj), ShouldResemble, Fixture_C_test_msgs__Arrays())
	})
	Convey("test_msgs.BasicTypes", t, func() {
		goObj := test_msgs.NewBasicTypes()
		test_msgs.BasicTypesTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__BasicTypes()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__BasicTypes())
		cobj := test_msgs.BasicTypesTypeSupport.PrepareMemory()
		defer test_msgs.BasicTypesTypeSupport.ReleaseMemory(cobj)
		test_msgs.BasicTypesTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__BasicTypes)(cobj), ShouldResemble, Fixture_C_test_msgs__BasicTypes())
	})
	Convey("test_msgs.Builtins", t, func() {
		goObj := test_msgs.NewBuiltins()
		test_msgs.BuiltinsTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__Builtins()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Builtins())
		cobj := test_msgs.BuiltinsTypeSupport.PrepareMemory()
		defer test_msgs.BuiltinsTypeSupport.ReleaseMemory(cobj)
		test_msgs.BuiltinsTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__Builtins)(cobj), ShouldResemble, Fixture_C_test_msgs__Builtins())
	})
	Convey("test_msgs.BoundedSequences", t, func() {
		goObj := test_msgs.NewBoundedSequences()
		test_msgs.BoundedSequencesTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__BoundedSequences()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__BoundedSequences())
		cobj := test_msgs.BoundedSequencesTypeSupport.PrepareMemory()
		defer test_msgs.BoundedSequencesTypeSupport.ReleaseMemory(cobj)
		test_msgs.BoundedSequencesTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__BoundedSequences)(cobj), ShouldResemble, Fixture_C_test_msgs__BoundedSequences())
	})
	Convey("test_msgs.Constants", t, func() {
		goObj := test_msgs.NewConstants()
		test_msgs.ConstantsTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__Constants()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Constants())
		cobj := test_msgs.ConstantsTypeSupport.PrepareMemory()
		defer test_msgs.ConstantsTypeSupport.ReleaseMemory(cobj)
		test_msgs.ConstantsTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__Constants)(cobj), ShouldResemble, Fixture_C_test_msgs__Constants())
	})
	Convey("test_msgs.Defaults", t, func() {
		goObj := test_msgs.NewDefaults()
		test_msgs.DefaultsTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__Defaults()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Defaults())
		cobj := test_msgs.DefaultsTypeSupport.PrepareMemory()
		defer test_msgs.DefaultsTypeSupport.ReleaseMemory(cobj)
		test_msgs.DefaultsTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__Defaults)(cobj), ShouldResemble, Fixture_C_test_msgs__Defaults())
	})
	Convey("test_msgs.Empty", t, func() {
		goObj := test_msgs.NewEmpty()
		test_msgs.EmptyTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__Empty()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Empty())
		cobj := test_msgs.EmptyTypeSupport.PrepareMemory()
		defer test_msgs.EmptyTypeSupport.ReleaseMemory(cobj)
		test_msgs.EmptyTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__Empty)(cobj), ShouldResemble, Fixture_C_test_msgs__Empty())
	})
	Convey("test_msgs.MultiNested", t, func() {
		goObj := test_msgs.NewMultiNested()
		test_msgs.MultiNestedTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__MultiNested()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__MultiNested())
		cobj := test_msgs.MultiNestedTypeSupport.PrepareMemory()
		defer test_msgs.MultiNestedTypeSupport.ReleaseMemory(cobj)
		test_msgs.MultiNestedTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__MultiNested)(cobj), ShouldResemble, Fixture_C_test_msgs__MultiNested())
	})
	Convey("test_msgs.Nested", t, func() {
		goObj := test_msgs.NewNested()
		test_msgs.NestedTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__Nested()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__Nested())
		cobj := test_msgs.NestedTypeSupport.PrepareMemory()
		defer test_msgs.NestedTypeSupport.ReleaseMemory(cobj)
		test_msgs.NestedTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__Nested)(cobj), ShouldResemble, Fixture_C_test_msgs__Nested())
	})
	Convey("test_msgs.UnboundedSequences do not allocate memory for empty slices", t, func() {
		goObj := test_msgs.NewUnboundedSequences()
		test_msgs.UnboundedSequencesTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__UnboundedSequences_no_allocate_memory_on_empty_slice()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__UnboundedSequences_no_allocate_memory_on_empty_slice())
		cobj := test_msgs.UnboundedSequencesTypeSupport.PrepareMemory()
		defer test_msgs.UnboundedSequencesTypeSupport.ReleaseMemory(cobj)
		test_msgs.UnboundedSequencesTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__UnboundedSequences)(cobj), ShouldResemble, Fixture_C_test_msgs__UnboundedSequences_no_allocate_memory_on_empty_slice())
	})
	Convey("test_msgs.UnboundedSequences", t, func() {
		goObj := test_msgs.NewUnboundedSequences()
		test_msgs.UnboundedSequencesTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__UnboundedSequences()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__UnboundedSequences())
		cobj := test_msgs.UnboundedSequencesTypeSupport.PrepareMemory()
		defer test_msgs.UnboundedSequencesTypeSupport.ReleaseMemory(cobj)
		test_msgs.UnboundedSequencesTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__UnboundedSequences)(cobj), ShouldResemble, Fixture_C_test_msgs__UnboundedSequences())
	})
	Convey("test_msgs.WStrings", t, func() {
		goObj := test_msgs.NewWStrings()
		test_msgs.WStringsTypeSupport.AsGoStruct(goObj, unsafe.Pointer(Fixture_C_test_msgs__WStrings()))
		So(goObj, ShouldResemble, Fixture_Go_test_msgs__WStrings())
		cobj := test_msgs.WStringsTypeSupport.PrepareMemory()
		defer test_msgs.WStringsTypeSupport.ReleaseMemory(cobj)
		test_msgs.WStringsTypeSupport.AsCStruct(cobj, goObj)
		So((*_Ctype_struct_test_msgs__msg__WStrings)(cobj), ShouldResemble, Fixture_C_test_msgs__WStrings())
	})
}

func TestClone(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)

	Convey("Cloning works", t, func() {
		{
			goObj := Fixture_Go_std_msgs__ColorRGBA()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
			clone.B = 200
			So(clone, ShouldNotResemble, goObj)
		}
		{
			goObj := Fixture_Go_std_msgs__String()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
			clone.Data = "cloned string"
			So(clone, ShouldNotResemble, goObj)
		}
		{
			goObj := Fixture_Go_sensor_msgs__ChannelFloat32()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
			clone.Name = "cloned string"
			clone.Values[2] = 4.2
			So(clone, ShouldNotResemble, goObj)
		}
		{
			goObj := Fixture_Go_sensor_msgs__Illuminance()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
			clone.Header.Stamp.Nanosec = 1231234321
			So(clone, ShouldNotResemble, goObj)
		}
		{
			goObj := Fixture_Go_std_msgs__Int64MultiArray()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
			clone.Data[1] = 99
			So(clone, ShouldNotResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__Arrays()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
			clone.BasicTypesValues[2].Uint64Value = 4312
			So(clone, ShouldNotResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__BasicTypes()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__Builtins()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__BoundedSequences()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__Constants()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__Defaults()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__Empty()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__MultiNested()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__Nested()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__UnboundedSequences_no_allocate_memory_on_empty_slice()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__UnboundedSequences()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
		}
		{
			goObj := Fixture_Go_test_msgs__WStrings()
			clone := goObj.Clone()
			So(clone, ShouldResemble, goObj)
		}
	})
}
