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
#include <rcl/arguments.h>
#include <rcl/graph.h>
#include <rcl/init.h>
#include <rcl/init_options.h>
#include <rcl/logging.h>
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
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/google/shlex"
	"github.com/hashicorp/go-multierror"
	"github.com/kivilahtio/go-re/v0"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
)

/*
Keeps track of all the C entities initialized, so we can later free them
*/
type rclEntityWrapper struct {
	rcl_allocator_t    *C.rcutils_allocator_t
	rcl_context_t      *C.rcl_context_t
	Clock              *Clock
	rcl_init_options_t *C.rcl_init_options_t
	Nodes              list.List // []*Node
	Publishers         list.List // []*Publisher
	Subscriptions      list.List // []*Subscription
	Timers             list.List // []*Timer
	WaitSets           list.List // []*WaitSet
}

/*
Close frees the allocated memory
*/
func (self *rclEntityWrapper) Close() error {
	var errs error
	var rc C.rcl_ret_t
	rc = C.rcl_init_options_fini(self.rcl_init_options_t)
	if rc != C.RCL_RET_OK {
		errs = multierror.Append(errs, ErrorsCastC(rc, fmt.Sprintf("C.rcl_init_options_fini(%+v)", self.rcl_init_options_t)))
	} else {
		self.rcl_init_options_t = nil
	}
	rc = C.rcl_shutdown(self.rcl_context_t)
	if rc != C.RCL_RET_OK {
		errs = multierror.Append(errs, ErrorsCastC(rc, fmt.Sprintf("C.rcl_shutdown(%+v)", self.rcl_context_t)))
	} else {
		C.free(unsafe.Pointer(self.rcl_context_t))
		self.rcl_context_t = nil
	}
	C.free(unsafe.Pointer(self.rcl_allocator_t))
	self.rcl_allocator_t = nil

	rc = C.rcl_clock_fini(self.Clock.rcl_clock_t)
	if rc != C.RCL_RET_OK {
		errs = multierror.Append(errs, ErrorsCastC(rc, fmt.Sprintf("C.rcl_clock_fini(%+v)", self.Clock.rcl_clock_t)))
	} else {
		self.Clock = nil
	}
	return errs
}

type Clock struct {
	rcl_clock_t *C.rcl_clock_t
}

type Node struct {
	rcl_node_t *C.rcl_node_t
	context    *Context
}

type Publisher struct {
	TopicName               string
	Ros2MsgType             ros2types.ROS2Msg
	node                    *Node
	rcl_publisher_options_t *C.rcl_publisher_options_t
	rcl_publisher_t         *C.rcl_publisher_t
}

