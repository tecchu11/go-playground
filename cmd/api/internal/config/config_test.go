package config_test

import (
	"go-playground/cmd/api/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tests := map[string]struct {
		env          string
		expectedConf *config.Config
	}{
		"env is local": {
			env: "local",
			expectedConf: &config.Config{
				AppName: "go-playground",
				Svr: config.ConfigServer{
					Addr:         ":8080",
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  120 * time.Second,
					GraceTimeout: 20 * time.Second,
				},
			},
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			actualConf, actualErr := config.Load(v.env)
			require.NoError(t, actualErr)
			require.Equal(t, v.expectedConf, actualConf)
		})
	}
}
