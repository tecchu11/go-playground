package config_test

import (
	"go-playground/cmd/api/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	os.Setenv("DB_USER", "test-user")
	os.Setenv("DB_PASSWORD", "test-password")
	os.Setenv("DB_HOST", "test.localhost:3306")
	os.Setenv("DB_NAME", "test-db")
	actual := config.Load()
	expected := config.Configuration{
		Env:        "test",
		DBUser:     "test-user",
		DBPassword: "test-password",
		DBAddr:     "test.localhost:3306",
		DBName:     "test-db",
	}
	assert.Equal(t, &expected, actual)
}
