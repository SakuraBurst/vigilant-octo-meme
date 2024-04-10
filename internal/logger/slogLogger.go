package logger

import (
	"log/slog"
	"os"
)

func NewSlogLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	return log
}
