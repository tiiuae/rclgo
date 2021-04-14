/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package cmd

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tiiuae/rclgo/pkg/ros2"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"

	_ "github.com/tiiuae/rclgo/pkg/ros2/msgs" // Load all the available ROS2 Message types. In Go one cannot dynamically import.
)

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:   "echo <topic-name> <msg-type>",
	Short: "Output messages from a topic",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rclContext, err := ros2.NewRCLContext(nil, 0, ros2.NewRCLArgsMust(viper.GetString("ros-args")))
		if err != nil {
			panic(fmt.Sprintf("Error '%+v' ros2.NewRCLContext.\n", err))
		}

		rclNode, err := ros2.NewNode(rclContext, viper.GetString("node-name"), viper.GetString("namespace"))
		if err != nil {
			panic(fmt.Sprintf("Error '%+v' ros2.NewNode.\n", err))
		}

		ros2msg := ros2_type_dispatcher.TranslateROS2MsgTypeNameToTypeMust(viper.GetString("msg-type"))
		ros2msgClone := ros2msg.Clone()
		subscription, err := rclNode.NewSubscription(viper.GetString("topic-name"), ros2msgClone,
			func(subscription *ros2.Subscription, ros2_msg_receive_buffer unsafe.Pointer, rmwMessageInfo *ros2.RmwMessageInfo) {
				ros2msgClone.AsGoStruct(ros2_msg_receive_buffer)
				fmt.Printf("%+v ", ros2msgClone)
				fmt.Printf("SourceTimestamp='%s' ReceivedTimestamp='%s'\n", rmwMessageInfo.SourceTimestamp.Format(time.RFC3339Nano), rmwMessageInfo.ReceivedTimestamp.Format(time.RFC3339Nano))
			})
		if err != nil {
			panic(fmt.Sprintf("Error '%+v' SubscriptionCreate.\n", err))
		}

		subscriptions := []*ros2.Subscription{subscription}
		waitSet, err := ros2.NewWaitSet(rclContext, subscriptions, nil, 1000*time.Millisecond)
		if err != nil {
			panic(fmt.Sprintf("Error '%+v' WaitSetCreate.\n", err))
		}

		err = waitSet.Run()
		if err != nil {
			panic(fmt.Sprintf("Error '%+v' WaitSetRun.\n", err))
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet("topic-name") {
			if len(args) > 0 {
				viper.Set("topic-name", args[0])
			} else {
				return fmt.Errorf("expecting rcl topic name as the first argument")
			}
		}
		if !viper.IsSet("msg-type") {
			if len(args) > 1 {
				viper.Set("msg-type", args[1])
			} else {
				return fmt.Errorf("expecting ROS2 message type as the second argument")
			}
		}

		return nil
	},
}

func init() {
	topicCmd.AddCommand(echoCmd)
	// Defined flags
	viper.BindPFlags(echoCmd.PersistentFlags())
	viper.BindPFlags(echoCmd.LocalFlags())
}
