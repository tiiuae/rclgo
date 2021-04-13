/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
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
static void setString(const char*argv[], int i, const char *str) { // TODO replace with Go pointer arithmetics
	argv[i] = str;
}

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

int gowrapper_find_rcutils_string_array_index(rcutils_string_array_t* haystack, char* needle, int needle_size) {
	int i;
	for (i = 0 ; i < haystack->size ; i++) {
		char** data = haystack[i].data;
		if (strncmp(*data, needle, needle_size) == 0) {
			return i;
		}
	}

	return -1;
}


void print_gid(rmw_gid_t gid) {
	printf("gid:\n'%s'\n", gid.data); // gid.data looks like gibberish
}

*/
import "C"
import (
	"container/list"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
)

/*
Keeps track of all the C entities initialized, so we can later free them
*/
type rclEntityWrapper struct {
	rcl_allocator_t    *C.rcutils_allocator_t
	rcl_context_t      *C.rcl_context_t
	clock              *Clock
	rcl_init_options_t *C.rcl_init_options_t
	publishers         list.List // []*Publisher
	subscriptions      list.List // []*Subscription
	timers             list.List // []*Timer
}

func (self *rclEntityWrapper) Fini() *RCLErrors {
	var rclErrors *RCLErrors
	var rc C.rcl_ret_t
	rc = C.rcl_init_options_fini(self.rcl_init_options_t)
	if rc != C.RCL_RET_OK {
		rclErrors = RCLErrorsPut(rclErrors, ErrorsCastC(rc, fmt.Sprintf("C.rcl_init_options_fini(%+v)", self.rcl_init_options_t)))
	} else {
		self.rcl_init_options_t = nil
	}
	rc = C.rcl_shutdown(self.rcl_context_t)
	if rc != C.RCL_RET_OK {
		rclErrors = RCLErrorsPut(rclErrors, ErrorsCastC(rc, fmt.Sprintf("C.rcl_shutdown(%+v)", self.rcl_context_t)))
	} else {
		C.free(unsafe.Pointer(self.rcl_context_t))
		self.rcl_context_t = nil
	}
	C.free(unsafe.Pointer(self.rcl_allocator_t))
	self.rcl_allocator_t = nil

	rc = C.rcl_clock_fini(self.clock.rcl_clock_t)
	if rc != C.RCL_RET_OK {
		rclErrors = RCLErrorsPut(rclErrors, ErrorsCastC(rc, fmt.Sprintf("C.rcl_clock_fini(%+v)", self.clock.rcl_clock_t)))
	} else {
		self.clock = nil
	}
	return rclErrors
}

/*
RCLContext has a key rclEntities which points to the rclEntityWrapper
*/
type RCLContext context.Context

func getRCLEntities(ctx RCLContext) *rclEntityWrapper {
	return ctx.Value("rclEntities").(*rclEntityWrapper)
}

type Clock struct {
	rcl_clock_t *C.rcl_clock_t
}

type Node struct {
	rcl_node_t *C.rcl_node_t
	rclContext RCLContext
}

type Publisher struct {
	TopicName               string
	rcl_publisher_options_t *C.rcl_publisher_options_t
	rcl_publisher_t         *C.rcl_publisher_t
}

type Subscription struct {
	TopicName                  string
	Ros2MsgType                ros2types.ROS2Msg
	rcl_subscription_t         *C.rcl_subscription_t
	rcl_subscription_options_t *C.rcl_subscription_options_t
	Callback                   func(*Subscription, unsafe.Pointer, *RmwMessageInfo)
}

type Timer struct {
	rcl_timer_t *C.rcl_timer_t
	Callback    func(*Timer)
}

type WaitSet struct {
	Timeout        time.Duration
	Subscriptions  []*Subscription
	Timers         []*Timer
	rcl_wait_set_t *C.rcl_wait_set_t
}

type RmwMessageInfo struct {
	SourceTimestamp   time.Time
	ReceivedTimestamp time.Time
	FromIntraProcess  bool
}

type Rcl_clock_type_t uint32

var RCL_CLOCK_UNINITIALIZED Rcl_clock_type_t = 0
var RCL_ROS_TIME Rcl_clock_type_t = 1
var RCL_SYSTEM_TIME Rcl_clock_type_t = 2
var RCL_STEADY_TIME Rcl_clock_type_t = 3

