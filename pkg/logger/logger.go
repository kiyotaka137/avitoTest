package logger

import (
	"avitoTest/internal/config"
	"log/slog"
	"os"
)

var Log *slog.Logger

func InitLogger(cfg *config.Config) {
	var level slog.Level

	switch cfg.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	Log = slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: level,
			},
		),
	)
}
