/*
  Deliberate trial and error have been conducted in finding the best way of interfacing with rcl or rclc.

  rclc was initialyl considered, but:
  Executor subscription callback doesn't include the subscription, only the ros2 message.
  Thus we cannot intelligently and dynamically dispatch the ros2 message to the correct
  subscription callback on the golang layer.
  rcl wait_set has much more granular way of defining how the received messages are handled and
  allows for a more Golang-way of handling dynamic callbacks
*/
package ros2

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lstd_msgs__rosidl_generator_c -lstd_msgs__rosidl_typesupport_c -lrcutils -lrmw_implementation -lpx4_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_generator_c
#cgo CFLAGS: -I/opt/ros/foxy/include

#include <stdlib.h>

#include <rcutils/allocator.h>
#include <rcl/init.h>
#include <rcl/init_options.h>
#include <rcl/subscription.h>
#include <rcl/timer.h>
#include <rcl/time.h>
#include <rcl/wait.h>
#include <rcl/node_options.h>
#include <rcl/node.h>

///
/// These gowrappers are needed to access C arrays
///
rcl_subscription_t* gowrapper_get_subscription(rcl_subscription_t** subscriptions, ulong i) {
        return subscriptions[i];
}
rcl_timer_t* gowrapper_get_timer(rcl_timer_t** timers, ulong i) {
        return timers[i];
}

*/
import "C"
import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
)

type RclContext struct {
	rcl_allocator *C.rcutils_allocator_t
	rcl_context   *C.rcl_context_t
	rcl_clock     *C.rcl_clock_t
}

type Subscription struct {
	topic              string
	ros2MsgType        ros2types.ROS2Msg
	rcl_subscription_t C.rcl_subscription_t
}

func main() {
	rclContext, err := RclInit()
	if err != nil {
		fmt.Printf("Error '%+v' RclInit.\n", err)
		panic(err)
	}

	rcl_node, err := NodeCreate(rclContext, "node_name", "")
	if err != nil {
		fmt.Printf("Error '%+v' node_create.\n", err)
		panic(err)
	}

	ros2_msg_placeholder := &ros2types.StdMsgs_ColorRGBA{}
	rcl_publisher, err := PublisherCreate(rclContext, rcl_node, "topic_name", ros2_msg_placeholder)
	if err != nil {
		fmt.Printf("Error '%+v' publisher_create.\n", err)
		panic(err)
	}

	rcl_timer, err := TimerCreate(rclContext, 0)
	if err != nil {
		fmt.Printf("Error '%+v' TimerCreate.\n", err)
		panic(err)
	}

	subscription, err := SubscriptionCreate(rclContext, rcl_node, "topic_name", &ros2types.StdMsgs_ColorRGBA{}, nil)

	subscriptions := []Subscription{subscription}
	timers := []C.rcl_timer_t{*rcl_timer}
	rcl_wait_set, err := WaitSetCreate(rclContext, subscriptions, timers)
	if err != nil {
		fmt.Printf("Error '%+v' WaitSetCreate.\n", err)
		panic(err)
	}

	err = WaitSetRun(rcl_wait_set, subscription.rcl_subscription_t, ros2_msg_placeholder)
	if err != nil {
		fmt.Printf("Error '%+v' WaitSetRun.\n", err)
		panic(err)
	}

	fmt.Printf("%v", rcl_publisher)
}

func RclInit() (RclContext, RCLError) {
	var rc C.rcl_ret_t

	var argc C.int = 0
	var argv **C.char

	rclContext := RclContext{}
	rclContext.rcl_context = (*C.rcl_context_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_context_t{}))))
	*rclContext.rcl_context = C.rcl_get_zero_initialized_context()

	/* Instead of receiving the rcl_allocator_t as a golang struct,
	   prepare C memory from heap to receive a copy of the rcl allocator.
	   This way Golang wont mess with the rcl_allocator_t memory location
	   and complaing about nested Golang pointer passed over cgo */
	rclContext.rcl_allocator = (*C.rcl_allocator_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_allocator_t{}))))
	*rclContext.rcl_allocator = C.rcl_get_default_allocator()
	// TODO: Free C.free(rclContext.rcl_allocator)

	rcl_init_options := (*C.rcl_init_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_init_options_t{}))))
	*rcl_init_options = C.rcl_get_zero_initialized_init_options()
	rc = C.rcl_init_options_init(rcl_init_options, *rclContext.rcl_allocator)
	if rc != C.RCL_RET_OK {
		return rclContext, ErrorsCast(rc)
	}

	rc = C.rcl_init(argc, argv, rcl_init_options, rclContext.rcl_context)
	if rc != C.RCL_RET_OK {
		return rclContext, ErrorsCast(rc)
	}

	return rclContext, nil
}