/*
NewRCLContext initializes a new RCL context.

parent can be nil, a new context.Background is created
clockType can be nil, then no clock is initialized, you can later initialize it with NewClock()
osArgs can be nil
*/
func NewRCLContext(parent context.Context, clockType Rcl_clock_type_t, osArgs []string) (RCLContext, RCLError) {
	rclEntities, rclError := rclInit(osArgs)
	if rclError != nil {
		return nil, rclError
	}

	if parent == nil {
		parent = context.Background()
	}
	newCtx := (RCLContext)(context.WithValue(parent, "rclEntities", rclEntities))

	if clockType != 0 {
		_, err := NewClock(newCtx, clockType)
		if err != nil {
			return newCtx, err
		}
	}

	return newCtx, nil

	/*	go func() {
		<-newCtx.Done()
		rclContext.Fini()
	}()*/
}

/*
NewRCLContextChild TODO:
- Example usage of nested contexts to init ROS2 and then create nodes etc for a nested context.
- Cleanup partially one context at a time.
*/
func NewRCLContextChild(parent context.Context) (*RCLContext, RCLError) {
	return nil, nil
}

func rclInit(osArgs []string) (*rclEntityWrapper, RCLError) {
	var rc C.rcl_ret_t

	rclEntityWrapper := rclEntityWrapper{}
	rclEntityWrapper.rcl_context_t = (*C.rcl_context_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_context_t{}))))
	*rclEntityWrapper.rcl_context_t = C.rcl_get_zero_initialized_context()

	/* Instead of receiving the rcl_allocator_t as a golang struct,
	   prepare C memory from heap to receive a copy of the rcl allocator.
	   This way Golang wont mess with the rcl_allocator_t memory location
	   and complaing about nested Golang pointer passed over cgo */
	rclEntityWrapper.rcl_allocator_t = (*C.rcl_allocator_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_allocator_t{}))))
	*rclEntityWrapper.rcl_allocator_t = C.rcl_get_default_allocator()
	// TODO: Free C.free(rclEntityWrapper.rcl_allocator)

	rclEntityWrapper.rcl_init_options_t = (*C.rcl_init_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_init_options_t{}))))
	*rclEntityWrapper.rcl_init_options_t = C.rcl_get_zero_initialized_init_options()
	rc = C.rcl_init_options_init(rclEntityWrapper.rcl_init_options_t, *rclEntityWrapper.rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		rclEntityWrapper.Fini()
		return nil, ErrorsCast(rc)
	}
	rc = rclInitWithGoARGV(osArgs, rclEntityWrapper)
	if rc != C.RCL_RET_OK {
		rclEntityWrapper.Fini()
		return nil, ErrorsCast(rc)
	}

	return &rclEntityWrapper, nil
}

func rclInitWithGoARGV(osArgs []string, rclEntityWrapper rclEntityWrapper) C.int {
	if osArgs == nil {
		osArgs = os.Args
	}
	argc := C.int(len(osArgs))
	argv := (**C.char)(C.malloc((C.size_t)((C.int)(unsafe.Sizeof(uintptr(1))) * argc)))
	for i, arg := range osArgs {
		str := C.CString(arg)
		C.setString(argv, C.int(i), str)
		defer C.free(unsafe.Pointer(str))
	}

	defer C.free(unsafe.Pointer(argv))

	return C.rcl_init(argc, argv, rclEntityWrapper.rcl_init_options_t, rclEntityWrapper.rcl_context_t)
}

func NewNode(rclContext RCLContext, node_name string, namespace string) (*Node, RCLError) {
	ns := strings.ReplaceAll(namespace, "/", "")
	ns = strings.ReplaceAll(ns, "-", "")

	rcl_node_options := (*C.rcl_node_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_node_options_t{}))))
	*rcl_node_options = C.rcl_node_get_default_options()

	rcl_node := (*C.rcl_node_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_node_t{}))))
	*rcl_node = C.rcl_get_zero_initialized_node()

	var rc C.rcl_ret_t = C.rcl_node_init(rcl_node, C.CString(node_name), C.CString(ns), getRCLEntities(rclContext).rcl_context_t, rcl_node_options)
	if rc != C.RCL_RET_OK {
		fmt.Printf("Error '%d' in rcl_node_init\n", (int)(rc))
		return nil, ErrorsCast(rc)
	}

	return &Node{rcl_node_t: rcl_node, rclContext: rclContext}, nil
}

