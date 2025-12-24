package logger

import (
	"log/slog"
	"os"
)

func InitLogger(level string) {
	var logLevel slog.Level

	switch level {
	case "DEBUG":
		logLevel = slog.LevelDebug

	case "INFO":
		logLevel = slog.LevelInfo

	case "WARN":
		logLevel = slog.LevelWarn

	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)
}
