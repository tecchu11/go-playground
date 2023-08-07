package configs

import (
	"embed"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

// ApplicationProperties is struct for application configuration.
type ApplicationProperties struct {
	AppName     string       `yaml:"app_name"`
	AuthConfigs []AuthConfig `yaml:"auth"`
}

// AuthConfig hold auth configuration.
type AuthConfig struct {
	Name    string `yaml:"name"`
	RoleStr string `yaml:"role"`
	Key     string `yaml:"key"`
}

var (
	ErrConfigNotFound  = errors.New("cofiguration not found by env")
	ErrConfigUnmarshal = errors.New("failed to unmarshal to ApplicationProperties")
)

//go:embed *.yaml
var configs embed.FS

func Load(env string) (*ApplicationProperties, error) {
	key := fmt.Sprintf("config-%s.yaml", env)
	f, err := configs.ReadFile(key)
	if err != nil {
		return nil, errors.Join(ErrConfigNotFound, err)
	}
	var prop ApplicationProperties
	if err := yaml.Unmarshal(f, &prop); err != nil {
		return nil, errors.Join(ErrConfigUnmarshal, err)
	}
	return &prop, nil
}
