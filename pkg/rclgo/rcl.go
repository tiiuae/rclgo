/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package rclgo

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrmw -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
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
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/google/shlex"
	"github.com/hashicorp/go-multierror"
	"github.com/kivilahtio/go-re/v0"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
)

type RmwMessageInfo struct {
	SourceTimestamp   time.Time
	ReceivedTimestamp time.Time
	FromIntraProcess  bool
}

type ClockType uint32

const (
	ClockTypeUninitialized ClockType = 0
	ClockTypeROSTime       ClockType = 1
	ClockTypeSystemTime    ClockType = 2
	ClockTypeSteadyTime    ClockType = 3
)

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
			return nil, &InvalidArgument{rclRetStruct{1003, fmt.Sprintf("%s", err), ""}}
		}
	} else if rclArgs != "" {
		goArgs, err = shlex.Split("--ros-args " + rclArgs)
		if err != nil {
			return nil, &InvalidArgument{rclRetStruct{1003, fmt.Sprintf("%s", err), ""}}
		}
	} else if r := re.Mr(strings.Join(os.Args, " "), `m!--ros-args\s+(.+?)\s*(?:--|$)!`); r.Matches > 0 {
		goArgs = os.Args
	}

	ra := &RCLArgs{GoArgs: goArgs, CArgc: C.int(len(goArgs))}
	if len(goArgs) > 0 {
		// Turn the Golang []string into stone, erm. **C.char
		argc := C.size_t(len(goArgs))
		ra.CArgv = (**C.char)(C.malloc(C.size_t(unsafe.Sizeof(uintptr(1))) * argc))
		cargv := (*[1 << 30]*C.char)(unsafe.Pointer(ra.CArgv))
		for i, arg := range goArgs {
			cargv[i] = C.CString(arg)
		}
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
func (self *RCLArgs) Close() error {
	if self.CArgv == nil {
		return closeErr("rcl args")
	}
	cargv := (*[1 << 30]*C.char)(unsafe.Pointer(self.CArgv))
	for i := 0; i < int(self.CArgc); i++ {
		C.free(unsafe.Pointer(cargv[i]))
	}
	C.free(unsafe.Pointer(self.CArgv))
	self.CArgv = nil
	return nil
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

	ctx.rcl_allocator_t = (*C.rcl_allocator_t)(C.malloc(C.sizeof_rcl_allocator_t))
	*ctx.rcl_allocator_t = C.rcl_get_default_allocator()
	ctx.rcl_context_t = (*C.rcl_context_t)(C.malloc(C.sizeof_rcl_context_t))
	*ctx.rcl_context_t = C.rcl_get_zero_initialized_context()

	rclInitLogging(rclArgs, ctx.rcl_allocator_t)

	rcl_init_options_t := C.rcl_get_zero_initialized_init_options()
	rc = C.rcl_init_options_init(&rcl_init_options_t, *ctx.rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		return errorsCast(rc)
	}

	rc = C.rcl_init(rclArgs.CArgc, rclArgs.CArgv, &rcl_init_options_t, ctx.rcl_context_t)
	if rc != C.RCL_RET_OK {
		return errorsCast(rc)
	}
	return nil
}

func rclInitLogging(rclArgs *RCLArgs, allocator *C.rcl_allocator_t) error {
	rcl_arguments := C.rcl_get_zero_initialized_arguments()
	rc := C.rcl_parse_arguments(rclArgs.CArgc, rclArgs.CArgv, *allocator, &rcl_arguments)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, "rclInitLogging -> rcl_parse_arguments()")
	}

	rc = C.rcl_logging_configure(&rcl_arguments, allocator)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, "rclInitLogging -> rcl_logging_configure()")
	}
	return nil
}

type Node struct {
	rosID
	rosResourceStore
	rcl_node_t      *C.rcl_node_t
	context         *Context
	name, namespace *C.char
}

