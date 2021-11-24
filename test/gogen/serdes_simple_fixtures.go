/*
Fixtures here represent some easier data structures to compile and test.
These are handy when you need to make bigger changes to the way the templates are generated and you can get some confirmation that things work.
*/
package test

/*
#cgo LDFLAGS: -L/opt/ros/galactic/lib -Wl,-rpath=/opt/ros/galactic/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation -lrcl_interfaces__rosidl_typesupport_c
#cgo LDFLAGS: -lsensor_msgs__rosidl_typesupport_c -lsensor_msgs__rosidl_generator_c
#cgo LDFLAGS: -lstd_msgs__rosidl_typesupport_c -lstd_msgs__rosidl_generator_c
#cgo CFLAGS: -I/opt/ros/galactic/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <rosidl_runtime_c/primitives_sequence_functions.h>
#include <rosidl_runtime_c/string_functions.h>
#include <sensor_msgs/msg/channel_float32.h>
#include <sensor_msgs/msg/illuminance.h>
#include <std_msgs/msg/color_rgba.h>
#include <std_msgs/msg/int64_multi_array.h>
#include <std_msgs/msg/multi_array_dimension.h>
#include <std_msgs/msg/multi_array_layout.h>
#include <std_msgs/msg/string.h>

sensor_msgs__msg__ChannelFloat32 * sensor_msgs__msg__ChannelFloat32__fixture() {
	sensor_msgs__msg__ChannelFloat32 * obj = sensor_msgs__msg__ChannelFloat32__create();
	obj->name.data = "Always outnumbered, never outgunned";
	obj->name.size = 35;
	obj->name.capacity = 36;
	obj->values.data = malloc(sizeof(float)*6);
	obj->values.data[0] = 0.0;
	obj->values.data[1] = 1.1;
	obj->values.data[2] = 2.2;
	obj->values.data[3] = 3.3;
	obj->values.data[4] = 4.4;
	obj->values.data[5] = 5.5;
	obj->values.size = 6;
	obj->values.capacity = 6;
	return obj;
}
sensor_msgs__msg__Illuminance * sensor_msgs__msg__Illuminance__fixture() {
	sensor_msgs__msg__Illuminance * obj = sensor_msgs__msg__Illuminance__create();
	obj->illuminance          = 123456789.987654321;
	obj->variance             = 918273645.546372819;
	obj->header.frame_id.data = "Illuminati is here!\x00";
	obj->header.frame_id.size = 19;
	obj->header.frame_id.capacity = 20;
	obj->header.stamp.sec     = 3600;
	obj->header.stamp.nanosec = 7200;
	return obj;
}
std_msgs__msg__ColorRGBA * std_msgs__msg__ColorRGBA__fixture() {
	std_msgs__msg__ColorRGBA * obj = std_msgs__msg__ColorRGBA__create();
	obj->r = 256.652;
	obj->g = 128.821;
	obj->b = 0.0;
	obj->a = 512.215;
	return obj;
}
std_msgs__msg__Int64MultiArray * std_msgs__msg__Int64MultiArray__fixture() {
	std_msgs__msg__Int64MultiArray * obj = std_msgs__msg__Int64MultiArray__create();
	obj->layout.data_offset = 32;
	std_msgs__msg__MultiArrayDimension__Sequence__init(&obj->layout.dim, 1);
	rosidl_runtime_c__String__init(&obj->layout.dim.data[0].label);
	rosidl_runtime_c__String__assignn(&obj->layout.dim.data[0].label, "MAD_Sequence_0\x00", 14);
	obj->layout.dim.data[0].size = 10;
	obj->layout.dim.data[0].stride = 15;
//	obj->layout.dim.data[0].label.data = "MAD_Sequence_0\x00";
//	obj->layout.dim.data[0].label.size = 14
//	obj->layout.dim.data[0].label.capacity = 15
	rosidl_runtime_c__int64__Sequence__init(&obj->data, 5);
//	obj->data = malloc(sizeof(float)*5);
	obj->data.data[0] = -99999;
	obj->data.data[1] = -8888888888;
	obj->data.data[2] = -777777777777777;
	obj->data.data[3] = -6666666666666666666;
	obj->data.data[4] = -5555555555555555555;
//	obj->data.size = 5;
//	obj->data.capacity = 5;
	return obj;
}
std_msgs__msg__String * std_msgs__msg__String__fixture() {
	std_msgs__msg__String * obj = std_msgs__msg__String__create();
	obj->data.data = "this is a string\x00";
	obj->data.size = 16;
	obj->data.capacity = 17;
	return obj;
}

*/
import "C"
import (
	builtin_interfaces "github.com/tiiuae/rclgo/internal/msgs/builtin_interfaces/msg"
	sensor_msgs "github.com/tiiuae/rclgo/internal/msgs/sensor_msgs/msg"
	std_msgs "github.com/tiiuae/rclgo/internal/msgs/std_msgs/msg"
)

func Fixture_C_sensor_msgs__ChannelFloat32() *C.sensor_msgs__msg__ChannelFloat32 {
	return C.sensor_msgs__msg__ChannelFloat32__fixture()
}
func Fixture_Go_sensor_msgs__ChannelFloat32() *sensor_msgs.ChannelFloat32 {
	return &sensor_msgs.ChannelFloat32{
		Name:   "Always outnumbered, never outgunned",
		Values: []float32{0.0, 1.1, 2.2, 3.3, 4.4, 5.5},
	}
}
func Fixture_C_sensor_msgs__Illuminance() *C.sensor_msgs__msg__Illuminance {
	return C.sensor_msgs__msg__Illuminance__fixture()
}
func Fixture_Go_sensor_msgs__Illuminance() *sensor_msgs.Illuminance {
	return &sensor_msgs.Illuminance{
		Header: std_msgs.Header{
			Stamp: builtin_interfaces.Time{
				Sec:     3600,
				Nanosec: 7200,
			},
			FrameId: "Illuminati is here!",
		},
		Illuminance: 123456789.987654321,
		Variance:    918273645.546372819,
	}
}
func Fixture_C_std_msgs__ColorRGBA() *C.std_msgs__msg__ColorRGBA {
	return C.std_msgs__msg__ColorRGBA__fixture()
}
func Fixture_Go_std_msgs__ColorRGBA() *std_msgs.ColorRGBA {
	return &std_msgs.ColorRGBA{
		R: 256.652,
		G: 128.821,
		B: 0.0,
		A: 512.215,
	}
}
func Fixture_C_std_msgs__String() *C.std_msgs__msg__String {
	return C.std_msgs__msg__String__fixture()
}
func Fixture_Go_std_msgs__String() *std_msgs.String {
	return &std_msgs.String{
		Data: "this is a string",
	}
}
func Fixture_C_std_msgs__Int64MultiArray() *C.std_msgs__msg__Int64MultiArray {
	return C.std_msgs__msg__Int64MultiArray__fixture()
}
func Fixture_Go_std_msgs__Int64MultiArray() *std_msgs.Int64MultiArray {
	return &std_msgs.Int64MultiArray{
		Layout: std_msgs.MultiArrayLayout{
			Dim: []std_msgs.MultiArrayDimension{
				{
					Label:  "MAD_Sequence_0",
					Size:   10,
					Stride: 15,
				},
			},
			DataOffset: 32,
		},
		Data: []int64{-99999, -8888888888, -777777777777777, -6666666666666666666, -5555555555555555555},
	}
}
