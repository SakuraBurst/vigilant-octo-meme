package main

import (
	"github.com/SakuraBurst/vigilant-octo-meme/internal/app"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.NewSlogLogger(cfg)
	application, err := app.NewApp(cfg, log)
	if err != nil {
		panic(err)
	}
	if err := application.Run(); err != nil {
		panic(err)
	}
}