func (self *Node) NewPublisher(topic_name string, ros2msg ros2types.ROS2Msg) (*Publisher, RCLError) {
	rcl_publisher := (*C.rcl_publisher_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_publisher_t{}))))
	*rcl_publisher = C.rcl_get_zero_initialized_publisher()

	rcl_publisher_options := (*C.rcl_publisher_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_publisher_options_t{}))))
	*rcl_publisher_options = C.rcl_publisher_get_default_options()

	err := ValidateTopicName(topic_name)
	if err != nil {
		return nil, err
	}

	var rc C.rcl_ret_t = C.rcl_publisher_init(
		rcl_publisher,
		self.rcl_node_t,
		(*C.rosidl_message_type_support_t)(ros2msg.TypeSupport()),
		C.CString(topic_name),
		rcl_publisher_options)
	if rc != C.RCL_RET_OK {
		return nil, ErrorsCast(rc)
	}

	publisher := &Publisher{
		TopicName:               topic_name,
		rcl_publisher_options_t: rcl_publisher_options,
		rcl_publisher_t:         rcl_publisher,
	}

	ctx := getRCLEntities(self.rclContext).publishers
	ctx.PushFront(publisher)
	return publisher, nil
}

func (self *Publisher) Publish(ros2msg ros2types.ROS2Msg) RCLError {
	var rc C.rcl_ret_t

	ptr := ros2msg.AsCStruct()
	defer ros2msg.ReleaseMemory(unsafe.Pointer(ptr))

	rc = C.rcl_publish(self.rcl_publisher_t, ptr, nil)
	if rc != C.RCL_RET_OK {
		return ErrorsCastC(rc, fmt.Sprintf("rcl_publish() failed for publisher '%+v'", self))
	}
	return nil
}

func NewClock(rclContext RCLContext, clockType Rcl_clock_type_t) (*Clock, RCLError) {
	if clockType == 0 {
		clockType = RCL_ROS_TIME
	}
	rcl_clock := (*C.rcl_clock_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_clock_t{})))) //rcl_clock_init() doc says "This will allocate all necessary internal structures, and initialize variables.". The parameter is invalid if no memory allocated beforehand.
	var rc C.rcl_ret_t = C.rcl_clock_init(uint32(clockType), rcl_clock, getRCLEntities(rclContext).rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		return nil, ErrorsCast(rc)
	}
	c := &Clock{rcl_clock_t: rcl_clock}
	re := getRCLEntities(rclContext)
	re.clock = c
	return c, nil
}

func NewTimer(rclContext RCLContext, timeout time.Duration, timer_callback func(*Timer)) (*Timer, RCLError) {
	var rc C.rcl_ret_t
	re := getRCLEntities(rclContext)

	if timeout == 0 {
		timeout = 1000 * time.Millisecond
	}
	timer := &Timer{}
	timer.Callback = timer_callback

	timer.rcl_timer_t = (*C.rcl_timer_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_timer_t{}))))
	*timer.rcl_timer_t = C.rcl_get_zero_initialized_timer()

	if re.clock.rcl_clock_t == nil {
		var err RCLError
		_, err = NewClock(rclContext, C.RCL_SYSTEM_TIME)
		if err != nil {
			return timer, ErrorsCastC(C.int(err.rcl_ret()), fmt.Sprintf("Forwarding error from '%s'", err.Error()))
		}
	}

	rc = C.rcl_timer_init(
		timer.rcl_timer_t,
		re.clock.rcl_clock_t,
		re.rcl_context_t,
		(C.long)(timeout),
		nil,
		*re.rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		return timer, ErrorsCast(rc)
	}

	return timer, nil
}
func (self *Timer) GetTimeUntilNextCall() (int64, RCLError) {
	var rc C.rcl_ret_t
	time_until_next_call := (*C.int64_t)(C.malloc((C.size_t)(8)))
	defer C.free(unsafe.Pointer(time_until_next_call))

	rc = C.rcl_timer_get_time_until_next_call(self.rcl_timer_t, time_until_next_call)
	if rc != C.RCL_RET_OK {
		return 0, ErrorsCast(rc)
	}
	return int64(*time_until_next_call), nil
}

func (self *Timer) Reset() RCLError {
	var rc C.rcl_ret_t
	rc = C.rcl_timer_reset(self.rcl_timer_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}
	return nil
}

