package config

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

/*
Transparently parse all configurations from multiple sources
	(ENV, CLI, conf-file)
when the program starts.
Configurations take a logical precedence order. See viper docs.
Config is globally accessible with
	viper.Get("device-id")
*/
func init() {
	// Define CLI options
	flag.StringP("config-path", "c", "", "Load configuration-file 'communication_link' from the given directory (defaults \"'.', '/etc/communication_link/'\")")
	flag.StringP("device-id", "d", "", "The provisioned device id")
	flag.StringP("mqtt-broker", "m", "ssl://mqtt.googleapis.com:8883", "MQTT broker protocol, address and port")
	flag.StringP("private-key", "k", "/enclave/rsa_private.pem", "The private key for the MQTT authentication")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	// Configure parsing ENV, COML_DEVICE_ID => viper.GetString("device-id")
	viper.SetEnvPrefix("COML")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	// Configure loading configuration files
	viper.SetConfigName("communication_link")
	if viper.GetString("config-path") != "" { // ENV|flag overload
		viper.AddConfigPath(viper.GetString("config-path"))
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/communication_link/")
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	} else {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
		})
	}
}
