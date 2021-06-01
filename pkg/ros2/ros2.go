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
#include <rcl/service.h>
#include <rcl/client.h>
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
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/google/shlex"
	"github.com/hashicorp/go-multierror"
	"github.com/kivilahtio/go-re/v0"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
)

type Clock struct {
	rosID
	rcl_clock_t *C.rcl_clock_t
	context     *Context
}

type Node struct {
	rosID
	rosResourceStore
	rcl_node_t *C.rcl_node_t
	context    *Context
}

type Publisher struct {
	rosID
	TopicName               string
	Ros2MsgType             ros2types.ROS2Msg
	node                    *Node
	rcl_publisher_options_t *C.rcl_publisher_options_t
	rcl_publisher_t         *C.rcl_publisher_t
}

type Timer struct {
	rosID
	rcl_timer_t *C.rcl_timer_t
	Callback    func(*Timer)
	context     *Context
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

func rclInit(rclArgs *RCLArgs, ctx *Context) error {
	var rc C.rcl_ret_t
	if rclArgs == nil {
		rclArgs, _ = NewRCLArgs("")
	}

	/* Instead of receiving the rcl_allocator_t as a golang struct,
	   prepare C memory from heap to receive a copy of the rcl allocator.
	   This way Golang wont mess with the rcl_allocator_t memory location
	   and complaing about nested Golang pointer passed over cgo */
	ctx.rcl_allocator_t = (*C.rcl_allocator_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_allocator_t{}))))
	*ctx.rcl_allocator_t = C.rcl_get_default_allocator()
	// TODO: Free C.free(ctx.rcl_allocator)

	rclInitLogging(rclArgs)

	ctx.rcl_context_t = (*C.rcl_context_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_context_t{}))))
	*ctx.rcl_context_t = C.rcl_get_zero_initialized_context()

	ctx.rcl_init_options_t = (*C.rcl_init_options_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_init_options_t{}))))
	*ctx.rcl_init_options_t = C.rcl_get_zero_initialized_init_options()
	rc = C.rcl_init_options_init(ctx.rcl_init_options_t, *ctx.rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}

	rc = C.rcl_init(rclArgs.CArgc, rclArgs.CArgv, ctx.rcl_init_options_t, ctx.rcl_context_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}
	return nil
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

	var rc C.rcl_ret_t = C.rcl_node_init(rcl_node, C.CString(node_name), C.CString(ns), c.rcl_context_t, rcl_node_options)
	if rc != C.RCL_RET_OK {
		fmt.Printf("Error '%d' in rcl_node_init\n", (int)(rc))
		return nil, ErrorsCast(rc)
	}

	node := &Node{rcl_node_t: rcl_node, context: c}
	c.addResource(node)
	return node, nil
}

/*
Close frees the allocated memory
*/
func (self *Node) Close() error {
	if self.rcl_node_t == nil {
		return closeErr("node")
	}
	var err *multierror.Error
	err = multierror.Append(err, self.rosResourceStore.Close())
	self.context.removeResource(self)
	rc := C.rcl_node_fini(self.rcl_node_t)
	if rc != C.RCL_RET_OK {
		err = multierror.Append(err, ErrorsCast(rc))
	}
	C.free(unsafe.Pointer(self.rcl_node_t))
	self.rcl_node_t = nil
	return err.ErrorOrNil()
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

	self.addResource(publisher)
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
	self.node.removeResource(self)
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
	var rc C.rcl_ret_t = C.rcl_clock_init(uint32(clockType), rcl_clock, c.rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		return nil, ErrorsCast(rc)
	}
	clock := &Clock{
		context:     c,
		rcl_clock_t: rcl_clock,
	}
	c.addResource(clock)
	return clock, nil
}

