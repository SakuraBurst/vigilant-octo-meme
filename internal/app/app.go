package app

import (
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	bannerrouter "github.com/SakuraBurst/vigilant-octo-meme/internal/http/banner"
	bannerservice "github.com/SakuraBurst/vigilant-octo-meme/internal/services/banner"
	jwtservice "github.com/SakuraBurst/vigilant-octo-meme/internal/services/jwt"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/storage/cache"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/storage/postgres"
	"io"
	"log/slog"
)

type Router interface {
	Run() error
	Close() error
}

type App struct {
	router  Router
	log     *slog.Logger
	storage io.Closer
}

func (a *App) Run() error {
	return a.router.Run()
}

func (a *App) Stop() {
	log := a.log.With(slog.String("method", "GracefulShutdown"))
	err := a.router.Close()
	if err != nil {
		log.Error(err.Error())
	}
	err = a.storage.Close()
	if err != nil {
		log.Error(err.Error())
	}
}

func NewApp(cfg *config.Config, log *slog.Logger) (*App, error) {
	storage, err := postgres.New(cfg, log)
	if err != nil {
		return nil, err
	}

	cacheStore := cache.New(cfg)
	tokenService := jwtservice.New(cfg, log)
	bannerService := bannerservice.New(storage, cacheStore, tokenService, log)
	router := bannerrouter.New(cfg, bannerService)

	return &App{router: router, storage: storage}, nil
}
