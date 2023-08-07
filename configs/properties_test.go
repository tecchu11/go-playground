package configs_test

import (
	"go-playground/configs"
	"reflect"
	"testing"
	"time"
)

func TestPropertiesLoader_Load(t *testing.T) {
	cases := []struct {
		name      string
		env       string
		expected  *configs.ApplicationProperties
		expectErr bool
	}{
		{
			name: "case of successful loading local configuration",
			env:  "local",
			expected: &configs.ApplicationProperties{
				AppName: "go-playground",
				ServerConfig: configs.ServerConfig{
					Address:      ":8080",
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  120 * time.Second,
				},
				AuthConfigs: []configs.AuthConfig{
					{
						Name:    "tecchu11(ADMIN)",
						RoleStr: "ADMIN",
						Key:     "admin",
					},
					{
						Name:    "tecchu11(USER)",
						RoleStr: "USER",
						Key:     "user",
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, _ := configs.Load(c.env)
			if !reflect.DeepEqual(c.expected, actual) {
				t.Errorf("Failed to match between expected = %v and actual = %v with testdata = %v", c.expected, actual, c.env)
			}
		})
	}
}