func NodeCreate(rclContext RclContext, node_name string, namespace string) (*C.rcl_node_t, RCLError) {
	ns := strings.ReplaceAll(namespace, "/", "")
	ns = strings.ReplaceAll(ns, "-", "")

	rcl_node_options := (*C.rcl_node_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_node_options_t{}))))
	*rcl_node_options = C.rcl_node_get_default_options()

	rcl_node := (*C.rcl_node_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_node_t{}))))
	*rcl_node = C.rcl_get_zero_initialized_node()

	var rc C.rcl_ret_t = C.rcl_node_init(rcl_node, C.CString(node_name), C.CString(ns), rclContext.rcl_context, rcl_node_options)
	if rc != C.RCL_RET_OK {
		fmt.Printf("Error '%d' in rcl_node_init\n", (int)(rc))
		return rcl_node, ErrorsCast(rc)
	}

	return rcl_node, nil
}

func PublisherCreate(rclContext RclContext, rcl_node *C.rcl_node_t, topic_name string, ros2msg ros2types.ROS2Msg) (*C.rcl_publisher_t, RCLError) {
	rcl_publisher := (*C.rcl_publisher_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_publisher_t{}))))
	*rcl_publisher = C.rcl_get_zero_initialized_publisher()

	rcl_publisher_options := (*C.rcl_publisher_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_publisher_options_t{}))))
	*rcl_publisher_options = C.rcl_publisher_get_default_options()

	var rc C.rcl_ret_t = C.rcl_publisher_init(
		rcl_publisher,
		rcl_node,
		ros2msg.TypeSupport(),
		C.CString(topic_name),
		rcl_publisher_options)
	if rc != C.RCL_RET_OK {
		return rcl_publisher, ErrorsCast(rc)
	}

	return rcl_publisher, nil
}

func TimerCreate(rclContext RclContext, timeout time.Duration) (*C.rcl_timer_t, RCLError) {
	if timeout == 0 {
		timeout = 1000 * time.Millisecond
	}

	rcl_timer := (*C.rcl_timer_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_timer_t{}))))
	*rcl_timer = C.rcl_get_zero_initialized_timer()

	var rc C.rcl_ret_t = C.rcl_system_clock_init(rclContext.rcl_clock, rclContext.rcl_allocator)
	if rc != C.RCL_RET_OK {
		return rcl_timer, ErrorsCastC(rc, fmt.Sprint("rcl_system_clock_init() failed for timer '%+v'", rcl_timer))
	}

	rc = C.rcl_timer_init(
		rcl_timer,
		rclContext.rcl_clock,
		rclContext.rcl_context,
		(C.long)(timeout),
		nil,
		*rclContext.rcl_allocator)
	if rc != C.RCL_RET_OK {
		return rcl_timer, ErrorsCast(rc)
	}

	return rcl_timer, nil
}

func SubscriptionCreate(rclContext RclContext, rcl_node *C.rcl_node_t, topic_name string, ros2msg ros2types.ROS2Msg, rcl_subscription_callback interface{}) (Subscription, RCLError) {
	var subscription Subscription
	subscription.rcl_subscription_t = C.rcl_get_zero_initialized_subscription()
	subscription.ros2MsgType = ros2msg
	subscription.topic = topic_name

	rcl_subscription_options_t := (*C.rcl_subscription_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_subscription_options_t{}))))
	*rcl_subscription_options_t = C.rcl_subscription_get_default_options()

	var rc C.rcl_ret_t = C.rcl_subscription_init(
		&subscription.rcl_subscription_t,
		rcl_node,
		ros2msg.TypeSupport(),
		C.CString(topic_name),
		rcl_subscription_options_t)
	if rc != C.RCL_RET_OK {
		return subscription, ErrorsCastC(rc, fmt.Sprintf("Topic name '%s'", topic_name))
	}

	if rcl_subscription_callback != nil {

	}

	return subscription, nil
}