type Timer struct {
	rcl_timer_t *C.rcl_timer_t
	Callback    func(*Timer)
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

// ROS2 RCL is configured via CLI arguments, so merge them from different sources. See. http://design.ros2.org/articles/ros_command_line_arguments.html
type RCLArgs struct {
	GoArgs []string
	CArgv  **C.char
	CArgc  C.int
}

/*
NewRCLArgs parses ROS2 RCL commandline arguments from the given rclArgs or from
os.Args. If rclArgs is nil returns a string containing the prepared parameters
in the form the RCL can understand them.

C memory is freed when the RCLArgs-object is GC'd.

Example

    oldOSArgs := os.Args
    defer func() { os.Args = oldOSArgs }()

    os.Args = []string{"--extra0", "args0", "--ros-args", "--log-level", "DEBUG", "--", "--extra1", "args1"}
    rosArgs, err := NewRCLArgs("")
    if err != nil {
        panic(err)
    }
    fmt.Printf("rosArgs: %+v\n", rosArgs.GoArgs) // -> [--extra0 args0 --ros-args --log-level DEBUG -- --extra1 args1]

    rosArgs, err = NewRCLArgs("--log-level INFO")
    if err != nil {
        panic(err)
    }
    fmt.Printf("rosArgs: %+v\n", rosArgs.GoArgs) // -> [--ros-args --log-level INFO]

    os.Args = []string{"--extra0", "args0", "--extra1", "args1"}
    rosArgs, err = NewRCLArgs("")
    if err != nil {
        panic(err)
    }
    fmt.Printf("rosArgs: %+v\n", rosArgs.GoArgs) // -> []

    // Output: rosArgs: [--extra0 args0 --ros-args --log-level DEBUG -- --extra1 args1]
    // rosArgs: [--ros-args --log-level INFO]
    // rosArgs: []
*/
func NewRCLArgs(rclArgs string) (*RCLArgs, error) {
	var goArgs []string
	var err error
	rclArgs = strings.Trim(rclArgs, `'"`)
	if r := re.Mr(rclArgs, `m!--ros-args\s+(.+?)\s*(?:--|$)!`); r.Matches > 0 {
		goArgs, err = shlex.Split(rclArgs)
		if err != nil {
			return nil, &RCL_RET_INVALID_ARGUMENT{rclRetStruct{1003, fmt.Sprintf("%s", err), ""}}
		}
	} else if rclArgs != "" {
		goArgs, err = shlex.Split("--ros-args " + rclArgs)
		if err != nil {
			return nil, &RCL_RET_INVALID_ARGUMENT{rclRetStruct{1003, fmt.Sprintf("%s", err), ""}}
		}
	} else if r := re.Mr(strings.Join(os.Args, " "), `m!--ros-args\s+(.+?)\s*(?:--|$)!`); r.Matches > 0 {
		goArgs = os.Args
	}

	ra := &RCLArgs{GoArgs: goArgs, CArgc: C.int(len(goArgs))}
	if len(goArgs) > 0 {
		// Turn the Golang []string into stone, erm. **C.char
		argc := C.int(len(goArgs))
		argv := (**C.char)(C.malloc((C.size_t)((C.int)(unsafe.Sizeof(uintptr(1))) * argc)))
		for i, arg := range goArgs {
			str := C.CString(arg)
			C.setString(argv, C.int(i), str)
		}
		ra.CArgv = argv
	} else {
		ra.CArgv = nil
	}

	runtime.SetFinalizer(ra, func(obj *RCLArgs) {
		ra.Close()
	})
	return ra, nil
}

/*
Close frees the allocated memory
*/
func (self *RCLArgs) Close() {
	for i := 0; i < (int)(self.CArgc); i++ {
		cIdx := unsafe.Pointer(
			uintptr(unsafe.Pointer(self.CArgv)) + (unsafe.Sizeof(uintptr(1)) * uintptr(i)),
		)
		C.free(cIdx)
	}
	C.free(unsafe.Pointer(self.CArgv))
	C.free(unsafe.Pointer(&self.CArgc))
}

/*
NewRCLArgsMust behaves the same as NewRCLArgs except on error it panic()s!
*/
func NewRCLArgsMust(rclArgs string) *RCLArgs {
	args, err := NewRCLArgs(rclArgs)
	if err != nil {
		panic(err)
	}
	return args
}

func rclInit(rclArgs *RCLArgs) (*rclEntityWrapper, error) {
	var rc C.rcl_ret_t
	if rclArgs == nil {
		rclArgs, _ = NewRCLArgs("")
	}

	rclEntityWrapper := &rclEntityWrapper{}
	/* Instead of receiving the rcl_allocator_t as a golang struct,
	   prepare C memory from heap to receive a copy of the rcl allocator.
	   This way Golang wont mess with the rcl_allocator_t memory location
	   and complaing about nested Golang pointer passed over cgo */
	rclEntityWrapper.rcl_allocator_t = (*C.rcl_allocator_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_allocator_t{}))))
	*rclEntityWrapper.rcl_allocator_t = C.rcl_get_default_allocator()
	// TODO: Free C.free(rclEntityWrapper.rcl_allocator)

	rclInitLogging(rclArgs)

	rclEntityWrapper.rcl_context_t = (*C.rcl_context_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_context_t{}))))
	*rclEntityWrapper.rcl_context_t = C.rcl_get_zero_initialized_context()

	rclEntityWrapper.rcl_init_options_t = (*C.rcl_init_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_init_options_t{}))))
	*rclEntityWrapper.rcl_init_options_t = C.rcl_get_zero_initialized_init_options()
	rc = C.rcl_init_options_init(rclEntityWrapper.rcl_init_options_t, *rclEntityWrapper.rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		rclEntityWrapper.Close()
		return nil, ErrorsCast(rc)
	}

	rc = C.rcl_init(rclArgs.CArgc, rclArgs.CArgv, rclEntityWrapper.rcl_init_options_t, rclEntityWrapper.rcl_context_t)
	if rc != C.RCL_RET_OK {
		rclEntityWrapper.Close()
		return nil, ErrorsCast(rc)
	}

	return rclEntityWrapper, nil
}

