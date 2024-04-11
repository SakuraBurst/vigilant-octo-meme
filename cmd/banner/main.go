package main

import (
	"context"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/app"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/logger"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()
	log := logger.NewSlogLogger()
	log.LogAttrs(context.Background(), slog.LevelInfo, "cfg is", slog.Any("cfg", cfg))
	application, err := app.NewApp(cfg)
	if err != nil {
		panic(err)
	}
	if err := application.Run(); err != nil {
		panic(err)
	}
}
