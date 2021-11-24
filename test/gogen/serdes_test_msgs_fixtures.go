/*
ROS2 has helpfully included test_msgs/msg/* which seem to cover all cases of .msg-file serdes.
Here we define test fixtures for each message type for a rather complete test coverage to the .msg parsing functionality.
*/
package test

/*
#cgo LDFLAGS: -L/opt/ros/galactic/lib -Wl,-rpath=/opt/ros/galactic/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation -lrcl_interfaces__rosidl_typesupport_c
#cgo LDFLAGS: -ltest_msgs__rosidl_typesupport_c -ltest_msgs__rosidl_generator_c
#cgo CFLAGS: -I/opt/ros/galactic/include

#include <float.h>
#include <limits.h>

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <rosidl_runtime_c/primitives_sequence_functions.h>
#include <rosidl_runtime_c/string_functions.h>
#include <rosidl_runtime_c/u16string_functions.h>
#include <test_msgs/msg/arrays.h>
#include <test_msgs/msg/basic_types.h>
#include <test_msgs/msg/bounded_sequences.h>
#include <test_msgs/msg/builtins.h>
#include <test_msgs/msg/constants.h>
#include <test_msgs/msg/defaults.h>
#include <test_msgs/msg/empty.h>
#include <test_msgs/msg/multi_nested.h>
#include <test_msgs/msg/nested.h>
#include <test_msgs/msg/unbounded_sequences.h>
#include <test_msgs/msg/w_strings.h>

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

test_msgs__msg__BoundedSequences * test_msgs__msg__BoundedSequences__fixture() {
	test_msgs__msg__BoundedSequences * obj = test_msgs__msg__BoundedSequences__create();
	rosidl_runtime_c__bool__Sequence__init(&obj->bool_values, 5);
	obj->bool_values.data[0] = true;
	obj->bool_values.data[1] = false;
	obj->bool_values.data[2] = true; // Message definition declares bool_values as a bounded sequence of <=3.
	obj->bool_values.data[3] = true; // There doesn't seem to be any limit on ROS2's side on allocating more indexes than "allowed".
	obj->bool_values.data[4] = true; // We can leave this trap here in case such limits are enforced in the future.
	rosidl_runtime_c__String__Sequence__init(&obj->string_values, 4);
	rosidl_runtime_c__String__init(&obj->string_values.data[0]);
	rosidl_runtime_c__String__init(&obj->string_values.data[1]);
	rosidl_runtime_c__String__init(&obj->string_values.data[2]);
	rosidl_runtime_c__String__init(&obj->string_values.data[3]);
	rosidl_runtime_c__String__assignn(&obj->string_values.data[0], "Bared on your tomb\x00", 18);
	rosidl_runtime_c__String__assignn(&obj->string_values.data[1], "I'm a prayer for your loneliness\x00", 32);
	rosidl_runtime_c__String__assignn(&obj->string_values.data[2], "And would you ever soon\x00", 23);
	rosidl_runtime_c__String__assignn(&obj->string_values.data[3], "Come above unto me?\x00", 19);
	test_msgs__msg__BasicTypes__Sequence__init(&obj->basic_types_values, 1);
	obj->basic_types_values.data[0].byte_value = 12;
	obj->byte_values_default.data[0] = 10; // Override defaults here
	obj->byte_values_default.data[1] = 11;
	obj->byte_values_default.data[2] = 12;
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

test_msgs__msg__MultiNested * test_msgs__msg__MultiNested__fixture() {
	test_msgs__msg__MultiNested * obj = test_msgs__msg__MultiNested__create();
	obj->array_of_arrays[0].bool_values[0] = true;
	obj->array_of_arrays[0].bool_values[1] = false;
	obj->array_of_arrays[0].bool_values[2] = true;
	obj->array_of_arrays[1].basic_types_values[0].float32_value = 3.14159;
	obj->array_of_arrays[1].basic_types_values[1].float32_value;
	obj->array_of_arrays[1].basic_types_values[2].float32_value = 1.61803;
	obj->array_of_arrays[2].int8_values_default[1] = 64;
	obj->array_of_arrays[2].alignment_check = 32;

	test_msgs__msg__BoundedSequences__Sequence__init(&obj->bounded_sequence_of_bounded_sequences, 1);
	rosidl_runtime_c__bool__Sequence__init(&obj->bounded_sequence_of_bounded_sequences.data[0].bool_values, 1);
	obj->bounded_sequence_of_bounded_sequences.data[0].bool_values.data[0] = true;
	test_msgs__msg__BasicTypes__Sequence__init(&obj->bounded_sequence_of_bounded_sequences.data[0].basic_types_values, 1);
	obj->bounded_sequence_of_bounded_sequences.data[0].basic_types_values.data[0].byte_value = 12;
	obj->bounded_sequence_of_bounded_sequences.data[0].alignment_check = 32;

	test_msgs__msg__UnboundedSequences__Sequence__init(&obj->unbounded_sequence_of_unbounded_sequences, 1);
	rosidl_runtime_c__uint8__Sequence__init(&obj->unbounded_sequence_of_unbounded_sequences.data[0].char_values, 1);
	obj->unbounded_sequence_of_unbounded_sequences.data[0].char_values.data[0] = 'g';
	obj->unbounded_sequence_of_unbounded_sequences.data[0].alignment_check = 32;
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

test_msgs__msg__WStrings * test_msgs__msg__WStrings__fixture() {
	uint16_t* str = malloc(sizeof(uint16_t)*9);
	str[0] = 0x26A1; str[1] = 0x0048; str[2] = 0x2708; str[3] = 0x0048; str[4] = 0x2622; str[5] = 0x0048; str[6] = 0x2650; str[7] = 0x0048; str[8] = 0x0000;
	//str[0] = 0x0048; str[1] = 0x0048; str[2] = 0x0048; str[3] = 0x0048; str[4] = 0x0048; str[5] = 0x0048; str[6] = 0x0048; str[7] = 0x0048; str[8] = 0x0000;
	test_msgs__msg__WStrings * obj = test_msgs__msg__WStrings__create();
	rosidl_runtime_c__U16String__init(&obj->wstring_value);
	rosidl_runtime_c__U16String__assignn(&obj->wstring_value, str, 8);
	rosidl_runtime_c__U16String__init(&obj->array_of_wstrings[0]);
	rosidl_runtime_c__U16String__assignn(&obj->array_of_wstrings[0], str, 8);

	rosidl_runtime_c__U16String__Sequence__init(&obj->bounded_sequence_of_wstrings, 1);
	rosidl_runtime_c__U16String__init(&obj->bounded_sequence_of_wstrings.data[0]);
	rosidl_runtime_c__U16String__assignn(&obj->bounded_sequence_of_wstrings.data[0], str, 8);

	rosidl_runtime_c__U16String__Sequence__init(&obj->unbounded_sequence_of_wstrings, 1);
	rosidl_runtime_c__U16String__init(&obj->unbounded_sequence_of_wstrings.data[0]);
	rosidl_runtime_c__U16String__assignn(&obj->unbounded_sequence_of_wstrings.data[0], str, 8);

	return obj;
}

*/
import "C"
import (
	test_msgs "github.com/tiiuae/rclgo/internal/msgs/test_msgs/msg"
)

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

