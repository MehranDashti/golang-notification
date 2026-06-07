package logger

import (
	"log/slog"
	"os"
	"strings"
)

func SetUp(env string, logLevel string) {
	level := parseLogLevel(logLevel)
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	}

	var handler slog.Handler
	if strings.ToLower(env) == "production" || strings.ToLower(env) == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
	slog.Info("Logger successfully configed",
		"level", level,
		"env", env,
		"format", formatName(handler),
	)
}

func parseLogLevel(logLevel string) slog.Level {
	switch strings.ToLower(logLevel) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func formatName(h slog.Handler) string {
	switch h.(type) {
	case *slog.JSONHandler:
		return "json"
	default:
		return "text"
	}
}