/*
Close frees the allocated memory
*/
func (self *Clock) Close() error {
	self.context.removeResource(self)
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
	timer := &Timer{context: c}
	timer.Callback = timer_callback

	timer.rcl_timer_t = (*C.rcl_timer_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rcl_timer_t{}))))
	*timer.rcl_timer_t = C.rcl_get_zero_initialized_timer()

	if c.Clock == nil {
		var err error
		c.Clock, err = c.NewClock(RCL_ROS_TIME) // http://design.ros2.org/articles/clock_and_time.html // It is expected that the default choice of time will be to use the ROSTime source
		if err != nil {
			return nil, err
		}
	}

	rc = C.rcl_timer_init(
		timer.rcl_timer_t,
		c.Clock.rcl_clock_t,
		c.rcl_context_t,
		(C.long)(timeout),
		nil,
		*c.rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		return nil, ErrorsCast(rc)
	}

	c.addResource(timer)
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
	self.context.removeResource(self)
	rc := C.rcl_timer_fini(self.rcl_timer_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}
	C.free(unsafe.Pointer(self.rcl_timer_t))
	return nil
}

type SubscriptionCallback func(*Subscription)

type Subscription struct {
	rosID
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

	self.addResource(subscription)
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
	self.node.removeResource(self)
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
	rosID
	Timeout        time.Duration
	Subscriptions  []*Subscription
	Timers         []*Timer
	services       []*Service
	clients        []*Client
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
		services:       []*Service{},
		clients:        []*Client{},
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
		c.rcl_context_t,
		*c.rcl_allocator_t,
	)
	if rc != C.RCL_RET_OK {
		return nil, ErrorsCast(rc)
	}

	c.addResource(waitSet)
	return waitSet, nil
}

func (w *WaitSet) AddSubscriptions(subs ...*Subscription) {
	w.Subscriptions = append(w.Subscriptions, subs...)
}

func (w *WaitSet) AddTimers(timers ...*Timer) {
	w.Timers = append(w.Timers, timers...)
}

func (w *WaitSet) AddServices(services ...*Service) {
	w.services = append(w.services, services...)
}