func rclInitLogging(rclArgs *RCLArgs) error {
	var rc C.rcl_ret_t
	var allocator C.rcl_allocator_t = C.rcl_get_default_allocator()

	var rcl_arguments C.rcl_arguments_t = C.rcl_get_zero_initialized_arguments()
	rc = C.rcl_parse_arguments(rclArgs.CArgc, rclArgs.CArgv, allocator, &rcl_arguments)
	if rc != C.RCL_RET_OK {
		return ErrorsCastC(rc, "rclInitLogging -> rcl_parse_arguments()")
	}

	rc = C.rcl_logging_configure(&rcl_arguments, &allocator)
	if rc != C.RCL_RET_OK {
		return ErrorsCastC(rc, "rclInitLogging -> rcl_logging_configure()")
	}
	return nil
}

func (c *Context) NewNode(node_name, namespace string) (*Node, error) {
	ns := strings.ReplaceAll(namespace, "/", "")
	ns = strings.ReplaceAll(ns, "-", "")

	rcl_node_options := (*C.rcl_node_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_node_options_t{}))))
	*rcl_node_options = C.rcl_node_get_default_options()

	rcl_node := (*C.rcl_node_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_node_t{}))))
	*rcl_node = C.rcl_get_zero_initialized_node()

	var rc C.rcl_ret_t = C.rcl_node_init(rcl_node, C.CString(node_name), C.CString(ns), c.entities.rcl_context_t, rcl_node_options)
	if rc != C.RCL_RET_OK {
		fmt.Printf("Error '%d' in rcl_node_init\n", (int)(rc))
		return nil, ErrorsCast(rc)
	}

	node := &Node{rcl_node_t: rcl_node, context: c}
	c.entities.Nodes.PushFront(node)
	return node, nil
}

/*
Close frees the allocated memory
*/
func (self *Node) Close() error {
	rc := C.rcl_node_fini(self.rcl_node_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}
	return nil
}

func (self *Node) NewPublisher(topicName string, ros2msg ros2types.ROS2Msg) (*Publisher, error) {
	rcl_publisher := (*C.rcl_publisher_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_publisher_t{}))))
	*rcl_publisher = C.rcl_get_zero_initialized_publisher()

	rcl_publisher_options := (*C.rcl_publisher_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_publisher_options_t{}))))
	*rcl_publisher_options = C.rcl_publisher_get_default_options()
	rcl_publisher_options.qos.reliability = C.RMW_QOS_POLICY_RELIABILITY_SYSTEM_DEFAULT

	err := ValidateTopicName(topicName)
	if err != nil {
		return nil, err
	}

	var rc C.rcl_ret_t = C.rcl_publisher_init(
		rcl_publisher,
		self.rcl_node_t,
		(*C.rosidl_message_type_support_t)(ros2msg.TypeSupport()),
		C.CString(topicName),
		rcl_publisher_options)
	if rc != C.RCL_RET_OK {
		return nil, ErrorsCast(rc)
	}

	publisher := &Publisher{
		TopicName:               topicName,
		Ros2MsgType:             ros2msg,
		node:                    self,
		rcl_publisher_options_t: rcl_publisher_options,
		rcl_publisher_t:         rcl_publisher,
	}

	self.context.entities.Publishers.PushFront(publisher)
	return publisher, nil
}

