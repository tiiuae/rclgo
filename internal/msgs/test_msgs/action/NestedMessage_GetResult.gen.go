/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

// Code generated by rclgo-gen. DO NOT EDIT.

package test_msgs_action

/*
#include <rosidl_runtime_c/message_type_support_struct.h>
#include <test_msgs/action/nested_message.h>
*/
import "C"

import (
	"context"
	"errors"
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
)

func init() {
	typemap.RegisterService("test_msgs/NestedMessage_GetResult", NestedMessage_GetResultTypeSupport)
}

type _NestedMessage_GetResultTypeSupport struct {}

func (s _NestedMessage_GetResultTypeSupport) Request() types.MessageTypeSupport {
	return NestedMessage_GetResult_RequestTypeSupport
}

func (s _NestedMessage_GetResultTypeSupport) Response() types.MessageTypeSupport {
	return NestedMessage_GetResult_ResponseTypeSupport
}

func (s _NestedMessage_GetResultTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_service_type_support_handle__test_msgs__action__NestedMessage_GetResult())
}

// Modifying this variable is undefined behavior.
var NestedMessage_GetResultTypeSupport types.ServiceTypeSupport = _NestedMessage_GetResultTypeSupport{}

// NestedMessage_GetResultClient wraps rclgo.Client to provide type safe helper
// functions
type NestedMessage_GetResultClient struct {
	*rclgo.Client
}

// NewNestedMessage_GetResultClient creates and returns a new client for the
// NestedMessage_GetResult
func NewNestedMessage_GetResultClient(node *rclgo.Node, serviceName string, options *rclgo.ClientOptions) (*NestedMessage_GetResultClient, error) {
	client, err := node.NewClient(serviceName, NestedMessage_GetResultTypeSupport, options)
	if err != nil {
		return nil, err
	}
	return &NestedMessage_GetResultClient{client}, nil
}

func (s *NestedMessage_GetResultClient) Send(ctx context.Context, req *NestedMessage_GetResult_Request) (*NestedMessage_GetResult_Response, *rclgo.RmwServiceInfo, error) {
	msg, rmw, err := s.Client.Send(ctx, req)
	if err != nil {
		return nil, rmw, err
	}
	typedMessage, ok := msg.(*NestedMessage_GetResult_Response)
	if !ok {
		return nil, rmw, errors.New("invalid message type returned")
	}
	return typedMessage, rmw, err
}

type NestedMessage_GetResultServiceResponseSender struct {
	sender rclgo.ServiceResponseSender
}

func (s NestedMessage_GetResultServiceResponseSender) SendResponse(resp *NestedMessage_GetResult_Response) error {
	return s.sender.SendResponse(resp)
}

type NestedMessage_GetResultServiceRequestHandler func(*rclgo.RmwServiceInfo, *NestedMessage_GetResult_Request, NestedMessage_GetResultServiceResponseSender)

// NestedMessage_GetResultService wraps rclgo.Service to provide type safe helper
// functions
type NestedMessage_GetResultService struct {
	*rclgo.Service
}

// NewNestedMessage_GetResultService creates and returns a new service for the
// NestedMessage_GetResult
func NewNestedMessage_GetResultService(node *rclgo.Node, name string, options *rclgo.ServiceOptions, handler NestedMessage_GetResultServiceRequestHandler) (*NestedMessage_GetResultService, error) {
	h := func(rmw *rclgo.RmwServiceInfo, msg types.Message, rs rclgo.ServiceResponseSender) {
		m := msg.(*NestedMessage_GetResult_Request)
		responseSender := NestedMessage_GetResultServiceResponseSender{sender: rs} 
		handler(rmw, m, responseSender)
	}
	service, err := node.NewService(name, NestedMessage_GetResultTypeSupport, options, h)
	if err != nil {
		return nil, err
	}
	return &NestedMessage_GetResultService{service}, nil
}