func (w *WaitSet) AddClients(clients ...*Client) {
	w.clients = append(w.clients, clients...)
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
			panicIfCountMismatch("services", self.rcl_wait_set_t.size_of_services, len(self.services))
			panicIfCountMismatch("clients", self.rcl_wait_set_t.size_of_clients, len(self.clients))

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
			svcs := (*[1 << 30]*C.struct_rcl_service_t)(unsafe.Pointer(self.rcl_wait_set_t.services))
			for i, s := range self.services {
				if svcs[i] != nil {
					s.handleRequest()
				}
			}
			clients := (*[1 << 30]*C.struct_rcl_client_t)(unsafe.Pointer(self.rcl_wait_set_t.clients))
			for i, c := range self.clients {
				if clients[i] != nil {
					c.handleResponse()
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
		C.size_t(len(self.clients)),
		C.size_t(len(self.services)),
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
	for _, service := range self.services {
		rc = C.rcl_wait_set_add_service(&self.rcl_wait_set_t, service.rclService, nil)
		if rc != C.RCL_RET_OK {
			return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_service() failed for wait_set='%v'", self))
		}
	}
	for _, client := range self.clients {
		rc = C.rcl_wait_set_add_client(&self.rcl_wait_set_t, client.rclClient, nil)
		if rc != C.RCL_RET_OK {
			return ErrorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_client() failed for wait_set='%v'", self))
		}
	}
	return nil
}

/*
Close frees the allocated memory
*/
func (self *WaitSet) Close() error {
	if self.context == nil {
		return closeErr("wait set")
	}
	self.context.removeResource(self)
	self.context = nil
	rc := C.rcl_wait_set_fini(&self.rcl_wait_set_t)
	if rc != C.RCL_RET_OK {
		return ErrorsCast(rc)
	}
	return nil
}

type RmwRequestID struct {
	WriterGUID     [16]int8
	SequenceNumber int64
}

type RmwServiceInfo struct {
	SourceTimestamp   time.Time
	ReceivedTimestamp time.Time
	RequestID         RmwRequestID
}

type ServiceOptions struct {
	Qos RmwQosProfile
}

func NewDefaultServiceOptions() *ServiceOptions {
	return &ServiceOptions{Qos: NewRmwQosProfileServicesDefault()}
}

type ServiceResponseSender interface {
	SendResponse(resp ros2types.ROS2Msg) error
}

type serviceResponseSender func(resp ros2types.ROS2Msg) error

func (s serviceResponseSender) SendResponse(resp ros2types.ROS2Msg) error {
	return s(resp)
}

type ServiceRequestHandler func(*RmwServiceInfo, ros2types.ROS2Msg, ServiceResponseSender)

type Service struct {
	rosID
	node        *Node
	rclService  *C.rcl_service_t
	name        *C.char
	handler     ServiceRequestHandler
	typeSupport ros2types.Service
}

// NewService creates a new service.
//
// options must not be modified after passing it to this function. If options is
// nil, default options are used.
func (n *Node) NewService(
	name string,
	typeSupport ros2types.Service,
	options *ServiceOptions,
	handler ServiceRequestHandler,
) (s *Service, err error) {
	if options == nil {
		options = NewDefaultServiceOptions()
	}
	s = &Service{
		node:        n,
		rclService:  (*C.rcl_service_t)(C.malloc(C.sizeof_struct_rcl_service_t)),
		name:        C.CString(name),
		handler:     handler,
		typeSupport: typeSupport,
	}
	defer onErr(&err, s.Close)
	*s.rclService = C.rcl_get_zero_initialized_service()
	opts := C.rcl_service_options_t{allocator: *n.context.rcl_allocator_t}
	options.Qos.asCStruct(&opts.qos)
	retCode := C.rcl_service_init(
		s.rclService,
		n.rcl_node_t,
		(*C.struct_rosidl_service_type_support_t)(typeSupport.TypeSupport()),
		s.name,
		&opts,
	)
	if retCode != C.RCL_RET_OK {
		return nil, ErrorsCastC(retCode, "failed to create service")
	}
	n.addResource(s)
	return s, nil
}

func (s *Service) Close() (err error) {
	if s.name == nil {
		return closeErr("service")
	}
	s.node.removeResource(s)
	rc := C.rcl_service_fini(s.rclService, s.node.rcl_node_t)
	if rc != C.RCL_RET_OK {
		err = ErrorsCastC(rc, "failed to finalize service")
	}
	C.free(unsafe.Pointer(s.name))
	s.name = nil
	return err
}

func (s *Service) handleRequest() {
	var reqHeader C.struct_rmw_service_info_t
	reqBuffer := s.typeSupport.Request().PrepareMemory()
	defer s.typeSupport.Request().ReleaseMemory(reqBuffer)
	rc := C.rcl_take_request_with_info(s.rclService, &reqHeader, reqBuffer)
	if rc != C.RCL_RET_OK {
		log.Println(ErrorsCastC(rc, "failed to take request"))
		return
	}
	info := RmwServiceInfo{
		SourceTimestamp:   time.Unix(0, int64(reqHeader.source_timestamp)),
		ReceivedTimestamp: time.Unix(0, int64(reqHeader.received_timestamp)),
		RequestID: RmwRequestID{
			WriterGUID:     *(*[16]int8)(unsafe.Pointer(&reqHeader.request_id.writer_guid)),
			SequenceNumber: int64(reqHeader.request_id.sequence_number),
		},
	}
	req := s.typeSupport.Request().Clone()
	req.AsGoStruct(reqBuffer)
	s.handler(
		&info,
		req,
		serviceResponseSender(func(resp ros2types.ROS2Msg) error {
			respBuffer := resp.AsCStruct()
			defer resp.ReleaseMemory(respBuffer)
			rc := C.rcl_send_response(s.rclService, &reqHeader.request_id, respBuffer)
			if rc != C.RCL_RET_OK {
				return ErrorsCastC(rc, "failed to send response")
			}
			return nil
		}),
	)
}

type ClientOptions struct {
	Qos RmwQosProfile
}

func NewDefaultClientOptions() *ClientOptions {
	return &ClientOptions{Qos: NewRmwQosProfileServicesDefault()}
}

// Client is used to send requests to and receive responses from a service.
//
// Calling Send is thread-safe. Creating and closing clients is not thread-safe.
type Client struct {
	rosID
	node                 *Node
	rclClient            *C.struct_rcl_client_t
	serviceName          *C.char
	pendingRequests      map[C.long]chan *clientSendResult
	pendingRequestsMutex sync.Mutex
	typeSupport          ros2types.Service
}

// NewClient creates a new client.
//
// options must not be modified after passing it to this function. If options is
// nil, default options are used.
func (n *Node) NewClient(
	serviceName string,
	typeSupport ros2types.Service,
	options *ClientOptions,
) (c *Client, err error) {
	if options == nil {
		options = NewDefaultClientOptions()
	}
	c = &Client{
		node:            n,
		rclClient:       (*C.struct_rcl_client_t)(C.malloc(C.sizeof_struct_rcl_client_t)),
		serviceName:     C.CString(serviceName),
		pendingRequests: make(map[C.long]chan *clientSendResult),
		typeSupport:     typeSupport,
	}
	defer onErr(&err, c.Close)
	*c.rclClient = C.rcl_get_zero_initialized_client()
	opts := C.struct_rcl_client_options_t{allocator: *n.context.rcl_allocator_t}
	options.Qos.asCStruct(&opts.qos)
	rc := C.rcl_client_init(
		c.rclClient,
		n.rcl_node_t,
		(*C.struct_rosidl_service_type_support_t)(typeSupport.TypeSupport()),
		c.serviceName,
		&opts,
	)
	if rc != C.RCL_RET_OK {
		return nil, ErrorsCastC(rc, "failed to create client")
	}
	n.addResource(c)
	return c, nil
}

func (c *Client) Close() error {
	if c.rclClient == nil {
		return closeErr("client")
	}
	var err *multierror.Error
	c.node.removeResource(c)
	func() {
		c.pendingRequestsMutex.Lock()
		defer c.pendingRequestsMutex.Unlock()
		for _, resultChan := range c.pendingRequests {
			close(resultChan)
		}
	}()

	rc := C.rcl_client_fini(c.rclClient, c.node.rcl_node_t)
	if rc != C.RCL_RET_OK {
		err = multierror.Append(err, ErrorsCastC(rc, "failed to finalize client"))
	}
	c.rclClient = nil

	C.free(unsafe.Pointer(c.serviceName))

	return err.ErrorOrNil()
}

func (c *Client) Send(ctx context.Context, req ros2types.ROS2Msg) (ros2types.ROS2Msg, *RmwServiceInfo, error) {
	resultChan, err := c.sendRequest(req)
	if err != nil {
		return nil, nil, err
	}
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	case result := <-resultChan:
		if result == nil {
			return nil, nil, errors.New("client was closed before a response was received")
		}
		return result.resp, result.info, nil
	}
}

func (c *Client) sendRequest(req ros2types.ROS2Msg) (chan *clientSendResult, error) {
	c.pendingRequestsMutex.Lock()
	defer c.pendingRequestsMutex.Unlock()

	reqBuf := req.AsCStruct()
	defer req.ReleaseMemory(reqBuf)

	var sequenceNumber C.long
	rc := C.rcl_send_request(c.rclClient, reqBuf, &sequenceNumber)
	if rc != C.RCL_RET_OK {
		return nil, ErrorsCastC(rc, "failed to send request")
	}
	resultChan := make(chan *clientSendResult, 1)
	c.pendingRequests[sequenceNumber] = resultChan

	return resultChan, nil
}

func (c *Client) handleResponse() {
	c.pendingRequestsMutex.Lock()
	defer c.pendingRequestsMutex.Unlock()

	var respHeader C.struct_rmw_service_info_t
	respBuf := c.typeSupport.Response().PrepareMemory()
	defer c.typeSupport.Response().ReleaseMemory(respBuf)
	rc := C.rcl_take_response_with_info(c.rclClient, &respHeader, respBuf)
	if rc != C.RCL_RET_OK {
		log.Println(ErrorsCastC(rc, "failed to take response"))
		return
	}
	defer delete(c.pendingRequests, respHeader.request_id.sequence_number)
	defer close(c.pendingRequests[respHeader.request_id.sequence_number])

	result := &clientSendResult{
		resp: c.typeSupport.Response().Clone(),
		info: &RmwServiceInfo{
			SourceTimestamp:   time.Unix(0, int64(respHeader.source_timestamp)),
			ReceivedTimestamp: time.Unix(0, int64(respHeader.received_timestamp)),
			RequestID: RmwRequestID{
				WriterGUID:     *(*[16]int8)(unsafe.Pointer(&respHeader.request_id.writer_guid)),
				SequenceNumber: int64(respHeader.request_id.sequence_number),
			},
		},
	}
	result.resp.AsGoStruct(respBuf)

	c.pendingRequests[respHeader.request_id.sequence_number] <- result
}

type clientSendResult struct {
	resp ros2types.ROS2Msg
	info *RmwServiceInfo
}
