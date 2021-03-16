package ros2_type_dispatcher

import (
	"github.com/tiiuae/rclgo/pkg/ros2"
	"gopkg.in/yaml.v3"
)

/*
	Seed the type string to implementation dispatcher, so the correct type can be dynamically chosen.
	Needs to be defined for all supported ROS2 message types.
	Could also be implemented using Golang-plugins.
	Or maybe some pretty HC Golang hacks to dynamically cast a string to a Golang type.
*/
var ROS2MsgTypeNameToGoROS2Msg = make(map[string]ros2.ROS2Msg)

// Seed the type string to implementation dispatcher, so the correct type can be dynamically chosen.
func RegisterROS2MsgTypeNameAlias(alias string, msgType ros2.ROS2Msg) {
	ROS2MsgTypeNameToGoROS2Msg[alias] = msgType
}

func TranslateROS2MsgTypeNameToType(msgType string) ros2.ROS2Msg {
	return ROS2MsgTypeNameToGoROS2Msg[msgType]
}

/*
Returns a new instance of the given ROS2Msg-object

The ROS2 official cli-client uses YAML to define the data payload, so do we.
*/
func TranslateMsgPayloadYAMLToROS2Msg(yamlString string, ros2msg ros2.ROS2Msg) (ros2.ROS2Msg, error) {
	yamlBytes := []byte(yamlString)
	ros2msgClone := ros2msg.Clone()
	err := yaml.Unmarshal(yamlBytes, ros2msgClone)
	return ros2msgClone, err
}