func (self *Publisher) Publish(ros2msg ros2types.ROS2Msg) error {
	var rc C.rcl_ret_t

	ptr := ros2msg.AsCStruct()
	defer ros2msg.ReleaseMemory(unsafe.Pointer(ptr))

	rc = C.rcl_publish(self.rcl_publisher_t, ptr, nil)
	if rc != C.RCL_RET_OK {
		return ErrorsCastC(rc, fmt.Sprintf("rcl_publish() failed for publisher '%+v'", self))
	}
	return nil
}

/*
Close frees the allocated memory
*/
func (self *Publisher) Close() error {
	rc := C.rcl_publisher_fini(self.rcl_publisher_t, self.node.rcl_node_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}
	return nil
}

func (c *Context) NewClock(clockType Rcl_clock_type_t) (*Clock, error) {
	if clockType == 0 {
		clockType = RCL_ROS_TIME
	}
	rcl_clock := (*C.rcl_clock_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_clock_t{})))) //rcl_clock_init() doc says "This will allocate all necessary internal structures, and initialize variables.". The parameter is invalid if no memory allocated beforehand.
	var rc C.rcl_ret_t = C.rcl_clock_init(uint32(clockType), rcl_clock, c.entities.rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		return nil, ErrorsCast(rc)
	}
	c.entities.Clock = &Clock{rcl_clock_t: rcl_clock}
	return c.entities.Clock, nil
}

/*
Close frees the allocated memory
*/
func (self *Clock) Close() error {
	rc := C.rcl_clock_fini(self.rcl_clock_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}
	return nil
}

func (c *Context) NewTimer(timeout time.Duration, timer_callback func(*Timer)) (*Timer, error) {
	var rc C.rcl_ret_t

	if timeout == 0 {
		timeout = 1000 * time.Millisecond
	}
	timer := &Timer{}
	timer.Callback = timer_callback

	timer.rcl_timer_t = (*C.rcl_timer_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_timer_t{}))))
	*timer.rcl_timer_t = C.rcl_get_zero_initialized_timer()

	if c.entities.Clock == nil {
		var err error
		_, err = c.NewClock(RCL_ROS_TIME) // http://design.ros2.org/articles/clock_and_time.html // It is expected that the default choice of time will be to use the ROSTime source
		if err != nil {
			return timer, ErrorsCastC(C.int(err.(*rclRetStruct).rclRetCode), fmt.Sprintf("Forwarding error from '%s'", err.Error()))
		}
	}

	rc = C.rcl_timer_init(
		timer.rcl_timer_t,
		c.entities.Clock.rcl_clock_t,
		c.entities.rcl_context_t,
		(C.long)(timeout),
		nil,
		*c.entities.rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		return timer, ErrorsCast(rc)
	}

	c.entities.Timers.PushFront(timer)
	return timer, nil
}
func (self *Timer) GetTimeUntilNextCall() (int64, error) {
	var rc C.rcl_ret_t
	time_until_next_call := (*C.int64_t)(C.malloc((C.size_t)(8)))
	defer C.free(unsafe.Pointer(time_until_next_call))

	rc = C.rcl_timer_get_time_until_next_call(self.rcl_timer_t, time_until_next_call)
	if rc != C.RCL_RET_OK {
		return 0, ErrorsCast(rc)
	}
	return int64(*time_until_next_call), nil
}

func (self *Timer) Reset() error {
	var rc C.rcl_ret_t
	rc = C.rcl_timer_reset(self.rcl_timer_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}
	return nil
}

/*
Close frees the allocated memory
*/
func (self *Timer) Close() error {
	rc := C.rcl_timer_fini(self.rcl_timer_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}
	return nil
}

type SubscriptionCallback func(*Subscription)

type Subscription struct {
	TopicName                  string
	Ros2MsgType                ros2types.ROS2Msg
	node                       *Node
	rcl_subscription_t         *C.rcl_subscription_t
	rcl_subscription_options_t *C.rcl_subscription_options_t
	Callback                   SubscriptionCallback
}

