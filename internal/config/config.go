package config

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	envDevelopment = "development"
	envProduction  = "production"
)

// Env stores configuration settings extract from enviromental variables
// by using https://github.com/kelseyhightower/envconfig
type Env struct {
	// Env is environment where application is running The value must be
	// "development" or "production".
	Env string `envconfig:"ENV" required:"true"`

	// Port is http serve port.
	Port int `envconfig:"PORT" default:"8000"`
}

// ReadFromEnv reads configuration from environmental variables
// defined by Env struct.
func ReadFromEnv() (*Env, error) {
	var env Env
	if err := envconfig.Process("", &env); err != nil {
		return nil, err
	}
	return &env, nil
}
