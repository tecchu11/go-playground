package config

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
)

// ApplicationProperties is struct for application configuration.
type ApplicationProperties struct {
	AppName     string       `json:"app_name"`
	AuthConfigs []AuthConfig `json:"auth"`
}

// AuthConfig hold auth configuration.
type AuthConfig struct {
	Name    string `json:"name"`
	RoleStr string `json:"role"`
	Key     string `json:"key"`
}

var (
	ErrConfigNotFound  = errors.New("cofiguration not found by env")
	ErrConfigUnmarshal = errors.New("failed to unmarshal to ApplicationProperties")
)

//go:embed *.json
var configs embed.FS

func Load(env string) (*ApplicationProperties, error) {
	key := fmt.Sprintf("config-%s.json", env)
	f, err := configs.ReadFile(key)
	if err != nil {
		return nil, errors.Join(ErrConfigNotFound, err)
	}
	var prop ApplicationProperties
	if err := json.Unmarshal(f, &prop); err != nil {
		return nil, errors.Join(ErrConfigUnmarshal, err)
	}
	return &prop, nil
}
