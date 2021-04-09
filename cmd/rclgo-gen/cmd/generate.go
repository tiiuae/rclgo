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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tiiuae/rclgo/pkg/gogen"
)

// topicCmd represents the topic command
var generateCmd = &cobra.Command{
	Use:   "generate <root-path>",
	Short: "Generate Golang code from available ROS2 message definitions",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		gogen.Generate(viper.GetString("root-path"), viper.GetString("dest-path"))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			viper.Set("root-path", args[0])
		}
		if viper.GetString("root-path") == "" {
			if os.Getenv("AMENT_PREFIX_PATH") == "" {
				return fmt.Errorf("You haven't sourced your ROS2 environment! Cannot autodetect --root-path. Source your ROS2 or pass --root-path")
			}
			return fmt.Errorf("expecting root-path as the first argument")
		}
		_, err := os.Stat(viper.GetString("root-path"))
		if err != nil {
			return fmt.Errorf("root-path error: %v", err)
		}

		if len(args) > 1 {
			viper.Set("dest-path", args[1])
		}
		if viper.GetString("dest-path") == "" {
			return fmt.Errorf("expecting dest-path as the second argument")
		}
		_, err = os.Stat(viper.GetString("dest-path"))
		if err != nil {
			return fmt.Errorf("dest-path error: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringP("root-path", "r", os.Getenv("AMENT_PREFIX_PATH"), "Root lookup path for ROS2 .msg files. If ROS2 environment is sourced, is autodetected.")
	generateCmd.PersistentFlags().StringP("dest-path", "d", gogen.GetGoConvertedROS2MsgPackagesDir(), "Destination directory for the Golang typed converted ROS2 messages. ROS2 Message structure is preserved as <ros2-package>/msg/<msg-name>")
	viper.BindPFlags(generateCmd.PersistentFlags())
	viper.BindPFlags(generateCmd.LocalFlags())
}