func (self *Node) NewSubscription(topic_name string, ros2msg ros2types.ROS2Msg, subscriptionCallback func(*Subscription, unsafe.Pointer, *RmwMessageInfo)) (*Subscription, RCLError) {
	var subscription Subscription
	subscription.rcl_subscription_t = (*C.rcl_subscription_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_subscription_t{}))))
	*subscription.rcl_subscription_t = C.rcl_get_zero_initialized_subscription()
	subscription.Ros2MsgType = ros2msg
	subscription.TopicName = topic_name
	subscription.Callback = subscriptionCallback

	err := ValidateTopicName(subscription.TopicName)
	if err != nil {
		return &subscription, err
	}

	sops := C.rcl_subscription_get_default_options()
	subscription.rcl_subscription_options_t = &sops

	var rc C.rcl_ret_t = C.rcl_subscription_init(
		subscription.rcl_subscription_t,
		self.rcl_node_t,
		(*C.rosidl_message_type_support_t)(ros2msg.TypeSupport()),
		C.CString(topic_name),
		subscription.rcl_subscription_options_t)
	if rc != C.RCL_RET_OK {
		return &subscription, ErrorsCastC(rc, fmt.Sprintf("Topic name '%s'", topic_name))
	}

	return &subscription, nil
}

func PublishersInfoByTopic(rclContext RCLContext, rcl_node *C.rcl_node_t, topic_name string) (*C.rmw_topic_endpoint_info_array_t, RCLError) {
	re := getRCLEntities(rclContext)
	//TODO: This is actually an array of arrays and the memory allocation mechanisms inside ROS2 rcl are more complex! Need to review this on what to do here.
	rmw_topic_endpoint_info_array := (*C.rmw_topic_endpoint_info_array_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rmw_topic_endpoint_info_array_t{}))))
	*rmw_topic_endpoint_info_array = C.rcl_get_zero_initialized_topic_endpoint_info_array()
	var rc C.rcl_ret_t = C.rcl_get_publishers_info_by_topic(rcl_node, re.rcl_allocator_t, C.CString(topic_name), false, rmw_topic_endpoint_info_array)
	if rc != C.RCL_RET_OK {
		return rmw_topic_endpoint_info_array, ErrorsCast(rc)
	}
	return rmw_topic_endpoint_info_array, nil
}

func TopicGetEndpointInfo(rclContext RCLContext, rcl_node *C.rcl_node_t, topic_name string) RCLError {
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

func TopicGetTopicTypeString(rclContext RCLContext, rcl_node *C.rcl_node_t, topic_name string) (string, RCLError) {
	rmw_names_and_types, err := TopicGetTopicNamesAndTypes(rclContext, rcl_node)
	if err != nil {
		return "", err
	}

	var i C.int = C.gowrapper_find_rcutils_string_array_index(&rmw_names_and_types.names, C.CString(topic_name), (C.int)(len(topic_name)))
	if i == -1 {
		return "", nil
	}
	var data *C.char = C.gowrapper_get_rcutils_string_array_index(rmw_names_and_types.types, i)
	return C.GoString(data), nil
}

func TopicGetTopicNamesAndTypes(rclContext RCLContext, rcl_node *C.rcl_node_t) (*C.rmw_names_and_types_t, RCLError) {
	re := getRCLEntities(rclContext)
	var rmw_node *C.rmw_node_t = C.rcl_node_get_rmw_handle(rcl_node)

	rmw_names_and_types := (*C.rmw_names_and_types_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rmw_names_and_types_t{}))))
	*rmw_names_and_types = C.rmw_get_zero_initialized_names_and_types() // TODO: Array mnemory handling here

	var rc C.rcl_ret_t = (C.rcl_ret_t)(C.rmw_get_topic_names_and_types(rmw_node, re.rcl_allocator_t, false, rmw_names_and_types)) // rmw_ret_t is aliased to rcl_ret_t
	if rc != 0 {
		return rmw_names_and_types, ErrorsCast(rc)
	}

	return rmw_names_and_types, nil
}

func NewWaitSet(rclContext RCLContext, subscriptions []*Subscription, timers []*Timer, timeout time.Duration) (WaitSet, RCLError) {
	re := getRCLEntities(rclContext)
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
	waitSet.rcl_wait_set_t = &rcl_wait_set
	var rc C.rcl_ret_t = C.rcl_wait_set_init(waitSet.rcl_wait_set_t, number_of_subscriptions, number_of_guard_conditions, number_of_timers, number_of_clients, number_of_services, number_of_events, re.rcl_context_t, *re.rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		return waitSet, ErrorsCast(rc)
	}

	return waitSet, nil
}

