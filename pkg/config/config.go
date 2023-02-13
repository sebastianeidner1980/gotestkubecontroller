package config

import(
	"github.com/spf13/viper"
)

var (
	ConfigFile string = "config.yaml"
)

type Configuration struct {
	Namespace string `json: "namespace,omitempty"`
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}


