package config_test

import (
	"go-playground/config"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestPropertiesLoader_Load(t *testing.T) {
	cases := []struct {
		name     string
		testdata string
		expected *config.Properties
	}{
		{
			name:     "case of successful loading local configuration",
			testdata: "config-local.json",
			expected: &config.Properties{
				AppName: "go-playground",
				AuthConfigs: []config.AuthConfig{
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
	loader := config.NewPropertiesLoader(zap.NewExample())
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := loader.Load(c.testdata)
			if ok := reflect.DeepEqual(c.expected, actual); !ok {
				t.Errorf("Failed to match between expected = %v and actual = %v with testdata = %v", c.expected, actual, c.testdata)
			}
		})
	}
}
