/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tiiuae/rclgo/pkg/ros2"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"

	_ "github.com/tiiuae/rclgo/pkg/ros2/msgs" // Load all the available ROS2 Message types. In Go one cannot dynamically import.
)

// pubCmd represents the pub command
var pubCmd = &cobra.Command{
	Use:   "pub <topic-name> <msg-type> <payload>",
	Short: "Publish a message to a topic",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%+v\n", viper.AllSettings())
		// attach sigint & sigterm listeners
		terminationSignals := make(chan os.Signal, 1)
		signal.Notify(terminationSignals, syscall.SIGINT, syscall.SIGTERM)

		rclContext, errs := ros2.PublisherBundleTimer(nil, nil, viper.GetString("namespace"), viper.GetString("node-name"), viper.GetString("topic-name"), viper.GetString("msg-type"), ros2.NewRCLArgsMust(viper.GetString("ros-args")), 1000*time.Millisecond, viper.GetString("payload"),
			func(p *ros2.Publisher, m ros2types.ROS2Msg) bool {
				fmt.Printf("%+v\n", m)
				return true
			})
		if errs != nil {
			fmt.Println(errs)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}

		<-terminationSignals
		fmt.Printf("Closing topic pub\n")
		errs = ros2.RCLContextFini(rclContext)
		if errs != nil {
			fmt.Println(errs)
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
