package config

import (
	"embed"
	"encoding/json"
	"fmt"
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

//go:embed *.json
var configs embed.FS

func LoadConfigWith(key string) (*Properties, error) {
	f, err := configs.ReadFile(key)
	if err != nil {
		return nil, fmt.Errorf("failed read configuration properties by %s because %s", key, err.Error())
	}
	var prop Properties
	if err := json.Unmarshal(f, &prop); err != nil {
		return nil, fmt.Errorf("failed to unmarashal to Properties from config %s because %s", key, err.Error())
	}
	return &prop, nil
}
