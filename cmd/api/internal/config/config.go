package config

import (
	"go-playground/pkg/env"
)

type Configuration struct {
	Env        string `env:"APP_ENV"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBAddr     string `env:"DB_HOST"`
	DBName     string `env:"DB_NAME"`
}

// Load Configuration from environment variables.
func Load() *Configuration {
	var conf Configuration
	env.Decode(&conf)
	return &conf
}
