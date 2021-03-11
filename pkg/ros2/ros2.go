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
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrmw -lrosidl_runtime_c -lrosidl_typesupport_c -lstd_msgs__rosidl_generator_c -lstd_msgs__rosidl_typesupport_c -lrcutils -lrmw_implementation -lpx4_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_generator_c
#cgo CFLAGS: -I/opt/ros/foxy/include

#include <stdlib.h>
#include <string.h>

#include <rcutils/allocator.h>
#include <rcutils/types/string_array.h>
#include <rcl/graph.h>
#include <rcl/init.h>
#include <rcl/init_options.h>
#include <rcl/subscription.h>
#include <rcl/timer.h>
#include <rcl/time.h>
#include <rcl/wait.h>
#include <rcl/validate_topic_name.h>
#include <rcl/node_options.h>
#include <rcl/node.h>
#include <rmw/get_topic_names_and_types.h>
#include <rmw/names_and_types.h>
#include <rmw/types.h>
#include <rmw/topic_endpoint_info.h>
#include <rmw/topic_endpoint_info_array.h>


///
/// These gowrappers are needed to access C arrays
///
rcl_subscription_t* gowrapper_get_subscription(rcl_subscription_t** subscriptions, ulong i) {
        return subscriptions[i];
}
rcl_timer_t* gowrapper_get_timer(rcl_timer_t** timers, ulong i) {
        return timers[i];
}
rmw_topic_endpoint_info_t* gowrapper_get_rmw_topic_endpoint_info(rmw_topic_endpoint_info_array_t* infos, ulong i) {
	return &(infos->info_array[i]);
}

char* gowrapper_get_rcutils_string_array_index(rcutils_string_array_t* haystack, int i) {
	return *(haystack[i].data);
}

int gowrapper_find_rcutils_string_array_index(rcutils_string_array_t* haystack, char* needle) {
	int i;
	for (i = 0 ; i < haystack->size ; i++) {
		char** data = haystack[i].data;
		if (strcmp(*data, needle) == 0) {
			return i;
		}
	}

	return -1;
}

*/
import "C"
import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/std_msgs"
)

type ROS2Msg interface {
	TypeSupport() unsafe.Pointer //*C.rosidl_message_type_support_t
	PrepareMemory() unsafe.Pointer
	ReleaseMemory(unsafe.Pointer)
}

type RclContext struct {
	rcl_allocator *C.rcutils_allocator_t
	rcl_context   *C.rcl_context_t
	rcl_clock     *C.rcl_clock_t
}

type Subscription struct {
	Topic_name         string
	Ros2MsgType        ROS2Msg
	Rcl_subscription_t *C.rcl_subscription_t
	Callback           func(Subscription, ROS2Msg)
}

