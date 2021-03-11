package px4_msgs

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_cpp -lrosidl_typesupport_c -lstd_msgs__rosidl_generator_c -lpx4_msgs__rosidl_typesupport_c
#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <px4_msgs/msg/sensor_combined.h>
*/
import "C"
import "unsafe"

const (
	sensorCombinedRELATIVETIMESTAMPIINVALID int32 = 2147483647
	sensorCombinedCLIPPINGX                 uint8 = 1
	sensorCombinedCLIPPINGY                 uint8 = 2
	sensorCombinedCLIPPINGZ                 uint8 = 4
)

//SensorCombined ROS px4 message struct
type SensorCombined struct {
	Timestamp                      uint64
	GyroRad                        [3]float32
	GyroIntegralDt                 uint32
	AccelerometerTimestampRelative int32
	AccelerometerMS2               [3]float32
	AccelerometerIntegralDt        uint32
	AccelerometerClipping          uint8
}

func (t *SensorCombined) TypeSupport() *C.rosidl_message_type_support_t {
	return C.rosidl_typesupport_c__get_message_type_support_handle__px4_msgs__msg__SensorCombined()
}
func (t *SensorCombined) PrepareMemory() unsafe.Pointer { //returns *C.px4_msgs__msg__SensorCombined
	return (unsafe.Pointer)(C.px4_msgs__msg__SensorCombined__create())
}
func (t *SensorCombined) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.px4_msgs__msg__SensorCombined__destroy((*C.px4_msgs__msg__SensorCombined)(pointer_to_free))
}