func (c *Context) NewNode(node_name, namespace string) (node *Node, err error) {
	node = &Node{
		rcl_node_t: (*C.rcl_node_t)(C.malloc(C.sizeof_rcl_node_t)),
		context:    c,
		name:       C.CString(node_name),
		namespace:  C.CString(namespace),
	}
	defer onErr(&err, node.Close)
	*node.rcl_node_t = C.rcl_get_zero_initialized_node()

	rcl_node_options := C.rcl_node_get_default_options()
	rc := C.rcl_node_init(
		node.rcl_node_t,
		node.name,
		node.namespace,
		c.rcl_context_t,
		&rcl_node_options,
	)
	if rc != C.RCL_RET_OK {
		return nil, errorsCastC(rc, "failed to create node:")
	}

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
	self.context.removeResource(self)

	var err *multierror.Error
	err = multierror.Append(err, self.rosResourceStore.Close())

	rc := C.rcl_node_fini(self.rcl_node_t)
	if rc != C.RCL_RET_OK {
		err = multierror.Append(err, errorsCast(rc))
	}
	C.free(unsafe.Pointer(self.rcl_node_t))
	self.rcl_node_t = nil

	C.free(unsafe.Pointer(self.name))
	C.free(unsafe.Pointer(self.namespace))

	return err.ErrorOrNil()
}

type PublisherOptions struct {
	Qos RmwQosProfile
}

func NewDefaultPublisherOptions() *PublisherOptions {
	return &PublisherOptions{Qos: NewRmwQosProfileDefault()}
}

type Publisher struct {
	rosID
	TopicName       string
	typeSupport     types.MessageTypeSupport
	node            *Node
	rcl_publisher_t *C.rcl_publisher_t
	topicName       *C.char
}

// NewPublisher creates a new publisher.
//
// options must not be modified after passing it to this function. If options is
// nil, default options are used.
func (self *Node) NewPublisher(
	topicName string,
	ros2msg types.MessageTypeSupport,
	options *PublisherOptions,
) (pub *Publisher, err error) {
	if options == nil {
		options = NewDefaultPublisherOptions()
	}
	pub = &Publisher{
		TopicName:       topicName,
		typeSupport:     ros2msg,
		node:            self,
		rcl_publisher_t: (*C.rcl_publisher_t)(C.malloc(C.sizeof_rcl_publisher_t)),
		topicName:       C.CString(topicName),
	}
	defer onErr(&err, pub.Close)
	*pub.rcl_publisher_t = C.rcl_get_zero_initialized_publisher()
	rcl_publisher_options_t := C.rcl_publisher_get_default_options()
	rcl_publisher_options_t.allocator = *self.context.rcl_allocator_t
	options.Qos.asCStruct(&rcl_publisher_options_t.qos)

	var rc C.rcl_ret_t = C.rcl_publisher_init(
		pub.rcl_publisher_t,
		self.rcl_node_t,
		(*C.rosidl_message_type_support_t)(ros2msg.TypeSupport()),
		pub.topicName,
		&rcl_publisher_options_t,
	)
	if rc != C.RCL_RET_OK {
		return nil, errorsCast(rc)
	}

	self.addResource(pub)
	return pub, nil
}

func (self *Publisher) Publish(ros2msg types.Message) error {
	var rc C.rcl_ret_t

	ptr := self.typeSupport.PrepareMemory()
	defer self.typeSupport.ReleaseMemory(ptr)
	self.typeSupport.AsCStruct(ptr, ros2msg)

	rc = C.rcl_publish(self.rcl_publisher_t, ptr, nil)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, fmt.Sprintf("rcl_publish() failed for publisher '%+v'", self))
	}
	return nil
}

/*
Close frees the allocated memory
*/
func (self *Publisher) Close() error {
	if self.rcl_publisher_t == nil {
		return closeErr("publisher")
	}
	var err *multierror.Error
	self.node.removeResource(self)
	rc := C.rcl_publisher_fini(self.rcl_publisher_t, self.node.rcl_node_t)
	if rc != C.RCL_RET_OK {
		err = multierror.Append(err, errorsCast(rc))
	}
	C.free(unsafe.Pointer(self.rcl_publisher_t))
	self.rcl_publisher_t = nil
	C.free(unsafe.Pointer(self.topicName))
	return err.ErrorOrNil()
}

type Clock struct {
	rosID
	rcl_clock_t *C.rcl_clock_t
	context     *Context
}