type WaitSet struct {
	Timeout        time.Duration
	Subscriptions  []Subscription
	Timers         []C.rcl_timer_t
	Rcl_wait_set_t *C.rcl_wait_set_t
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

	ros2_msg_placeholder := &std_msgs.ColorRGBA{}
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

	subscription, err := SubscriptionCreate(rclContext, rcl_node, "topic_name", &std_msgs.ColorRGBA{}, nil)
	if err != nil {
		fmt.Printf("Error '%+v' SubscriptionCreate.\n", err)
		panic(err)
	}

	subscriptions := []Subscription{subscription}
	timers := []C.rcl_timer_t{*rcl_timer}
	waitSet, err := WaitSetCreate(rclContext, subscriptions, timers, 1000*time.Millisecond)
	if err != nil {
		fmt.Printf("Error '%+v' WaitSetCreate.\n", err)
		panic(err)
	}

	err = WaitSetRun(waitSet)
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

func PublisherCreate(rclContext RclContext, rcl_node *C.rcl_node_t, topic_name string, ros2msg ROS2Msg) (*C.rcl_publisher_t, RCLError) {
	rcl_publisher := (*C.rcl_publisher_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_publisher_t{}))))
	*rcl_publisher = C.rcl_get_zero_initialized_publisher()

	rcl_publisher_options := (*C.rcl_publisher_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_publisher_options_t{}))))
	*rcl_publisher_options = C.rcl_publisher_get_default_options()

	err := ValidateTopicName(topic_name)
	if err != nil {
		return rcl_publisher, err
	}

	var rc C.rcl_ret_t = C.rcl_publisher_init(
		rcl_publisher,
		rcl_node,
		(*C.rosidl_message_type_support_t)(ros2msg.TypeSupport()),
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
		return rcl_timer, ErrorsCastC(rc, fmt.Sprintf("rcl_system_clock_init() failed for timer '%+v'", rcl_timer))
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

func SubscriptionCreate(rclContext RclContext, rcl_node *C.rcl_node_t, topic_name string, ros2msg ROS2Msg, subscriptionCallback func(Subscription, ROS2Msg)) (Subscription, RCLError) {
	var subscription Subscription
	subscription.Rcl_subscription_t = (*C.rcl_subscription_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_subscription_t{}))))
	*subscription.Rcl_subscription_t = C.rcl_get_zero_initialized_subscription()
	subscription.Ros2MsgType = ros2msg
	subscription.Topic_name = topic_name
	subscription.Callback = subscriptionCallback

	err := ValidateTopicName(subscription.Topic_name)
	if err != nil {
		return subscription, err
	}

	rcl_subscription_options_t := (*C.rcl_subscription_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_subscription_options_t{}))))
	*rcl_subscription_options_t = C.rcl_subscription_get_default_options()

	var rc C.rcl_ret_t = C.rcl_subscription_init(
		subscription.Rcl_subscription_t,
		rcl_node,
		(*C.rosidl_message_type_support_t)(ros2msg.TypeSupport()),
		C.CString(topic_name),
		rcl_subscription_options_t)
	if rc != C.RCL_RET_OK {
		return subscription, ErrorsCastC(rc, fmt.Sprintf("Topic name '%s'", topic_name))
	}

	return subscription, nil
}

func PublishersInfoByTopic(rclContext RclContext, rcl_node *C.rcl_node_t, topic_name string) (*C.rmw_topic_endpoint_info_array_t, RCLError) {
	//TODO: This is actually an array of arrays and the memory allocation mechanisms inside ROS2 rcl are more complex! Need to review this on what to do here.
	rmw_topic_endpoint_info_array := (*C.rmw_topic_endpoint_info_array_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rmw_topic_endpoint_info_array_t{}))))
	*rmw_topic_endpoint_info_array = C.rcl_get_zero_initialized_topic_endpoint_info_array()
	var rc C.rcl_ret_t = C.rcl_get_publishers_info_by_topic(rcl_node, rclContext.rcl_allocator, C.CString(topic_name), false, rmw_topic_endpoint_info_array)
	if rc != C.RCL_RET_OK {
		return rmw_topic_endpoint_info_array, ErrorsCast(rc)
	}
	return rmw_topic_endpoint_info_array, nil
}

func TopicGetEndpointInfo(rclContext RclContext, rcl_node *C.rcl_node_t, topic_name string) RCLError {
	//rmw_topic_endpoint_info_array, err := PublishersInfoByTopic(rclContext, rcl_node, topic_name)
	//if err != nil {
	//	return err
	//}

	//var rmw_topic_endpoint_info C.rmw_topic_endpoint_info_t = C.gowrapper_get_rmw_topic_endpoint_info(rmw_topic_endpoint_info_array, 0)
	//rmw_topic_endpoint_info.
	return nil
}

/*func TopicGetTopicTypeSupport(rclContext RclContext, rcl_node *C.rcl_node_t, topic_name string) (C.rosidl_message_type_support_t, RCLError) {
	typeString, err := TopicGetTopicTypeString(rclContext, rcl_node, topic_name)
	if err == nil {
		return nil, err
	}
	parts := strings.Split(typeString, "/")
	if len(parts) == 0 {
		return nil, ErrorsCastC(C.RCL_RET_TOPIC_NAME_INVALID, topic_name)
	}

	//cFuncName := fmt.Sprintf("rosidl_typesupport_c__get_message_type_support_handle__%s__%s__%s", parts[0], parts[1], parts[2])
	return nil, nil
}*/

