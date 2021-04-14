/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tiiuae/rclgo/pkg/ros2"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"

	_ "github.com/tiiuae/rclgo/pkg/ros2/msgs" // Load all the available ROS2 Message types. In Go one cannot dynamically import.
)

// pubCmd represents the pub command
var pubCmd = &cobra.Command{
	Use:   "pub <topic-name> <msg-type> <payload>",
	Short: "Publish a message to a topic",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v\n", viper.AllSettings())

		rclContext, err := ros2.NewRCLContext(nil, 0, ros2.NewRCLArgsMust(viper.GetString("ros-args")))
		if err != nil {
			fmt.Printf("Error '%+v' ros2.NewRCLContext.\n", err)
			panic(err)
		}

		rclNode, err := ros2.NewNode(rclContext, viper.GetString("node-name"), viper.GetString("namespace"))
		if err != nil {
			fmt.Printf("Error '%+v' ros2.NewNode.\n", err)
			panic(err)
		}

		ros2msg := ros2_type_dispatcher.TranslateROS2MsgTypeNameToTypeMust(viper.GetString("msg-type"))
		ros2msg, err_yaml := ros2_type_dispatcher.TranslateMsgPayloadYAMLToROS2Msg(strings.ReplaceAll(viper.GetString("payload"), "\\n", "\n"), ros2msg)
		if err_yaml != nil {
			panic(fmt.Sprintf("Error '%v' unmarshalling YAML '%s' to ROS2 message type '%s'", err_yaml, viper.GetString("payload"), viper.GetString("msg-type")))
		}

		publisher, err := rclNode.NewPublisher(viper.GetString("topic-name"), ros2msg)
		if err != nil {
			fmt.Printf("Error '%+v' rclNode.NewPublisher.\n", err)
			panic(err)
		}

		timer, err := ros2.NewTimer(rclContext, 0, func(timer *ros2.Timer) {
			fmt.Printf("%+v\n", ros2msg)
			publisher.Publish(ros2msg)
		})
		if err != nil {
			fmt.Printf("Error '%+v' TimerCreate.\n", err)
			panic(err)
		}

		timers := []*ros2.Timer{timer}
		waitSet, err := ros2.NewWaitSet(rclContext, nil, timers, 1000*time.Millisecond)
		if err != nil {
			fmt.Printf("Error '%+v' WaitSetCreate.\n", err)
			panic(err)
		}

		err = waitSet.Run()
		if err != nil {
			fmt.Printf("Error '%+v' WaitSetRun.\n", err)
			panic(err)
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
		if !viper.IsSet("payload") {
			if len(args) > 2 {
				viper.Set("payload", args[2])
			} else {
				return fmt.Errorf("expecting ROS2 message payload as the third argument")
			}
		}

		return nil
	},
}

func init() {
	topicCmd.AddCommand(pubCmd)
	topicCmd.PersistentFlags().StringP("payload", "p", "", "Values to fill the message with in YAML format (e.g. 'data: Hello World'), otherwise the message will be published with default values")
	viper.BindPFlags(pubCmd.PersistentFlags())
	viper.BindPFlags(pubCmd.LocalFlags())
}
