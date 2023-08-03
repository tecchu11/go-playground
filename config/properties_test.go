package config_test

import (
	"go-playground/config"
	"reflect"
	"testing"
)

func TestPropertiesLoader_Load(t *testing.T) {
	cases := []struct {
		name      string
		testdata  string
		expected  *config.Properties
		expectErr bool
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
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, _ := config.LoadConfigWith(c.testdata)
			if !reflect.DeepEqual(c.expected, actual) {
				t.Errorf("Failed to match between expected = %v and actual = %v with testdata = %v", c.expected, actual, c.testdata)
			}
		})
	}
}
