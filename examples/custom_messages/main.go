package main

import (
	"github.com/tiiuae/rclgo/pkg/ros2"
	"github.com/tiiuae/rclgo/pkg/ros2/px4_msgs"
)


func main() {
	msg := px4_msgs.SensorCombined{}
	ros2.RclInit()
	ros2.NodeCreate()
	go ros2.SubscriptionCreate()
	go ros2.PublisherCreate()
	ros2...
}