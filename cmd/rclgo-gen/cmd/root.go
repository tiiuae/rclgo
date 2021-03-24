/*
Copyright © 2021 Solita Oy <olli-antti.kivilahti@solita.fi>

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
	"math/rand"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rclgo-gen",
	Short: "ROS2 client library in Golang - ROS2 Message generator",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