func WaitSetCreate(rclContext RclContext, subscriptions []Subscription, timers []C.rcl_timer_t) (C.rcl_wait_set_t, RCLError) {
	//var number_of_subscriptions C.ulong = 0
	var number_of_guard_conditions C.ulong = 0
	//var number_of_timers C.ulong = 0
	var number_of_clients C.ulong = 0
	var number_of_services C.ulong = 0
	var number_of_events C.ulong = 0

	var rcl_wait_set C.rcl_wait_set_t = C.rcl_get_zero_initialized_wait_set()
	var rc C.rcl_ret_t = C.rcl_wait_set_init(&rcl_wait_set, (C.ulong)(len(subscriptions)), (C.ulong)(number_of_guard_conditions), (C.ulong)(len(timers)), (C.ulong)(number_of_clients), (C.ulong)(number_of_services), (C.ulong)(number_of_events), rclContext.rcl_context, *rclContext.rcl_allocator)
	if rc != C.RCL_RET_OK {
		return rcl_wait_set, ErrorsCast(rc)
	}

	return rcl_wait_set, nil
}

/**
Using wait set manually to be able to better control the parameters to callback handlers.
rclc subscriptions do not pass the rcl_subscription_t to the callback,
making it impossible to dynamically dispatch messages to the corresponding callback handlers
*/
func WaitSetRun(rcl_wait_set C.rcl_wait_set_t, rcl_subscription C.rcl_subscription_t, ros2_msg_placeholder ros2types.ROS2Msg) RCLError {
	for {
		if !C.rcl_wait_set_is_valid(&rcl_wait_set) {
			//#define RCL_RET_WAIT_SET_INVALID 900
			return ErrorsCastC(900, fmt.Sprintf("rcl_wait_set_is_valid() failed for wait_set='%v'", rcl_wait_set))
		}
		var rc C.rcl_ret_t = C.rcl_wait_set_clear(&rcl_wait_set)
		if rc != C.RCL_RET_OK {
			return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_clear() failed for wait_set='%v'", rcl_wait_set))
		}
		rc = C.rcl_wait_set_add_subscription(&rcl_wait_set, &rcl_subscription, nil)
		if rc != C.RCL_RET_OK {
			return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_subscription() failed for wait_set='%v'", rcl_wait_set))
		}
		var wait_duration time.Duration = 1000 * time.Millisecond
		rc = C.rcl_wait(&rcl_wait_set, (C.long)(wait_duration))
		if rc == C.RCL_RET_TIMEOUT {
			continue
		}
		var i C.ulong
		for i = 0; i < rcl_wait_set.size_of_timers; i++ {
			var is_timer_ready_to_call C.bool = false
			if rcl_timer := C.gowrapper_get_timer(rcl_wait_set.timers, i); rcl_timer != nil {
				rc = C.rcl_timer_is_ready(rcl_timer, &is_timer_ready_to_call)
				if rc != C.RCL_RET_OK {
					return ErrorsCastC(rc, fmt.Sprintf("rcl_timer_is_ready() failed for wait_set='%v'", rcl_wait_set))
				}
				if is_timer_ready_to_call {
					timer_callback(*rcl_timer)
				}
			}
		}
		for i = 0; i < rcl_wait_set.size_of_subscriptions; i++ {
			if rcl_subscription := C.gowrapper_get_subscription(rcl_wait_set.subscriptions, i); rcl_subscription != nil {
				ros2_msg_receive_buffer := ros2_msg_placeholder.PrepareMemory()
				rc = C.rcl_take(rcl_subscription, ros2_msg_receive_buffer, nil, nil)
				if rc != C.RCL_RET_OK {
					return ErrorsCastC(rc, fmt.Sprintf("rcl_take() failed for wait_set='%v'", rcl_wait_set))
				} else {
					subscription_callback(*rcl_subscription, ros2_msg_receive_buffer)
				}
				ros2_msg_placeholder.ReleaseMemory(ros2_msg_receive_buffer)
			}
		}
	}
}

func subscription_callback(rcl_subscription C.rcl_subscription_t, ros2_msg_receive_buffer unsafe.Pointer) {
	eee := (*ros2types.StdMsgs_ColorRGBA)(ros2_msg_receive_buffer)
	fmt.Printf("ROS2 Message receive buffer: '%+v'", ros2_msg_receive_buffer)
	fmt.Printf("ROS2 Message receive buffer: '%+v'", eee)
}

func timer_callback(rcl_timer C.rcl_timer_t) {
	fmt.Printf("ROS2 timer callback for timer '%+v'", rcl_timer)
}
