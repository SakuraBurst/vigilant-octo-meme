package app

import (
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	bannerrouter "github.com/SakuraBurst/vigilant-octo-meme/internal/http/banner"
	bannerservice "github.com/SakuraBurst/vigilant-octo-meme/internal/services/banner"
	jwtservice "github.com/SakuraBurst/vigilant-octo-meme/internal/services/jwt"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/storage/cache"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/storage/postgres"
	"log/slog"
)

type Router interface {
	Run() error
}

type App struct {
	router Router
}

func (a *App) Run() error {
	return a.router.Run()
}

func NewApp(cfg *config.Config, log *slog.Logger) (*App, error) {
	storage, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	cacheStore := cache.New(cfg)
	tokenService := jwtservice.New(cfg)
	bannerService := bannerservice.New(storage, cacheStore, tokenService, log)
	router := bannerrouter.New(cfg, bannerService)

	return &App{router: router}, nil
}
