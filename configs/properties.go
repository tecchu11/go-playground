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

var (
	//go:embed *.yaml
	configs embed.FS
	envs    = map[string]string{"local": "config-local.yaml"}
)

// Load loads config of env.
func Load(env string) (*ApplicationProperties, error) {
	key, ok := envs[env]
	if !ok {
		return nil, fmt.Errorf("not found config by given env %s", env)
	}
	f, err := configs.ReadFile(key)
	if err != nil {
		return nil, fmt.Errorf("not found config by given env %s because %w", env, err)
	}
	var prop ApplicationProperties
	if err := yaml.Unmarshal(f, &prop); err != nil {
		return nil, fmt.Errorf("failed to unmarshal of %s because %w", f, err)
	}
	return &prop, nil
}