/*
WaitSetRun uses wait set manually to be able to better control the parameters to callback handlers.
rclc subscriptions do not pass the rcl_subscription_t to the callback,
making it impossible to dynamically dispatch messages to the corresponding callback handlers
*/
func (self *WaitSet) Run() RCLError {
	for {
		if !C.rcl_wait_set_is_valid(self.rcl_wait_set_t) {
			//#define RCL_RET_WAIT_SET_INVALID 900
			return ErrorsCastC(900, fmt.Sprintf("rcl_wait_set_is_valid() failed for wait_set='%v'", self))
		}
		var rc C.rcl_ret_t = C.rcl_wait_set_clear(self.rcl_wait_set_t)
		if rc != C.RCL_RET_OK {
			return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_clear() failed for wait_set='%v'", self))
		}
		for i := 0; i < len(self.Subscriptions); i++ {
			rc = C.rcl_wait_set_add_subscription(self.rcl_wait_set_t, self.Subscriptions[i].rcl_subscription_t, nil)
			if rc != C.RCL_RET_OK {
				return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_subscription() failed for wait_set='%v'", self))
			}
		}
		for i := 0; i < len(self.Timers); i++ {
			rc = C.rcl_wait_set_add_timer(self.rcl_wait_set_t, self.Timers[i].rcl_timer_t, nil)
			if rc != C.RCL_RET_OK {
				return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_timer() failed for wait_set='%v'", self))
			}
		}
		/*
			rcc := make(chan C.int)
			go func () {
				rc = C.rcl_wait(self.rcl_wait_set_t, (C.long)(self.Timeout))
				if rc == C.RCL_RET_TIMEOUT {
					continue
				}
				rcc <- rc
			}()
			select {
			case <- rclContext.Done:
				return nil
			case rc := <- rcc
				// TODO process WaitSet
			}
		*/
		var i C.ulong
		// Check timers. Guard against internal state representation mismatch. Due to some software bug the lists of timers could easily get out of sync. AND lead to very very difficult to detect bugs.
		if (int)(self.rcl_wait_set_t.size_of_timers) != len(self.Timers) {
			panic(fmt.Sprintf(
				"Wait set timers count mismatch! rcl_wait_set.size_of_timers='%d' != len(self.Timers)='%d'",
				(int)(self.rcl_wait_set_t.size_of_subscriptions),
				len(self.Subscriptions)))
		}
		for i = 0; i < self.rcl_wait_set_t.size_of_timers; i++ {
			var is_timer_ready_to_call C.bool = false
			timer := self.Timers[i]
			rc = C.rcl_timer_is_ready(timer.rcl_timer_t, &is_timer_ready_to_call)
			if rc != C.RCL_RET_OK {
				return ErrorsCastC(rc, fmt.Sprintf("rcl_timer_is_ready() failed for waitSet='%v', timer='%+v'", self, timer))
			}
			if is_timer_ready_to_call {
				timer.Reset()
				timer.Callback(timer)
			}
		}
		// Check subscriptions. Guard against internal state representation mismatch. Due to some software bug the lists of subscriptions could easily get out of sync. AND lead to very very difficult to detect bugs.
		if (int)(self.rcl_wait_set_t.size_of_subscriptions) != len(self.Subscriptions) {
			panic(fmt.Sprintf(
				"Wait set subscriptions count mismatch! rcl_wait_set.size_of_subscriptions='%d' != len(self.Subscriptions)='%d'",
				(int)(self.rcl_wait_set_t.size_of_subscriptions),
				len(self.Subscriptions)))
		}
		for i = 0; i < self.rcl_wait_set_t.size_of_subscriptions; i++ {
			subscription := self.Subscriptions[i]

			rmw_message_info := (*C.rmw_message_info_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rmw_message_info_t{}))))
			*rmw_message_info = C.rmw_get_zero_initialized_message_info()
			defer C.free(unsafe.Pointer(rmw_message_info))

			ros2_msg_receive_buffer := subscription.Ros2MsgType.PrepareMemory()
			defer subscription.Ros2MsgType.ReleaseMemory(ros2_msg_receive_buffer)

			rc = C.rcl_take(subscription.rcl_subscription_t, ros2_msg_receive_buffer, rmw_message_info, nil)
			if rc != C.RCL_RET_OK {
				return ErrorsCastC(rc, fmt.Sprintf("rcl_take() failed for waitSet='%+v', subscription='%+v'", self, subscription))
			}
			rmwMessageInfo := &RmwMessageInfo{
				SourceTimestamp:   time.Unix(0, int64(rmw_message_info.source_timestamp)),
				ReceivedTimestamp: time.Unix(0, int64(rmw_message_info.received_timestamp)),
				FromIntraProcess:  bool(rmw_message_info.from_intra_process),
			}
			subscription.Callback(subscription, ros2_msg_receive_buffer, rmwMessageInfo)
		}
	}
}
