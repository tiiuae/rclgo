/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package ros2

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	std_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/std_msgs/msg"
)

func TestContextClose(t *testing.T) {
	var context *Context
	defer func() {
		if context != nil {
			context.Close()
		}
	}()
	SetDefaultFailureMode(FailureContinues)
	Convey("Scenario: does Context handle closing resources correctly", t, func() {
		Convey("Given a context with resources", func() {
			var err error
			context, err = newDefaultRCLContext()
			So(err, ShouldBeNil)
			_, err = context.NewNode("node1", "/test/context_close")
			So(err, ShouldBeNil)
			node2, err := context.NewNode("node2", "/test/context_close")
			So(err, ShouldBeNil)
			_, err = context.NewWaitSet(time.Second)
			So(err, ShouldBeNil)
			_, err = context.NewNode("node3", "/test/context_close")
			So(err, ShouldBeNil)
			_, err = node2.NewPublisher("/test_topic", &std_msgs.String{})
			So(err, ShouldBeNil)
			_, err = node2.NewSubscription("/test_topic", &std_msgs.ColorRGBA{}, func(s *Subscription) {})
			So(err, ShouldBeNil)
			_, err = node2.NewPublisher("/test_topic2", &std_msgs.ColorRGBA{})
			So(err, ShouldBeNil)
		})
		Convey("When the context is closed the first time, no errors occur", func() {
			So(context.Close(), ShouldBeNil)
		})
		Convey("When the context is closed the second time, no errors occur", func() {
			So(context.Close(), ShouldBeNil)
		})
	})
}
