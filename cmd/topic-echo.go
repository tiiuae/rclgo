/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
)

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:   "echo <topic-name> <msg-type>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		rclContext, err := ros2.RclInit()
		if err != nil {
			fmt.Printf("Error '%+v' ros2.RclInit.\n", err)
			panic(err)
		}

		rcl_node, err := ros2.NodeCreate(rclContext, viper.GetString("node-name"), viper.GetString("namespace"))
		if err != nil {
			fmt.Printf("Error '%+v' ros2.NodeCreate.\n", err)
			panic(err)
		}

		ros2msg := ros2_type_dispatcher.TranslateROS2MsgTypeNameToType(viper.GetString("msg-type"))
		ros2msgClone := ros2msg.Clone()
		subscription, err := ros2.SubscriptionCreate(rclContext, rcl_node, viper.GetString("topic-name"), ros2msgClone,
			func(subscription *ros2.Subscription, ros2_msg_receive_buffer unsafe.Pointer, rmwMessageInfo *ros2.RmwMessageInfo) {
				ros2msgClone.AsGoStruct(ros2_msg_receive_buffer)
				fmt.Printf("%+v ", ros2msgClone)
				fmt.Printf("Source_timestamp='%s' Received_timestamp='%s'\n", rmwMessageInfo.Source_timestamp.Format(time.RFC3339Nano), rmwMessageInfo.Received_timestamp.Format(time.RFC3339Nano))
			})
		if err != nil {
			fmt.Printf("Error '%+v' SubscriptionCreate.\n", err)
			panic(err)
		}

		subscriptions := []ros2.Subscription{subscription}
		waitSet, err := ros2.WaitSetCreate(rclContext, subscriptions, nil, 1000*time.Millisecond)
		if err != nil {
			fmt.Printf("Error '%+v' WaitSetCreate.\n", err)
			panic(err)
		}

		err = ros2.WaitSetRun(waitSet)
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

		return nil
	},
}

func init() {
	topicCmd.AddCommand(echoCmd)
	// Defined flags
	viper.BindPFlags(echoCmd.PersistentFlags())
	viper.BindPFlags(echoCmd.LocalFlags())
}