func TopicGetTopicTypeString(rclContext RclContext, rcl_node *C.rcl_node_t, topic_name string) (string, RCLError) {
	rmw_names_and_types, err := TopicGetTopicNamesAndTypes(rclContext, rcl_node)
	if err != nil {
		return "", err
	}

	var i C.int = C.gowrapper_find_rcutils_string_array_index(&rmw_names_and_types.names, C.CString(topic_name))
	if i == -1 {
		return "", nil
	}
	var data *C.char = C.gowrapper_get_rcutils_string_array_index(rmw_names_and_types.types, i)
	return C.GoString(data), nil
}

func TopicGetTopicNamesAndTypes(rclContext RclContext, rcl_node *C.rcl_node_t) (*C.rmw_names_and_types_t, RCLError) {
	var rmw_node *C.rmw_node_t = C.rcl_node_get_rmw_handle(rcl_node)

	rmw_names_and_types := (*C.rmw_names_and_types_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rmw_names_and_types_t{}))))
	*rmw_names_and_types = C.rmw_get_zero_initialized_names_and_types() // TODO: Array mnemory handling here

	var rc C.rcl_ret_t = (C.rcl_ret_t)(C.rmw_get_topic_names_and_types(rmw_node, rclContext.rcl_allocator, false, rmw_names_and_types)) // rmw_ret_t is aliased to rcl_ret_t
	if rc != 0 {
		return rmw_names_and_types, ErrorsCast(rc)
	}

	return rmw_names_and_types, nil
}

func WaitSetCreate(rclContext RclContext, subscriptions []Subscription, timers []C.rcl_timer_t, timeout time.Duration) (WaitSet, RCLError) {
	waitSet := WaitSet{}
	waitSet.Timeout = timeout
	var number_of_subscriptions C.ulong = 0
	if subscriptions != nil {
		number_of_subscriptions = (C.ulong)(len(subscriptions))
		waitSet.Subscriptions = subscriptions
	}
	var number_of_guard_conditions C.ulong = 0
	var number_of_timers C.ulong = 0
	if timers != nil {
		number_of_timers = (C.ulong)(len(timers))
		waitSet.Timers = timers
	}
	var number_of_clients C.ulong = 0
	var number_of_services C.ulong = 0
	var number_of_events C.ulong = 0

	var rcl_wait_set C.rcl_wait_set_t = C.rcl_get_zero_initialized_wait_set()
	waitSet.Rcl_wait_set_t = &rcl_wait_set
	var rc C.rcl_ret_t = C.rcl_wait_set_init(waitSet.Rcl_wait_set_t, number_of_subscriptions, number_of_guard_conditions, number_of_timers, number_of_clients, number_of_services, number_of_events, rclContext.rcl_context, *rclContext.rcl_allocator)
	if rc != C.RCL_RET_OK {
		return waitSet, ErrorsCast(rc)
	}

	return waitSet, nil
}