func (self *Node) NewSubscription(topic_name string, ros2msg ros2types.ROS2Msg, subscriptionCallback SubscriptionCallback) (*Subscription, error) {
	subscription := &Subscription{}
	subscription.rcl_subscription_t = (*C.rcl_subscription_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_subscription_t{}))))
	*subscription.rcl_subscription_t = C.rcl_get_zero_initialized_subscription()
	subscription.node = self
	subscription.Ros2MsgType = ros2msg
	subscription.TopicName = topic_name
	subscription.Callback = subscriptionCallback

	err := ValidateTopicName(subscription.TopicName)
	if err != nil {
		return subscription, err
	}

	sops := C.rcl_subscription_get_default_options()
	sops.qos.reliability = C.RMW_QOS_POLICY_RELIABILITY_SYSTEM_DEFAULT
	subscription.rcl_subscription_options_t = &sops

	var rc C.rcl_ret_t = C.rcl_subscription_init(
		subscription.rcl_subscription_t,
		self.rcl_node_t,
		(*C.rosidl_message_type_support_t)(ros2msg.TypeSupport()),
		C.CString(topic_name),
		subscription.rcl_subscription_options_t)
	if rc != C.RCL_RET_OK {
		return subscription, ErrorsCastC(rc, fmt.Sprintf("Topic name '%s'", topic_name))
	}

	self.context.entities.Subscriptions.PushFront(subscription)
	return subscription, nil
}

func (s *Subscription) TakeMessage(out ros2types.ROS2Msg) (*RmwMessageInfo, error) {
	rmw_message_info := (*C.rmw_message_info_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rmw_message_info_t{}))))
	*rmw_message_info = C.rmw_get_zero_initialized_message_info()
	defer C.free(unsafe.Pointer(rmw_message_info))

	ros2_msg_receive_buffer := s.Ros2MsgType.PrepareMemory()
	defer s.Ros2MsgType.ReleaseMemory(ros2_msg_receive_buffer)

	rc := C.rcl_take(s.rcl_subscription_t, ros2_msg_receive_buffer, rmw_message_info, nil)
	if rc != C.RCL_RET_OK {
		return nil, ErrorsCastC(rc, fmt.Sprintf("rcl_take() failed for subscription='%+v'", s))
	}
	out.AsGoStruct(ros2_msg_receive_buffer)
	return &RmwMessageInfo{
		SourceTimestamp:   time.Unix(0, int64(rmw_message_info.source_timestamp)),
		ReceivedTimestamp: time.Unix(0, int64(rmw_message_info.received_timestamp)),
		FromIntraProcess:  bool(rmw_message_info.from_intra_process),
	}, nil
}

func spinErr(thing string, err error) error {
	return fmt.Errorf("failed to spin %s: %w", thing, err)
}

func (s *Subscription) Spin(ctx context.Context, timeout time.Duration) error {
	ws, err := s.node.context.NewWaitSet(timeout)
	if err != nil {
		return spinErr("subscription", err)
	}
	defer ws.Close()
	ws.AddSubscriptions(s)
	if err = ws.Run(ctx); err != nil {
		return spinErr("subscription", err)
	}
	return nil
}

/*
Close frees the allocated memory
*/
func (self *Subscription) Close() error {
	rc := C.rcl_subscription_fini(self.rcl_subscription_t, self.node.rcl_node_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}
	return nil
}

