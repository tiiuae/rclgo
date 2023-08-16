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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rclgo-gen",
	Short: "ROS2 client library in Golang - ROS2 Message generator",
	Long:  `Call this program to generate Go types for the ROS2 messages found in your ROS2 environment.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("config-file", "c", "", "config file (default is $HOME/.rclgo.yaml)")
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
