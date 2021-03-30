/*
ROS2 has helpfully included test_msgs/msg/* which seem to cover all cases of .msg-file serdes.
Here we define test fixtures for each message type for rather complete test coverage to the .msg parsing functionality.
*/
package test

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation -lrcl_interfaces__rosidl_typesupport_c
#cgo LDFLAGS: -ltest_msgs__rosidl_typesupport_c -ltest_msgs__rosidl_generator_c
#cgo CFLAGS: -I/opt/ros/foxy/include

#include <float.h>
#include <limits.h>

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <rosidl_runtime_c/primitives_sequence_functions.h>
#include <rosidl_runtime_c/string_functions.h>
#include <test_msgs/msg/arrays.h>
#include <test_msgs/msg/basic_types.h>
#include <test_msgs/msg/builtins.h>
#include <test_msgs/msg/constants.h>
#include <test_msgs/msg/defaults.h>
#include <test_msgs/msg/empty.h>
#include <test_msgs/msg/nested.h>
#include <test_msgs/msg/unbounded_sequences.h>

float float32_max() {
	return FLT_MAX;
}
float float32_min() {
	return FLT_MIN;
}
double float64_max() {
	return DBL_MAX;
}
double float64_min() {
	return DBL_MIN;
}
uint64_t uint64_max() {
	return ULONG_MAX;
}

test_msgs__msg__Arrays * test_msgs__msg__Arrays__fixture() {
	test_msgs__msg__Arrays * obj = test_msgs__msg__Arrays__create();
	obj->bool_values[0] = false;
	obj->bool_values[1] = true;
	obj->bool_values[2] = false;
	obj->byte_values[0] = 32;
	obj->byte_values[1] = 64;
	obj->byte_values[2] = 128;
	obj->char_values[0] = 'c';
	obj->char_values[1] = 'b';
	obj->char_values[2] = 'd';
	obj->float32_values[0] = float32_max();
	obj->float32_values[1] = float32_min();
	obj->float32_values[2] = 10.01;
	obj->float64_values[0] = float64_max();
	obj->float64_values[1] = float64_min();
	obj->float64_values[2] = 1010.0101;
	obj->int8_values[3];
	obj->uint8_values[3];
	obj->int16_values[3];
	obj->uint16_values[3];
	obj->int32_values[3];
	obj->uint32_values[3];
	obj->int64_values[3];
	obj->uint64_values[0] = uint64_max();
	obj->uint64_values[1] = 0;
	obj->uint64_values[2] = 3333333333;
	obj->string_values[3];
	obj->basic_types_values[3];
	obj->constants_values[3];
	obj->defaults_values[3];
	obj->bool_values_default[3];
	obj->byte_values_default[3];
	obj->char_values_default[3];
	obj->float32_values_default[3];
	obj->float64_values_default[3];
	obj->int8_values_default[3];
	obj->uint8_values_default[3];
	obj->int16_values_default[3];
	obj->uint16_values_default[3];
	obj->int32_values_default[3];
	obj->uint32_values_default[3];
	obj->int64_values_default[3];
	obj->uint64_values_default[3];
	obj->string_values_default[3];
	obj->alignment_check;
	return obj;
}

test_msgs__msg__BasicTypes * test_msgs__msg__BasicTypes__fixture() {
	test_msgs__msg__BasicTypes * obj = test_msgs__msg__BasicTypes__create();
	obj->bool_value = true;
	obj->byte_value = 32;
	obj->char_value = 'x';
	obj->float32_value = 5544.4455;
	obj->float64_value = 665544.445566;
	obj->int64_value = -1111111111;
	obj->uint64_value = 222222222222;
	return obj;
}

test_msgs__msg__Builtins * test_msgs__msg__Builtins__fixture() {
	test_msgs__msg__Builtins * obj = test_msgs__msg__Builtins__create();
	obj->duration_value.sec = 10;
	obj->duration_value.nanosec = 10001;
	obj->time_value.sec = 20;
	obj->time_value.nanosec = 20002;
	return obj;
}

test_msgs__msg__Constants * test_msgs__msg__Constants__fixture() {
	test_msgs__msg__Constants * obj = test_msgs__msg__Constants__create();
	return obj;
}

test_msgs__msg__Defaults * test_msgs__msg__Defaults__fixture() {
	test_msgs__msg__Defaults * obj = test_msgs__msg__Defaults__create();
	return obj;
}

test_msgs__msg__Empty * test_msgs__msg__Empty__fixture() {
	test_msgs__msg__Empty * obj = test_msgs__msg__Empty__create();
	return obj;
}

test_msgs__msg__Nested * test_msgs__msg__Nested__fixture() {
	test_msgs__msg__Nested * obj = test_msgs__msg__Nested__create();
	obj->basic_types_value.int8_value = 16;
	return obj;
}