/*
func PublishersInfoByTopic(rclContext RCLContext, rcl_node *C.rcl_node_t, topic_name string) (*C.rmw_topic_endpoint_info_array_t, error) {
	re := GetRCLEntities(rclContext)
	//TODO: This is actually an array of arrays and the memory allocation mechanisms inside ROS2 rcl are more complex! Need to review this on what to do here.
	rmw_topic_endpoint_info_array := (*C.rmw_topic_endpoint_info_array_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rmw_topic_endpoint_info_array_t{}))))
	*rmw_topic_endpoint_info_array = C.rcl_get_zero_initialized_topic_endpoint_info_array()
	var rc C.rcl_ret_t = C.rcl_get_publishers_info_by_topic(rcl_node, re.rcl_allocator_t, C.CString(topic_name), false, rmw_topic_endpoint_info_array)
	if rc != C.RCL_RET_OK {
		return rmw_topic_endpoint_info_array, ErrorsCast(rc)
	}
	return rmw_topic_endpoint_info_array, nil
}

func TopicGetEndpointInfo(rclContext RCLContext, rcl_node *C.rcl_node_t, topic_name string) error {
	//rmw_topic_endpoint_info_array, err := PublishersInfoByTopic(rclContext, rcl_node, topic_name)
	//if err != nil {
	//	return err
	//}

	//var rmw_topic_endpoint_info C.rmw_topic_endpoint_info_t = C.gowrapper_get_rmw_topic_endpoint_info(rmw_topic_endpoint_info_array, 0)
	//rmw_topic_endpoint_info.
	return nil
}

/*func TopicGetTopicTypeSupport(rclContext RclContext, rcl_node *C.rcl_node_t, topic_name string) (C.rosidl_message_type_support_t, error) {
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
}

func TopicGetTopicTypeString(rclContext RCLContext, rcl_node *C.rcl_node_t, topic_name string) (string, error) {
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

func TopicGetTopicNamesAndTypes(rclContext RCLContext, rcl_node *C.rcl_node_t) (*C.rmw_names_and_types_t, error) {
	re := GetRCLEntities(rclContext)
	var rmw_node *C.rmw_node_t = C.rcl_node_get_rmw_handle(rcl_node)

	rmw_names_and_types := (*C.rmw_names_and_types_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rmw_names_and_types_t{}))))
	*rmw_names_and_types = C.rmw_get_zero_initialized_names_and_types() // TODO: Array mnemory handling here

	var rc C.rcl_ret_t = (C.rcl_ret_t)(C.rmw_get_topic_names_and_types(rmw_node, re.rcl_allocator_t, false, rmw_names_and_types)) // rmw_ret_t is aliased to rcl_ret_t
	if rc != 0 {
		return rmw_names_and_types, ErrorsCast(rc)
	}

	return rmw_names_and_types, nil
}
*/

type WaitSet struct {
	Timeout        time.Duration
	Subscriptions  []*Subscription
	Timers         []*Timer
	Ready          bool // flag to notify the outside gothreads that this WaitSet is ready to receive messages. Use waitSet.WaitForReady() to synchronize
	rcl_wait_set_t C.rcl_wait_set_t
	context        *Context
}

func (c *Context) NewWaitSet(timeout time.Duration) (*WaitSet, error) {
	const (
		subscriptionsCount   = 0
		guardConditionsCount = 0
		timersCount          = 0
		clientsCount         = 0
		servicesCount        = 0
		eventsCount          = 0
	)
	waitSet := &WaitSet{
		context:        c,
		Timeout:        timeout,
		Subscriptions:  []*Subscription{},
		Timers:         []*Timer{},
		rcl_wait_set_t: C.rcl_get_zero_initialized_wait_set(),
	}
	var rc C.rcl_ret_t = C.rcl_wait_set_init(
		&waitSet.rcl_wait_set_t,
		subscriptionsCount,
		guardConditionsCount,
		timersCount,
		clientsCount,
		servicesCount,
		eventsCount,
		c.entities.rcl_context_t,
		*c.entities.rcl_allocator_t,
	)
	if rc != C.RCL_RET_OK {
		return nil, ErrorsCast(rc)
	}
	c.entities.WaitSets.PushFront(waitSet)
	return waitSet, nil
}

func (w *WaitSet) AddSubscriptions(subs ...*Subscription) {
	w.Subscriptions = append(w.Subscriptions, subs...)
}

func (w *WaitSet) AddTimers(timers ...*Timer) {
	w.Timers = append(w.Timers, timers...)
}

