package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	ConfigFile string = "config.yaml"
)

type Configuration struct {
	Namespace  string `json: "namespace,omitempty"`
	Kubeconfig string `json: "kubeconfig,omitempty"`
	Resource   string `json: "resource,omitempty"`
}

func NewConfiguration() *Configuration {
	var C Configuration
	viper.SetConfigName("config")         // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.Unmarshal(&C)

	return &C
}
