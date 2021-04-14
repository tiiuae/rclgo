/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tiiuae/rclgo/pkg/datagenerator"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rclgo",
	Short: "ROS2 client library in Golang",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rand.Seed(time.Now().UnixNano())
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("config-file", "c", "", "config file (default is $HOME/.rclgo.yaml)")
	rootCmd.PersistentFlags().StringP("node-name", "n", datagenerator.NodeName(), "Node name")
	rootCmd.PersistentFlags().StringP("namespace", "s", "/", "Namespace name")
	rootCmd.PersistentFlags().StringP("ros-args", "", "", "See. http://design.ros2.org/articles/ros_command_line_arguments.html . Anything which is between quotation marks (\"\") is forwarded to the ROS2 rcl init. Example --ros-args \"--log-level DEBUG --enclave=\"/foo/bar\"\"")
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindPFlags(rootCmd.LocalFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if viper.GetString("config-file") != "" {
		// Use config file from the flag.
		viper.SetConfigFile(viper.GetString("config-file"))
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".rclgo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rclgo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
