package ros2

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lstd_msgs__rosidl_generator_c -lstd_msgs__rosidl_typesupport_c -lrcutils -lrmw_implementation -lpx4_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_generator_c
#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>

#include <std_msgs/msg/string.h>
#include <std_msgs/msg/color_rgba.h>

const rosidl_message_type_support_t * rosidl_typesupport_handle_std_msgs__msg__String_gowrapper() {
        return rosidl_typesupport_c__get_message_type_support_handle__std_msgs__msg__String();
}

const rosidl_message_type_support_t * rosidl_typesupport_handle_std_msgs__msg__ColorRGBA_gowrapper() {
        return rosidl_typesupport_c__get_message_type_support_handle__std_msgs__msg__ColorRGBA();
}


*/
import "C"
import "unsafe"

type ROS2Type uintptr

type ROS2Msg interface {
	TypeSupport() *C.rosidl_message_type_support_t
	PrepareMemory() unsafe.Pointer
	ReleaseMemory(unsafe.Pointer)
}

type StdMsgs_ColorRGBA struct {
	R float32
	G float32
	B float32
	A float32
}

func (t *StdMsgs_ColorRGBA) TypeSupport() *C.rosidl_message_type_support_t {
	return C.rosidl_typesupport_c__get_message_type_support_handle__std_msgs__msg__ColorRGBA()
}
func (t *StdMsgs_ColorRGBA) PrepareMemory() unsafe.Pointer { //returns *C.std_msgs__msg__ColorRGBA
	return (unsafe.Pointer)(C.std_msgs__msg__ColorRGBA__create())
}
func (t *StdMsgs_ColorRGBA) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.std_msgs__msg__ColorRGBA__destroy((*C.std_msgs__msg__ColorRGBA)(pointer_to_free))
}

type StdMsgs_String struct {
	Data     *C.char
	Size     int
	Capacity int
}

func (t *StdMsgs_String) TypeSupport() *C.rosidl_message_type_support_t {
	return C.rosidl_typesupport_c__get_message_type_support_handle__std_msgs__msg__String()
}
func (t *StdMsgs_String) PrepareMemory() unsafe.Pointer { //returns *C.std_msgs__msg__String
	return (unsafe.Pointer)(C.std_msgs__msg__String__create())
}
func (t *StdMsgs_String) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.std_msgs__msg__String__destroy((*C.std_msgs__msg__String)(pointer_to_free))
}
