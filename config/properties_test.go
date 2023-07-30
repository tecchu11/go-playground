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
			name:     "case of successful loading with expected configuration json",
			testdata: "../testdata/config/config-test-1.json",
			expected: &config.Properties{
				AppName: "go-playground-test",
				AuthConfigs: []config.AuthConfig{
					{
						Name:    "test",
						RoleStr: "ADMIN",
						Key:     "test-api-key",
					},
				},
			},
		},
		{
			name:     "case to ensure that default values are storred when a required field is missing",
			testdata: "../testdata/config/config-test-2.json",
			expected: &config.Properties{
				AppName: "go-playground-test",
				AuthConfigs: []config.AuthConfig{
					{
						Name:    "test",
						RoleStr: "",
						Key:     "test-api-key",
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
