package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func New() *Logger {
	var handler slog.Handler

	level := new(slog.LevelVar)

	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level.Set(slog.LevelDebug)
	default:
		level.Set(slog.LevelInfo)
	}
	
	handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return &Logger{
		Logger: slog.New(handler),
	}
}
