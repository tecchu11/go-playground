package main

import (
	"embed"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

type ConfigServer struct {
	Addr         string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	GraceTimeout time.Duration `yaml:"grace_timeout"`
}

type Config struct {
	AppName string       `yaml:"app_name"`
	Svr     ConfigServer `yaml:"server"`
}

//go:embed *.yaml
var conf embed.FS

func LoadConfig(env string) (*Config, error) {
	f := fmt.Sprintf("config-%s.yaml", env)
	buf, err := conf.ReadFile(f)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(buf, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
