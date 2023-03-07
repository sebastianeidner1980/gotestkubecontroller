package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const ConfigFile = "config.yaml"

type Configuration struct {
	Namespace  string `json: "namespace,omitempty"`
	Kubeconfig string `json: "kubeconfig,omitempty"`
	Resource   string `json: "resource,omitempty"`
	LogLevel   string `json: "loglevel,omitempty"`
}

func NewConfiguration(appName string) *Configuration {
	var C Configuration
	cfgFile := strings.Split(ConfigFile, ".")
	viper.SetConfigName(cfgFile[0])          // name of config file (without extension)
	viper.SetConfigType(cfgFile[1])          // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/." + appName) // call multiple times to add many search paths
	viper.AddConfigPath("./" + cfgFile[0])   // optionally look for config in the working directory
	err := viper.ReadInConfig()              // Find and read the config file
	if err != nil {                          // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		panic(fmt.Errorf("fatal error could completly parse config file: %w", err))
	}

	return &C, nil
}
