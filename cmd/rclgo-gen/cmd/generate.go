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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tiiuae/rclgo/pkg/gogen"
)

func validateGenerateArgs(cmd *cobra.Command, args []string) error {
	if viper.GetString(getKey(cmd, "root-path")) == "" {
		if os.Getenv("AMENT_PREFIX_PATH") == "" {
			return fmt.Errorf("You haven't sourced your ROS2 environment! Cannot autodetect --root-path. Source your ROS2 or pass --root-path")
		}
		return fmt.Errorf("root-path is required")
	}
	_, err := os.Stat(viper.GetString(getKey(cmd, "root-path")))
	if err != nil {
		return fmt.Errorf("root-path error: %v", err)
	}
	if viper.GetString(getKey(cmd, "dest-path")) == "" {
		return fmt.Errorf("dest-path is required")
	}
	_, err = os.Stat(viper.GetString(getKey(cmd, "dest-path")))
	if err != nil {
		return fmt.Errorf("dest-path error: %v", err)
	}
	return nil
}

// topicCmd represents the topic command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go bindings for ROS2 interface definitions under <root-path>",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootPath := viper.GetString(getKey(cmd, "root-path"))
		destPath := viper.GetString(getKey(cmd, "dest-path"))
		if err := gogen.GenerateGolangMessageTypes(rootPath, destPath); err != nil {
			return fmt.Errorf("failed to generate interface bindings: %w", err)
		}
		if err := gogen.GenerateROS2AllMessagesImporter(destPath); err != nil {
			return fmt.Errorf("failed to generate all importer: %w", err)
		}
		return nil
	},
	Args: validateGenerateArgs,
}

var generateRclgoCmd = &cobra.Command{
	Use:   "generate-rclgo",
	Short: "Generate Go code that forms a part of rclgo",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootPath := viper.GetString(getKey(cmd, "root-path"))
		destPath := viper.GetString(getKey(cmd, "dest-path"))
		if err := gogen.GeneratePrimitives(destPath); err != nil {
			return fmt.Errorf("failed to generate primitive types: %w", err)
		}
		if err := gogen.GenerateROS2ErrorTypes(rootPath, destPath); err != nil {
			return fmt.Errorf("failed to generate error types: %w", err)
		}
		return nil
	},
	Args: validateGenerateArgs,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringP("root-path", "r", os.Getenv("AMENT_PREFIX_PATH"), "Root lookup path for ROS2 .msg files. If ROS2 environment is sourced, is autodetected.")
	generateCmd.PersistentFlags().StringP("dest-path", "d", ".", "Destination directory for the Golang typed converted ROS2 messages. ROS2 Message structure is preserved as <ros2-package>/msg/<msg-name>")
	bindPFlags(generateCmd)

	rootCmd.AddCommand(generateRclgoCmd)
	generateRclgoCmd.PersistentFlags().StringP("root-path", "r", os.Getenv("AMENT_PREFIX_PATH"), "Root lookup path for ROS2 .msg files. If ROS2 environment is sourced, is autodetected.")
	generateRclgoCmd.PersistentFlags().StringP("dest-path", "d", gogen.RclgoRepoRootPath(), "Path to the root of the rclgo repository")
	bindPFlags(generateRclgoCmd)
}

func getPrefix(cmd *cobra.Command) string {
	parts := []string{}
	for c := cmd; c != c.Root(); c = c.Parent() {
		parts = append(parts, c.Name())
	}
	for i := 0; i < len(parts)/2; i++ {
		parts[i], parts[len(parts)-i-1] = parts[len(parts)-i-1], parts[i]
	}
	prefix := strings.Join(parts, ".")
	if prefix != "" {
		prefix += "."
	}
	return prefix
}

func getKey(cmd *cobra.Command, key string) string {
	return getPrefix(cmd) + key
}

func bindPFlags(cmd *cobra.Command) {
	prefix := getPrefix(cmd)
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(prefix+f.Name, f)
	})
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(prefix+f.Name, f)
	})
}
