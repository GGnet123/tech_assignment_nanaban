package logger

import (
	"go.uber.org/zap"
)

func New(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	default:
		return zap.NewDevelopment()
	}
}
