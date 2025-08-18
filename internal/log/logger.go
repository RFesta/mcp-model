package log

import (
	"log/slog"
	"os"
)

func New(env, service string) *slog.Logger {
	opts := &slog.HandlerOptions{AddSource: false, Level: slog.LevelInfo}
	var h slog.Handler = slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(h)
	logger = logger.With("service", service, "env", env)
	return logger
}
