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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tiiuae/rclgo/pkg/ros2"
	"github.com/tiiuae/rclgo/pkg/ros2/std_msgs"
)

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:   "echo",
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

		subscription, err := ros2.SubscriptionCreate(rclContext, rcl_node, viper.GetString("topic-name"), &std_msgs.ColorRGBA{}, func(sub ros2.Subscription, msg ros2.ROS2Msg) {
			fmt.Printf("Cobra received message!\n")
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

		fmt.Printf("%v%v\n", rclContext, rcl_node)

	},
}

func init() {
	topicCmd.AddCommand(echoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// echoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// echoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
