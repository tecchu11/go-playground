package configs

import (
	"embed"
	"errors"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// ApplicationProperties is struct for application configuration.
type ApplicationProperties struct {
	AppName      string       `yaml:"app_name"`
	ServerConfig ServerConfig `yaml:"server"`
	AuthConfigs  []AuthConfig `yaml:"auth"`
}

// ServerConfig properties for http.Server.
type ServerConfig struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	GraceTimeout time.Duration `yaml:"grace_timeout"`
}

// AuthConfig hold auth configuration.
type AuthConfig struct {
	Name    string `yaml:"name"`
	RoleStr string `yaml:"role"`
	Key     string `yaml:"key"`
}

var (
	ErrConfigNotFound  = errors.New("configuration not found by env")
	ErrConfigUnmarshal = errors.New("failed to unmarshal to ApplicationProperties")
)

//go:embed *.yaml
var configs embed.FS

// Load loads config of env.
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
