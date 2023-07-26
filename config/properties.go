package config

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

// Properties is struct for application configuration.
type Properties struct {
	AppName     string       `json:"app_name"`
	AuthConfigs []AuthConfig `json:"auth"`
}

// AuthConfig hold auth configuration.
type AuthConfig struct {
	Name    string `json:"name"`
	RoleStr string `json:"role"`
	Key     string `json:"key"`
}

// PropertiesLoader load configuration.
type PropertiesLoader interface {
	// Load retrive configuration from configLocation, and then decode to Properties.
	Load(configLocation string) *Properties
}

type propertiesLoader struct {
	logger *zap.Logger
}

// NewPropertiesLoader initialize PropertiesLoader implementation.
func NewPropertiesLoader(logger *zap.Logger) PropertiesLoader {
	return &propertiesLoader{logger: logger}
}

func (pl *propertiesLoader) Load(configLocation string) *Properties {
	f, err := os.ReadFile(configLocation)
	if err != nil {
		pl.logger.Fatal("Failed to read condiguration", zap.Error(err), zap.String("fileName", configLocation))
	}
	var prop Properties
	if err := json.Unmarshal(f, &prop); err != nil {
		pl.logger.Fatal("Failed to decode confugiration", zap.Error(err), zap.String("fileName", configLocation))
	}
	return &prop
}
