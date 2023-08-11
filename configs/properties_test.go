package configs_test

import (
	"go-playground/configs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := map[string]struct {
		env       string
		expectErr bool
		expected  *configs.ApplicationProperties
	}{
		"local config mapping exactly": {
			env: "local",
			expected: &configs.ApplicationProperties{
				AppName:      "go-playground",
				ServerConfig: configs.ServerConfig{Address: ":8080", ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 120 * time.Second, GraceTimeout: 20 * time.Second},
				AuthConfigs: []configs.AuthConfig{
					{Name: "tecchu11(ADMIN)", RoleStr: "ADMIN", Key: "admin"},
					{Name: "tecchu11(USER)", RoleStr: "USER", Key: "user"},
				},
			},
		},
		"error is not nil when given env is invalid": {
			env:       "none",
			expectErr: true,
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			actual, err := configs.Load(v.env)
			if v.expectErr {
				assert.NotNil(t, err, "expect error is not nil when given env is invalid")
				return
			}
			assert.NoError(t, err, "config mapping exactly so no err")
			assert.Equal(t, v.expected, actual, "config mapping exactly")
		})
	}
}