func Fixture_C_test_msgs__BoundedSequences() *C.test_msgs__msg__BoundedSequences {
	return C.test_msgs__msg__BoundedSequences__fixture()
}
func Fixture_Go_test_msgs__BoundedSequences() *test_msgs.BoundedSequences {
	obj := test_msgs.NewBoundedSequences()
	obj.BoolValues = make([]bool, 5)
	obj.BoolValues[0] = true
	obj.BoolValues[1] = false
	obj.BoolValues[2] = true // Message definition declares bool_values as a bounded sequence of <=3.
	obj.BoolValues[3] = true // There doesn't seem to be any limit on ROS2's side on allocating more indexes than "allowed".
	obj.BoolValues[4] = true // We can leave this trap here in case such limits are enforced in the future.
	obj.StringValues = make([]string, 4)
	obj.StringValues[0] = "Bared on your tomb"
	obj.StringValues[1] = "I'm a prayer for your loneliness"
	obj.StringValues[2] = "And would you ever soon"
	obj.StringValues[3] = "Come above unto me?"
	obj.BasicTypesValues = make([]test_msgs.BasicTypes, 1)
	obj.BasicTypesValues[0].ByteValue = 12
	obj.ByteValuesDefault[0] = 10 // Override defaults here
	obj.ByteValuesDefault[1] = 11
	obj.ByteValuesDefault[2] = 12
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

func Fixture_C_test_msgs__MultiNested() *C.test_msgs__msg__MultiNested {
	return C.test_msgs__msg__MultiNested__fixture()
}
func Fixture_Go_test_msgs__MultiNested() *test_msgs.MultiNested {
	obj := test_msgs.NewMultiNested()
	obj.ArrayOfArrays[0].BoolValues[0] = true
	obj.ArrayOfArrays[0].BoolValues[1] = false
	obj.ArrayOfArrays[0].BoolValues[2] = true
	obj.ArrayOfArrays[1].BasicTypesValues[0].Float32Value = 3.14159
	//obj.ArrayOfArrays[1].BasicTypesValues[1].Float32Value
	obj.ArrayOfArrays[1].BasicTypesValues[2].Float32Value = 1.61803
	obj.ArrayOfArrays[2].Int8ValuesDefault[1] = 64
	obj.ArrayOfArrays[2].AlignmentCheck = 32

	obj.BoundedSequenceOfBoundedSequences = make([]test_msgs.BoundedSequences, 1)
	obj.BoundedSequenceOfBoundedSequences[0] = *test_msgs.NewBoundedSequences()
	obj.BoundedSequenceOfBoundedSequences[0].BoolValues = make([]bool, 1)
	obj.BoundedSequenceOfBoundedSequences[0].BoolValues[0] = true
	obj.BoundedSequenceOfBoundedSequences[0].BasicTypesValues = make([]test_msgs.BasicTypes, 1)
	obj.BoundedSequenceOfBoundedSequences[0].BasicTypesValues[0].ByteValue = 12
	obj.BoundedSequenceOfBoundedSequences[0].AlignmentCheck = 32

	obj.UnboundedSequenceOfUnboundedSequences = make([]test_msgs.UnboundedSequences, 1)
	obj.UnboundedSequenceOfUnboundedSequences[0] = *test_msgs.NewUnboundedSequences()
	obj.UnboundedSequenceOfUnboundedSequences[0].CharValues = make([]byte, 1)
	obj.UnboundedSequenceOfUnboundedSequences[0].CharValues[0] = 'g'
	obj.UnboundedSequenceOfUnboundedSequences[0].AlignmentCheck = 32

	return obj
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

func Fixture_C_test_msgs__WStrings() *C.test_msgs__msg__WStrings {
	return C.test_msgs__msg__WStrings__fixture()
}
func Fixture_Go_test_msgs__WStrings() *test_msgs.WStrings {
	//str := string("⚡H✈️H☢️H♐H") // This UTF16 string instantiator adds some strange Private Use Area bytes between emojis?
	strUTF16 := []rune{
		0x26A1, 0x0048, 0x2708, 0x0048, 0x2622, 0x0048, 0x2650, 0x0048, // So to make the equally rendering strings comparable, we need to do it the hard way.
	}
	str := string(strUTF16)
	obj := test_msgs.NewWStrings()
	obj.WstringValue = str
	obj.ArrayOfWstrings[0] = str
	obj.BoundedSequenceOfWstrings = make([]string, 1)
	obj.BoundedSequenceOfWstrings[0] = str
	obj.UnboundedSequenceOfWstrings = make([]string, 1)
	obj.UnboundedSequenceOfWstrings[0] = str
	return obj
}
