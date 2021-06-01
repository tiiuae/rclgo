/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package ros2_type_dispatcher

import (
	"fmt"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
	"gopkg.in/yaml.v3"
)

/*
ROS2MsgTypeNameToGoROS2Msg maps the ROS2 Message type name to the type implementation in Go,
so the correct type can be dynamically chosen.
Needs to be defined for all supported ROS2 message types.
Could also be implemented using Golang-plugins.
Or maybe some pretty HC Golang hacks to dynamically cast a string to a Golang type.
*/
var ROS2MsgTypeNameToGoROS2Msg = make(map[string]ros2types.ROS2Msg)

/*
RegisterROS2MsgTypeNameAlias sets the type string to implementation dispatcher, so the correct type can be dynamically chosen.
The Golang types of ROS2 Message use

    func init() {}

to automatically populate this when imported.
*/
func RegisterROS2MsgTypeNameAlias(alias string, msgType ros2types.ROS2Msg) {
	ROS2MsgTypeNameToGoROS2Msg[alias] = msgType
}

/*
TranslateROS2MsgTypeNameToType translates for ex
    "std_msgs/ColorRGBA"
to
    std_msgs.ColorRGBA -Go type

returns true if the type mapping is found
*/
func TranslateROS2MsgTypeNameToType(msgType string) (ros2types.ROS2Msg, bool) {
	ros2msg, ok := ROS2MsgTypeNameToGoROS2Msg[msgType]
	return ros2msg, ok
}

/*
TranslateROS2MsgTypeNameToTypeMust panics if there is no mapping
*/
func TranslateROS2MsgTypeNameToTypeMust(msgType string) ros2types.ROS2Msg {
	ros2msg, ok := ROS2MsgTypeNameToGoROS2Msg[msgType]
	if !ok {
		panic(fmt.Sprintf("No registered implementation for ROS2 message type '%s'!\n", msgType))
	}
	return ros2msg
}

/*
TranslateMsgPayloadYAMLToROS2Msg returns a new instance of the given ROS2Msg-object

The ROS2 official cli-client uses YAML to define the data payload, so do we.
*/
func TranslateMsgPayloadYAMLToROS2Msg(yamlString string, ros2msg ros2types.ROS2Msg) (ros2types.ROS2Msg, error) {
	yamlBytes := []byte(yamlString)
	ros2msgClone := ros2msg.Clone()
	ros2msgClone.SetDefaults(nil)
	err := yaml.Unmarshal(yamlBytes, ros2msgClone)
	return ros2msgClone, err
}

// serviceTypeToGoServiceDefinition is the ROS2MsgTypeNameToGoROS2Msg equivalent
// for services.
var serviceTypeToGoServiceDefinition = make(map[string]ros2types.Service)

// RegisterROS2ServiceTypeNameAlias is the RegisterROS2MsgTypeNameAlias
// equivalent for services.
func RegisterROS2ServiceTypeNameAlias(alias string, srvType ros2types.Service) {
	serviceTypeToGoServiceDefinition[alias] = srvType
}

// TranslateROS2ServiceTypeNameToType is the TranslateROS2MsgTypeNameToType
// equivalent for services.
func TranslateROS2ServiceTypeNameToType(srvType string) (ros2types.Service, bool) {
	srv, ok := serviceTypeToGoServiceDefinition[srvType]
	return srv, ok
}

// TranslateROS2ServiceTypeNameToTypeMust is the
// TranslateROS2MsgTypeNameToTypeMust equivalent for services.
func TranslateROS2ServiceTypeNameToTypeMust(srvType string) ros2types.Service {
	srv, ok := serviceTypeToGoServiceDefinition[srvType]
	if !ok {
		panic(fmt.Sprintf("No registered implementation for ROS2 message type '%s'!\n", srvType))
	}
	return srv
}
