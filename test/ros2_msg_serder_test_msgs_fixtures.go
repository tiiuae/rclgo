/*
ROS2 has helpfully included test_msgs/msg/* which seem to cover all cases of .msg-file serdes.
Here we define test fixtures for each message type for rather complete test coverage to the .msg parsing functionality.
*/
package test

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation -lrcl_interfaces__rosidl_typesupport_c
#cgo LDFLAGS: -ltest_msgs__rosidl_typesupport_c -ltest_msgs__rosidl_generator_c
#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <rosidl_runtime_c/primitives_sequence_functions.h>
#include <rosidl_runtime_c/string_functions.h>
#include <test_msgs/msg/arrays.h>

test_msgs__msg__Arrays * test_msgs__msg__Arrays__fixture() {
	test_msgs__msg__Arrays * obj = test_msgs__msg__Arrays__create();
	obj->bool_values[0] = false;
	obj->bool_values[1] = true;
	obj->bool_values[2] = true;
	obj->byte_values[3];
	obj->char_values[3];
	obj->float32_values[3];
	obj->float64_values[3];
	obj->int8_values[3];
	obj->uint8_values[3];
	obj->int16_values[3];
	obj->uint16_values[3];
	obj->int32_values[3];
	obj->uint32_values[3];
	obj->int64_values[3];
	obj->uint64_values[3];
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

*/
import "C"
import test_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/test_msgs/msg"

func Fixture_C_test_msgs__Arrays() *C.test_msgs__msg__Arrays {
	return C.test_msgs__msg__Arrays__fixture()
}
func Fixture_Go_test_msgs__Arrays() *test_msgs.Arrays {
	return &test_msgs.Arrays{}
}