test_msgs__msg__UnboundedSequences * test_msgs__msg__UnboundedSequences__fixture() {
	test_msgs__msg__UnboundedSequences * obj = test_msgs__msg__UnboundedSequences__create();
	rosidl_runtime_c__uint8__Sequence__init(&obj->uint8_values, 2);
	obj->uint8_values.data[0] = 12;
	obj->uint8_values.data[1] = 24;
	rosidl_runtime_c__uint8__Sequence__init(&obj->char_values, 3);
	obj->char_values.data[0] = 'o';
	obj->char_values.data[1] = 'm';
	obj->char_values.data[2] = 'g';
	rosidl_runtime_c__int64__Sequence__init(&obj->int64_values, 1);
	obj->int64_values.data[0] = 32;
	rosidl_runtime_c__bool__Sequence__init(&obj->bool_values, 0);

	return obj;
}

*/
import "C"
import test_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/test_msgs/msg"

func Fixture_C_test_msgs__Arrays() *C.test_msgs__msg__Arrays {
	return C.test_msgs__msg__Arrays__fixture()
}
func Fixture_Go_test_msgs__Arrays() *test_msgs.Arrays {
	obj := test_msgs.NewArrays()
	obj.BoolValues = [3]bool{false, true, false}
	obj.ByteValues = [3]byte{32, 64, 128}
	obj.CharValues = [3]byte{'c', 'b', 'd'}
	obj.Float32Values = [3]float32{float32(C.float32_max()), float32(C.float32_min()), 10.01}
	obj.Float64Values = [3]float64{float64(C.float64_max()), float64(C.float64_min()), 1010.0101}
	obj.Uint64Values = [3]uint64{uint64(C.uint64_max()), 0, 3333333333}
	return obj
}

func Fixture_C_test_msgs__BasicTypes() *C.test_msgs__msg__BasicTypes {
	return C.test_msgs__msg__BasicTypes__fixture()
}
func Fixture_Go_test_msgs__BasicTypes() *test_msgs.BasicTypes {
	obj := test_msgs.NewBasicTypes()
	obj.BoolValue = true
	obj.ByteValue = 32
	obj.CharValue = 'x'
	obj.Float32Value = 5544.4455
	obj.Float64Value = 665544.445566
	obj.Int64Value = -1111111111
	obj.Uint64Value = 222222222222
	return obj
}

func Fixture_C_test_msgs__Builtins() *C.test_msgs__msg__Builtins {
	return C.test_msgs__msg__Builtins__fixture()
}
func Fixture_Go_test_msgs__Builtins() *test_msgs.Builtins {
	obj := test_msgs.NewBuiltins()
	obj.DurationValue.Sec = 10
	obj.DurationValue.Nanosec = 10001
	obj.TimeValue.Sec = 20
	obj.TimeValue.Nanosec = 20002
	return obj
}

func Fixture_C_test_msgs__Constants() *C.test_msgs__msg__Constants {
	return C.test_msgs__msg__Constants__fixture()
}
func Fixture_Go_test_msgs__Constants() *test_msgs.Constants {
	return &test_msgs.Constants{}
}

func Fixture_C_test_msgs__Defaults() *C.test_msgs__msg__Defaults {
	return C.test_msgs__msg__Defaults__fixture()
}
func Fixture_Go_test_msgs__Defaults() *test_msgs.Defaults {
	return test_msgs.NewDefaults()
}

func Fixture_C_test_msgs__Empty() *C.test_msgs__msg__Empty {
	return C.test_msgs__msg__Empty__fixture()
}
func Fixture_Go_test_msgs__Empty() *test_msgs.Empty {
	return &test_msgs.Empty{}
}

func Fixture_C_test_msgs__Nested() *C.test_msgs__msg__Nested {
	return C.test_msgs__msg__Nested__fixture()
}
func Fixture_Go_test_msgs__Nested() *test_msgs.Nested {
	obj := test_msgs.NewNested()
	obj.BasicTypesValue.Int8Value = 16
	return obj
}

func Fixture_C_test_msgs__UnboundedSequences_no_allocate_memory_on_empty_slice() *C.test_msgs__msg__UnboundedSequences {
	return C.test_msgs__msg__UnboundedSequences__create()
}
func Fixture_Go_test_msgs__UnboundedSequences_no_allocate_memory_on_empty_slice() *test_msgs.UnboundedSequences {
	obj := test_msgs.NewUnboundedSequences()
	return obj
}

func Fixture_C_test_msgs__UnboundedSequences() *C.test_msgs__msg__UnboundedSequences {
	return C.test_msgs__msg__UnboundedSequences__fixture()
}
func Fixture_Go_test_msgs__UnboundedSequences() *test_msgs.UnboundedSequences {
	obj := test_msgs.NewUnboundedSequences()
	obj.Uint8Values = make([]uint8, 2)
	obj.Uint8Values[0] = 12
	obj.Uint8Values[1] = 24
	obj.CharValues = make([]byte, 3)
	obj.CharValues[0] = 'o'
	obj.CharValues[1] = 'm'
	obj.CharValues[2] = 'g'
	obj.Int64Values = make([]int64, 1)
	obj.Int64Values[0] = 32
	//obj.BoolValues = make([]bool, 0) // Even if we instantiate a zero-length rosidl-Sequence in C, we get a NULL pointer
	return obj
}
