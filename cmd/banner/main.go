package main

import (
	"context"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/app/router"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/logger"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()
	log := logger.NewSlogLogger()
	log.LogAttrs(context.Background(), slog.LevelInfo, "cfg is", slog.Any("cfg", cfg))
	rt := router.New(cfg)
	rt.Run()

}