func (self *WaitSet) WaitForReady(timeout, interval time.Duration) error {
	for !self.Ready {
		select {
		case <-time.After(timeout):
			if self.Ready {
				return nil
			} else {
				return ErrorsCast(2)
			}
		default:
			time.Sleep(interval)
		}
	}
	return nil
}

func (self *WaitSet) RunGoroutine(ctx context.Context) {
	self.context.WG.Add(1)
	go func() {
		defer self.context.WG.Done()
		if err := self.Run(ctx); err != nil {
			fmt.Printf("RunGoroutine error: '%+v'\n", err)
		}
	}()
}

/*
Run causes the current goroutine to block on this given WaitSet.
WaitSet executes the given timers and subscriptions and calls their callbacks on new events.
*/
func (self *WaitSet) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := self.initEntities()
			if err != nil {
				return err
			}

			var rc C.rcl_ret_t = C.rcl_wait(&self.rcl_wait_set_t, (C.long)(self.Timeout))
			self.Ready = true
			if rc == C.RCL_RET_TIMEOUT {
				continue
			}

			// Check if counts in rcl layer and Go layer differ. Guards against
			// internal state representation mismatch. Due to some software bug
			// the lists of waited resources could easily get out of sync. AND
			// lead to very very difficult to detect bugs.
			panicIfCountMismatch("timers", self.rcl_wait_set_t.size_of_timers, len(self.Timers))
			panicIfCountMismatch("subscriptions", self.rcl_wait_set_t.size_of_subscriptions, len(self.Subscriptions))

			timers := (*[1 << 30]*C.struct_rcl_timer_t)(unsafe.Pointer(self.rcl_wait_set_t.timers))
			for i, t := range self.Timers {
				if timers[i] != nil {
					t.Reset()
					t.Callback(t)
				}
			}
			subs := (*[1 << 30]*C.struct_rcl_subscription_t)(unsafe.Pointer(self.rcl_wait_set_t.subscriptions))
			for i, s := range self.Subscriptions {
				if subs[i] != nil {
					s.Callback(s)
				}
			}
		}
	}
}

func panicIfCountMismatch(typ string, expected C.ulong, actual int) {
	if int(expected) != actual {
		panic(fmt.Sprintf(
			"Wait set %s count mismatch! expected='%d' != actual='%d'",
			typ,
			expected,
			actual))
	}
}

func (self *WaitSet) initEntities() error {
	if !C.rcl_wait_set_is_valid(&self.rcl_wait_set_t) {
		//#define RCL_RET_WAIT_SET_INVALID 900
		return ErrorsCastC(900, fmt.Sprintf("rcl_wait_set_is_valid() failed for wait_set='%v'", self))
	}
	var rc C.rcl_ret_t = C.rcl_wait_set_clear(&self.rcl_wait_set_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_clear() failed for wait_set='%v'", self))
	}
	rc = C.rcl_wait_set_resize(
		&self.rcl_wait_set_t,
		C.size_t(len(self.Subscriptions)),
		self.rcl_wait_set_t.size_of_guard_conditions,
		C.size_t(len(self.Timers)),
		self.rcl_wait_set_t.size_of_clients,
		self.rcl_wait_set_t.size_of_services,
		self.rcl_wait_set_t.size_of_events,
	)
	if rc != C.RCL_RET_OK {
		return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_resize() failed for wait_set='%v'", self))
	}
	for _, sub := range self.Subscriptions {
		rc = C.rcl_wait_set_add_subscription(&self.rcl_wait_set_t, sub.rcl_subscription_t, nil)
		if rc != C.RCL_RET_OK {
			return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_subscription() failed for wait_set='%v'", self))
		}
	}
	for _, timer := range self.Timers {
		rc = C.rcl_wait_set_add_timer(&self.rcl_wait_set_t, timer.rcl_timer_t, nil)
		if rc != C.RCL_RET_OK {
			return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_timer() failed for wait_set='%v'", self))
		}
	}
	return nil
}

/*
Close frees the allocated memory
*/
func (self *WaitSet) Close() error {
	rc := C.rcl_wait_set_fini(&self.rcl_wait_set_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}
	return nil
}