func (c *Context) NewClock(clockType ClockType) (clock *Clock, err error) {
	if clockType == ClockTypeUninitialized {
		clockType = ClockTypeROSTime
	}
	clock = &Clock{
		context:     c,
		rcl_clock_t: (*C.rcl_clock_t)(C.malloc(C.sizeof_rcl_clock_t)),
	}
	defer onErr(&err, c.Close)
	rc := C.rcl_clock_init(
		uint32(clockType),
		clock.rcl_clock_t,
		c.rcl_allocator_t,
	)
	if rc != C.RCL_RET_OK {
		return nil, errorsCast(rc)
	}
	c.addResource(clock)
	return clock, nil
}

/*
Close frees the allocated memory
*/
func (self *Clock) Close() error {
	if self.rcl_clock_t == nil {
		return closeErr("clock")
	}
	var err *multierror.Error
	self.context.removeResource(self)
	rc := C.rcl_clock_fini(self.rcl_clock_t)
	if rc != C.RCL_RET_OK {
		err = multierror.Append(err, errorsCast(rc))
	}
	C.free(unsafe.Pointer(self.rcl_clock_t))
	self.rcl_clock_t = nil
	return err.ErrorOrNil()
}

type Timer struct {
	rosID
	rcl_timer_t *C.rcl_timer_t
	Callback    func(*Timer)
	context     *Context
}

func (c *Context) NewTimer(timeout time.Duration, timer_callback func(*Timer)) (timer *Timer, err error) {
	if timeout == 0 {
		timeout = 1000 * time.Millisecond
	}
	timer = &Timer{
		rcl_timer_t: (*C.rcl_timer_t)(C.malloc(C.sizeof_rcl_timer_t)),
		Callback:    timer_callback,
		context:     c,
	}
	defer onErr(&err, timer.Close)
	*timer.rcl_timer_t = C.rcl_get_zero_initialized_timer()

	if c.Clock == nil {
		var err error
		c.Clock, err = c.NewClock(ClockTypeROSTime) // http://design.ros2.org/articles/clock_and_time.html // It is expected that the default choice of time will be to use the ROSTime source
		if err != nil {
			return nil, err
		}
	}

	rc := C.rcl_timer_init(
		timer.rcl_timer_t,
		c.Clock.rcl_clock_t,
		c.rcl_context_t,
		(C.long)(timeout),
		nil,
		*c.rcl_allocator_t,
	)
	if rc != C.RCL_RET_OK {
		return nil, errorsCast(rc)
	}

	c.addResource(timer)
	return timer, nil
}

func (self *Timer) GetTimeUntilNextCall() (int64, error) {
	var time_until_next_call C.int64_t
	rc := C.rcl_timer_get_time_until_next_call(self.rcl_timer_t, &time_until_next_call)
	if rc != C.RCL_RET_OK {
		return 0, errorsCast(rc)
	}
	return int64(time_until_next_call), nil
}

func (self *Timer) Reset() error {
	var rc C.rcl_ret_t
	rc = C.rcl_timer_reset(self.rcl_timer_t)
	if rc != C.RCL_RET_OK {
		return errorsCast(rc)
	}
	return nil
}

/*
Close frees the allocated memory
*/
func (self *Timer) Close() error {
	if self.rcl_timer_t == nil {
		return closeErr("timer")
	}
	var err *multierror.Error
	self.context.removeResource(self)
	rc := C.rcl_timer_fini(self.rcl_timer_t)
	if rc != C.RCL_RET_OK {
		err = multierror.Append(err, errorsCast(rc))
	}
	C.free(unsafe.Pointer(self.rcl_timer_t))
	self.rcl_timer_t = nil
	return err.ErrorOrNil()
}

type SubscriptionCallback func(*Subscription)

type Subscription struct {
	rosID
	TopicName          string
	Ros2MsgType        types.MessageTypeSupport
	Callback           SubscriptionCallback
	node               *Node
	rcl_subscription_t *C.rcl_subscription_t
	topicName          *C.char
}

