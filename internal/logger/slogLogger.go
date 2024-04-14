package logger

import (
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	lj "gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func NewSlogLogger(cfg *config.Config) *slog.Logger {
	var logWriter io.Writer = os.Stdout
	if cfg.Env == "prod" {
		logFilePath := filepath.Join(".", "logs", "banners-log.jsonl")
		logWriter = &lj.Logger{
			Filename:   logFilePath,
			MaxBackups: 3,
			MaxSize:    1,
			MaxAge:     7,
		}
	}
	log := slog.New(slog.NewJSONHandler(logWriter, nil))
	return log
}