/**
Using wait set manually to be able to better control the parameters to callback handlers.
rclc subscriptions do not pass the rcl_subscription_t to the callback,
making it impossible to dynamically dispatch messages to the corresponding callback handlers
*/
func WaitSetRun(waitSet WaitSet) RCLError {
	for {
		if !C.rcl_wait_set_is_valid(waitSet.Rcl_wait_set_t) {
			//#define RCL_RET_WAIT_SET_INVALID 900
			return ErrorsCastC(900, fmt.Sprintf("rcl_wait_set_is_valid() failed for wait_set='%v'", waitSet))
		}
		var rc C.rcl_ret_t = C.rcl_wait_set_clear(waitSet.Rcl_wait_set_t)
		if rc != C.RCL_RET_OK {
			return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_clear() failed for wait_set='%v'", waitSet))
		}
		for i := 0; i < len(waitSet.Subscriptions); i++ {
			rc = C.rcl_wait_set_add_subscription(waitSet.Rcl_wait_set_t, waitSet.Subscriptions[i].Rcl_subscription_t, nil)
			if rc != C.RCL_RET_OK {
				return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_subscription() failed for wait_set='%v'", waitSet))
			}
		}

		rc = C.rcl_wait(waitSet.Rcl_wait_set_t, (C.long)(waitSet.Timeout))
		if rc == C.RCL_RET_TIMEOUT {
			continue
		}
		var i C.ulong
		// Check timers. Guard against internal state representation mismatch. Due to some software bug the lists of timers could easily get out of sync. AND lead to very very difficult to detect bugs.
		if (int)(waitSet.Rcl_wait_set_t.size_of_timers) != len(waitSet.Timers) {
			panic(fmt.Sprintf(
				"Wait set timers count mismatch! rcl_wait_set.size_of_timers='%d' != len(waitSet.Timers)='%d'",
				(int)(waitSet.Rcl_wait_set_t.size_of_subscriptions),
				len(waitSet.Subscriptions)))
		}
		for i = 0; i < waitSet.Rcl_wait_set_t.size_of_timers; i++ {
			var is_timer_ready_to_call C.bool = false
			if rcl_timer := C.gowrapper_get_timer(waitSet.Rcl_wait_set_t.timers, i); rcl_timer != nil {
				rc = C.rcl_timer_is_ready(rcl_timer, &is_timer_ready_to_call)
				if rc != C.RCL_RET_OK {
					return ErrorsCastC(rc, fmt.Sprintf("rcl_timer_is_ready() failed for wait_set='%v'", waitSet))
				}
				if is_timer_ready_to_call {
					timer_callback(*rcl_timer)
				}
			}
		}
		// Check subscriptions. Guard against internal state representation mismatch. Due to some software bug the lists of subscriptions could easily get out of sync. AND lead to very very difficult to detect bugs.
		if (int)(waitSet.Rcl_wait_set_t.size_of_subscriptions) != len(waitSet.Subscriptions) {
			panic(fmt.Sprintf(
				"Wait set subscriptions count mismatch! rcl_wait_set.size_of_subscriptions='%d' != len(waitSet.Subscriptions)='%d'",
				(int)(waitSet.Rcl_wait_set_t.size_of_subscriptions),
				len(waitSet.Subscriptions)))
		}
		for i = 0; i < waitSet.Rcl_wait_set_t.size_of_subscriptions; i++ {
			subscription := waitSet.Subscriptions[i]

			rmw_message_info := (*C.rmw_message_info_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rmw_message_info_t{}))))
			*rmw_message_info = C.rmw_get_zero_initialized_message_info()
			defer C.free(unsafe.Pointer(rmw_message_info))

			ros2_msg_receive_buffer := subscription.Ros2MsgType.PrepareMemory()
			defer subscription.Ros2MsgType.ReleaseMemory(ros2_msg_receive_buffer)

			rc = C.rcl_take(subscription.Rcl_subscription_t, ros2_msg_receive_buffer, rmw_message_info, nil)
			if rc != C.RCL_RET_OK {
				return ErrorsCastC(rc, fmt.Sprintf("rcl_take() failed for wait_set='%v'", waitSet))
			} else {
				subscription_callback(subscription, ros2_msg_receive_buffer, rmw_message_info)
			}

		}
	}
}

func subscription_callback(subscription Subscription, ros2_msg_receive_buffer unsafe.Pointer, rmw_message_info *C.rmw_message_info_t) {
	eee := (*std_msgs.ColorRGBA)(ros2_msg_receive_buffer)
	fmt.Printf("ROS2 Message receive buffer: '%+v'", ros2_msg_receive_buffer)
	fmt.Printf("ROS2 Message receive buffer: '%+v'", eee)
}

func timer_callback(rcl_timer C.rcl_timer_t) {
	fmt.Printf("ROS2 timer callback for timer '%+v'", rcl_timer)
}