func (self *Node) NewSubscription(topic_name string, ros2msg types.MessageTypeSupport, subscriptionCallback SubscriptionCallback) (sub *Subscription, err error) {
	sub = &Subscription{
		TopicName:          topic_name,
		Ros2MsgType:        ros2msg,
		Callback:           subscriptionCallback,
		node:               self,
		rcl_subscription_t: (*C.rcl_subscription_t)(C.malloc(C.sizeof_rcl_subscription_t)),
		topicName:          C.CString(topic_name),
	}
	defer onErr(&err, sub.Close)
	*sub.rcl_subscription_t = C.rcl_get_zero_initialized_subscription()

	opts := C.rcl_subscription_get_default_options()
	opts.qos.reliability = C.RMW_QOS_POLICY_RELIABILITY_SYSTEM_DEFAULT

	rc := C.rcl_subscription_init(
		sub.rcl_subscription_t,
		self.rcl_node_t,
		(*C.rosidl_message_type_support_t)(ros2msg.TypeSupport()),
		C.CString(topic_name),
		&opts)
	if rc != C.RCL_RET_OK {
		return sub, errorsCastC(rc, fmt.Sprintf("Topic name '%s'", topic_name))
	}

	self.addResource(sub)
	return sub, nil
}

func (s *Subscription) TakeMessage(out types.Message) (*RmwMessageInfo, error) {
	rmw_message_info := C.rmw_get_zero_initialized_message_info()

	ros2_msg_receive_buffer := s.Ros2MsgType.PrepareMemory()
	defer s.Ros2MsgType.ReleaseMemory(ros2_msg_receive_buffer)

	rc := C.rcl_take(s.rcl_subscription_t, ros2_msg_receive_buffer, &rmw_message_info, nil)
	if rc != C.RCL_RET_OK {
		return nil, errorsCastC(rc, fmt.Sprintf("rcl_take() failed for subscription='%+v'", s))
	}
	s.Ros2MsgType.AsGoStruct(out, ros2_msg_receive_buffer)
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
	if self.rcl_subscription_t == nil {
		return closeErr("subscription")
	}
	var err *multierror.Error
	self.node.removeResource(self)
	rc := C.rcl_subscription_fini(self.rcl_subscription_t, self.node.rcl_node_t)
	if rc != C.RCL_RET_OK {
		err = multierror.Append(err, errorsCast(rc))
	}
	C.free(unsafe.Pointer(self.rcl_subscription_t))
	self.rcl_subscription_t = nil
	C.free(unsafe.Pointer(self.topicName))
	return err.ErrorOrNil()
}

/*
func PublishersInfoByTopic(rclContext RCLContext, rcl_node *C.rcl_node_t, topic_name string) (*C.rmw_topic_endpoint_info_array_t, error) {
	re := GetRCLEntities(rclContext)
	//TODO: This is actually an array of arrays and the memory allocation mechanisms inside ROS2 rcl are more complex! Need to review this on what to do here.
	rmw_topic_endpoint_info_array := (*C.rmw_topic_endpoint_info_array_t)(C.malloc((C.size_t)(unsafe.Sizeof(C.rmw_topic_endpoint_info_array_t{}))))
	*rmw_topic_endpoint_info_array = C.rcl_get_zero_initialized_topic_endpoint_info_array()
	var rc C.rcl_ret_t = C.rcl_get_publishers_info_by_topic(rcl_node, re.rcl_allocator_t, C.CString(topic_name), false, rmw_topic_endpoint_info_array)
	if rc != C.RCL_RET_OK {
		return rmw_topic_endpoint_info_array, errorsCast(rc)
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
		return nil, errorsCastC(C.RCL_RET_TOPIC_NAME_INVALID, topic_name)
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
		return rmw_names_and_types, errorsCast(rc)
	}

	return rmw_names_and_types, nil
}
*/

type WaitSet struct {
	rosID
	Timeout       time.Duration
	Subscriptions []*Subscription
	Timers        []*Timer
	services      []*Service
	clients       []*Client
	// Flag to notify the outside gothreads that this WaitSet is ready to
	// receive messages. Use waitSet.WaitForReady() to synchronize. Value of 0
	// means not ready and any other value means ready.
	ready          uint32
	rcl_wait_set_t C.rcl_wait_set_t
	context        *Context
}

