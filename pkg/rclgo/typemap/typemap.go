/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package typemap

import (
	"fmt"

	"github.com/tiiuae/rclgo/pkg/rclgo/types"
)

/*
messageTypeMap maps the ROS2 Message type name to the type implementation in Go,
so the correct type can be dynamically chosen.
Needs to be defined for all supported ROS2 message types.
Could also be implemented using Golang-plugins.
Or maybe some pretty HC Golang hacks to dynamically cast a string to a Golang type.
*/
var messageTypeMap = make(map[string]types.MessageTypeSupport)

/*
RegisterMessage sets the type string to implementation dispatcher, so the
correct type can be dynamically chosen. The Golang types of ROS2 Message use

	func init() {}

to automatically populate this when imported.
*/
func RegisterMessage(alias string, msgType types.MessageTypeSupport) {
	messageTypeMap[alias] = msgType
}

/*
GetMessage translates for ex

	"std_msgs/ColorRGBA"

to

	std_msgs.ColorRGBA -Go type

returns true if the type mapping is found
*/
func GetMessage(msgType string) (types.MessageTypeSupport, bool) {
	ros2msg, ok := messageTypeMap[msgType]
	return ros2msg, ok
}

/*
GetMessageMust panics if there is no mapping
*/
func GetMessageMust(msgType string) types.MessageTypeSupport {
	ros2msg, ok := messageTypeMap[msgType]
	if !ok {
		panic(fmt.Sprintf("No registered implementation for ROS2 message type '%s'!\n", msgType))
	}
	return ros2msg
}

// serviceTypeMap is the messageTypeMap equivalent for services.
var serviceTypeMap = make(map[string]types.ServiceTypeSupport)

// RegisterService is the RegisterMessage equivalent for services.
func RegisterService(alias string, srvType types.ServiceTypeSupport) {
	serviceTypeMap[alias] = srvType
}

// GetService is the GetMessage equivalent for services.
func GetService(srvType string) (types.ServiceTypeSupport, bool) {
	srv, ok := serviceTypeMap[srvType]
	return srv, ok
}

// GetServiceMust is the GetMessageMust equivalent for services.
func GetServiceMust(srvType string) types.ServiceTypeSupport {
	srv, ok := serviceTypeMap[srvType]
	if !ok {
		panic(fmt.Sprintf("No registered implementation for ROS2 message type '%s'!\n", srvType))
	}
	return srv
}

// actionTypeMap is the messageTypeMap equivalent for actions.
var actionTypeMap = make(map[string]types.ActionTypeSupport)

// RegisterAction is the RegisterMessage equivalent for actions.
func RegisterAction(alias string, actionType types.ActionTypeSupport) {
	actionTypeMap[alias] = actionType
}

// GetAction is the GetMessage equivalent for actions.
func GetAction(actionType string) (types.ActionTypeSupport, bool) {
	action, ok := actionTypeMap[actionType]
	return action, ok
}

// GetActionMust is the GetMessageMust equivalent for actions.
func GetActionMust(actionType string) types.ActionTypeSupport {
	action, ok := actionTypeMap[actionType]
	if !ok {
		panic(fmt.Sprintf("No registered implementation for ROS2 message type '%s'!\n", actionType))
	}
	return action
}
