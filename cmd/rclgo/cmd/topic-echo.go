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
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tiiuae/rclgo/pkg/ros2"

	_ "github.com/tiiuae/rclgo/pkg/ros2/msgs" // Load all the available ROS2 Message types. In Go one cannot dynamically import.
)

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:   "echo <topic-name> <msg-type>",
	Short: "Output messages from a topic",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%+v\n", viper.AllSettings())
		// attach sigint & sigterm listeners
		terminationSignals := make(chan os.Signal, 1)
		signal.Notify(terminationSignals, syscall.SIGINT, syscall.SIGTERM)

		rclContext, errs := ros2.SubscriberBundle(nil, nil, viper.GetString("namespace"), viper.GetString("node-name"), viper.GetString("topic-name"), viper.GetString("msg-type"), ros2.NewRCLArgsMust(viper.GetString("ros-args")),
			func(s *ros2.Subscription, p unsafe.Pointer, rmi *ros2.RmwMessageInfo) {
				msgClone := s.Ros2MsgType.Clone()
				msgClone.AsGoStruct(p)
				fmt.Printf("%+v ", msgClone)
				fmt.Printf("SourceTimestamp='%s' ReceivedTimestamp='%s'\n", rmi.SourceTimestamp.Format(time.RFC3339Nano), rmi.ReceivedTimestamp.Format(time.RFC3339Nano))
			})
		if errs != nil {
			fmt.Println(errs)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}

		<-terminationSignals
		fmt.Printf("Closing topic echo\n")
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

		return nil
	},
}

func init() {
	topicCmd.AddCommand(echoCmd)
	// Defined flags
	viper.BindPFlags(echoCmd.PersistentFlags())
	viper.BindPFlags(echoCmd.LocalFlags())
}
