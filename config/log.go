package config

import (
	"fmt"

	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("Failed to create zap loggger")
	}
	return logger
}
