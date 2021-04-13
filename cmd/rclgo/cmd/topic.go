/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// topicCmd represents the topic command
var topicCmd = &cobra.Command{
	Use:   "topic",
	Short: "Topic operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("topic called")
	},
}

func init() {
	rootCmd.AddCommand(topicCmd)
	topicCmd.PersistentFlags().StringP("topic-name", "t", "", "Name of the ROS topic to publish to (e.g. '/chatter')")
	topicCmd.PersistentFlags().StringP("msg-type", "m", "", "Type of the ROS message (e.g. 'std_msgs/String')")
	viper.BindPFlags(topicCmd.PersistentFlags())
	viper.BindPFlags(topicCmd.LocalFlags())
}