func (c *Context) NewWaitSet(timeout time.Duration) (ws *WaitSet, err error) {
	const (
		subscriptionsCount   = 0
		guardConditionsCount = 0
		timersCount          = 0
		clientsCount         = 0
		servicesCount        = 0
		eventsCount          = 0
	)
	ws = &WaitSet{
		context:        c,
		Timeout:        timeout,
		Subscriptions:  []*Subscription{},
		Timers:         []*Timer{},
		services:       []*Service{},
		clients:        []*Client{},
		rcl_wait_set_t: C.rcl_get_zero_initialized_wait_set(),
	}
	defer onErr(&err, ws.Close)
	var rc C.rcl_ret_t = C.rcl_wait_set_init(
		&ws.rcl_wait_set_t,
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
		return nil, errorsCast(rc)
	}

	c.addResource(ws)
	return ws, nil
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

func (w *WaitSet) Ready() bool {
	return atomic.LoadUint32(&w.ready) > 0
}

func (w *WaitSet) markReady() {
	atomic.StoreUint32(&w.ready, 1)
}

func (self *WaitSet) WaitForReady(timeout, interval time.Duration) error {
	after := time.After(timeout)
	for {
		if self.Ready() {
			return nil
		}
		select {
		case <-after:
			if self.Ready() {
				return nil
			}
			return errorsCastC(2, "WaitForReady:")
		default:
			time.Sleep(interval)
		}
	}
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
			self.markReady()
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
		return errorsCastC(900, fmt.Sprintf("rcl_wait_set_is_valid() failed for wait_set='%v'", self))
	}
	var rc C.rcl_ret_t = C.rcl_wait_set_clear(&self.rcl_wait_set_t)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_clear() failed for wait_set='%v'", self))
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
		return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_resize() failed for wait_set='%v'", self))
	}
	for _, sub := range self.Subscriptions {
		rc = C.rcl_wait_set_add_subscription(&self.rcl_wait_set_t, sub.rcl_subscription_t, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_subscription() failed for wait_set='%v'", self))
		}
	}
	for _, timer := range self.Timers {
		rc = C.rcl_wait_set_add_timer(&self.rcl_wait_set_t, timer.rcl_timer_t, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_timer() failed for wait_set='%v'", self))
		}
	}
	for _, service := range self.services {
		rc = C.rcl_wait_set_add_service(&self.rcl_wait_set_t, service.rclService, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_service() failed for wait_set='%v'", self))
		}
	}
	for _, client := range self.clients {
		rc = C.rcl_wait_set_add_client(&self.rcl_wait_set_t, client.rclClient, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_client() failed for wait_set='%v'", self))
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
	var err *multierror.Error
	self.context.removeResource(self)
	self.context = nil
	rc := C.rcl_wait_set_fini(&self.rcl_wait_set_t)
	if rc != C.RCL_RET_OK {
		err = multierror.Append(err, errorsCast(rc))
	}
	return err.ErrorOrNil()
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
	SendResponse(resp types.Message) error
}

type serviceResponseSender func(resp types.Message) error

func (s serviceResponseSender) SendResponse(resp types.Message) error {
	return s(resp)
}

type ServiceRequestHandler func(*RmwServiceInfo, types.Message, ServiceResponseSender)

type Service struct {
	rosID
	node                *Node
	rclService          *C.rcl_service_t
	name                *C.char
	handler             ServiceRequestHandler
	requestTypeSupport  types.MessageTypeSupport
	responseTypeSupport types.MessageTypeSupport
}

// NewService creates a new service.
//
// options must not be modified after passing it to this function. If options is
// nil, default options are used.
func (n *Node) NewService(
	name string,
	typeSupport types.ServiceTypeSupport,
	options *ServiceOptions,
	handler ServiceRequestHandler,
) (s *Service, err error) {
	if options == nil {
		options = NewDefaultServiceOptions()
	}
	s = &Service{
		requestTypeSupport:  typeSupport.Request(),
		responseTypeSupport: typeSupport.Response(),
		node:                n,
		rclService:          (*C.rcl_service_t)(C.malloc(C.sizeof_struct_rcl_service_t)),
		name:                C.CString(name),
		handler:             handler,
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
		return nil, errorsCastC(retCode, "failed to create service")
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
		err = errorsCastC(rc, "failed to finalize service")
	}
	C.free(unsafe.Pointer(s.rclService))
	C.free(unsafe.Pointer(s.name))
	s.name = nil
	return err
}

func (s *Service) handleRequest() {
	var reqHeader C.struct_rmw_service_info_t
	reqBuffer := s.requestTypeSupport.PrepareMemory()
	defer s.requestTypeSupport.ReleaseMemory(reqBuffer)
	rc := C.rcl_take_request_with_info(s.rclService, &reqHeader, reqBuffer)
	if rc != C.RCL_RET_OK {
		log.Println(errorsCastC(rc, "failed to take request"))
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
	req := s.requestTypeSupport.New()
	s.requestTypeSupport.AsGoStruct(req, reqBuffer)
	s.handler(
		&info,
		req,
		serviceResponseSender(func(resp types.Message) error {
			respBuffer := s.responseTypeSupport.PrepareMemory()
			defer s.responseTypeSupport.ReleaseMemory(respBuffer)
			s.responseTypeSupport.AsCStruct(respBuffer, resp)
			rc := C.rcl_send_response(s.rclService, &reqHeader.request_id, respBuffer)
			if rc != C.RCL_RET_OK {
				return errorsCastC(rc, "failed to send response")
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
	requestTypeSupport   types.MessageTypeSupport
	responseTypeSupport  types.MessageTypeSupport
}

// NewClient creates a new client.
//
// options must not be modified after passing it to this function. If options is
// nil, default options are used.
func (n *Node) NewClient(
	serviceName string,
	typeSupport types.ServiceTypeSupport,
	options *ClientOptions,
) (c *Client, err error) {
	if options == nil {
		options = NewDefaultClientOptions()
	}
	c = &Client{
		requestTypeSupport:  typeSupport.Request(),
		responseTypeSupport: typeSupport.Response(),
		pendingRequests:     make(map[C.long]chan *clientSendResult),
		node:                n,
		rclClient:           (*C.struct_rcl_client_t)(C.malloc(C.sizeof_struct_rcl_client_t)),
		serviceName:         C.CString(serviceName),
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
		return nil, errorsCastC(rc, "failed to create client")
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
		err = multierror.Append(err, errorsCastC(rc, "failed to finalize client"))
	}
	C.free(unsafe.Pointer(c.rclClient))
	c.rclClient = nil

	C.free(unsafe.Pointer(c.serviceName))

	return err.ErrorOrNil()
}

func (c *Client) Send(ctx context.Context, req types.Message) (types.Message, *RmwServiceInfo, error) {
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

func (c *Client) sendRequest(req types.Message) (chan *clientSendResult, error) {
	c.pendingRequestsMutex.Lock()
	defer c.pendingRequestsMutex.Unlock()

	reqBuf := c.requestTypeSupport.PrepareMemory()
	defer c.requestTypeSupport.ReleaseMemory(reqBuf)
	c.requestTypeSupport.AsCStruct(reqBuf, req)

	var sequenceNumber C.long
	rc := C.rcl_send_request(c.rclClient, reqBuf, &sequenceNumber)
	if rc != C.RCL_RET_OK {
		return nil, errorsCastC(rc, "failed to send request")
	}
	resultChan := make(chan *clientSendResult, 1)
	c.pendingRequests[sequenceNumber] = resultChan

	return resultChan, nil
}

func (c *Client) handleResponse() {
	c.pendingRequestsMutex.Lock()
	defer c.pendingRequestsMutex.Unlock()

	var respHeader C.struct_rmw_service_info_t
	respBuf := c.responseTypeSupport.PrepareMemory()
	defer c.responseTypeSupport.ReleaseMemory(respBuf)
	rc := C.rcl_take_response_with_info(c.rclClient, &respHeader, respBuf)
	if rc != C.RCL_RET_OK {
		log.Println(errorsCastC(rc, "failed to take response"))
		return
	}
	defer delete(c.pendingRequests, respHeader.request_id.sequence_number)
	defer close(c.pendingRequests[respHeader.request_id.sequence_number])

	result := &clientSendResult{
		resp: c.responseTypeSupport.New(),
		info: &RmwServiceInfo{
			SourceTimestamp:   time.Unix(0, int64(respHeader.source_timestamp)),
			ReceivedTimestamp: time.Unix(0, int64(respHeader.received_timestamp)),
			RequestID: RmwRequestID{
				WriterGUID:     *(*[16]int8)(unsafe.Pointer(&respHeader.request_id.writer_guid)),
				SequenceNumber: int64(respHeader.request_id.sequence_number),
			},
		},
	}
	c.responseTypeSupport.AsGoStruct(result.resp, respBuf)

	c.pendingRequests[respHeader.request_id.sequence_number] <- result
}

type clientSendResult struct {
	resp types.Message
	info *RmwServiceInfo
}
