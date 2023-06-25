package config

import (
	"github.com/kelseyhightower/envconfig"
)

const ServiceName = "WOW"

// Configuration describe services config
type Configuration struct {
	ServerType    string `envconfig:"SERVER_TYPE"    required:"false" default:"tcp"`
	ServerAddress string `envconfig:"SERVER_ADDRESS" required:"false" default:"localhost"`
	ServerPort    int    `envconfig:"SERVER_PORT"    required:"false" default:"2425"`
}

// New initialize configuration
func New() (*Configuration, error) {
	var cfg Configuration
	err := envconfig.Process(ServiceName, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
