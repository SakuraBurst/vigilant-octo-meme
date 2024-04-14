package cache

import (
	"context"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/eko/gocache/lib/v4/cache"
	gocacheStore "github.com/eko/gocache/store/go_cache/v4"
	gocache "github.com/patrickmn/go-cache"
	"time"
)

type Store struct {
	cacheManager *cache.Cache[[]byte]
}

func (s Store) SetRequestCache(ctx context.Context, key, value []byte) error {
	return s.cacheManager.Set(ctx, key, value)
}

func (s Store) GetRequestCache(ctx context.Context, key []byte) ([]byte, error) {
	return s.cacheManager.Get(ctx, key)
}

func New(cfg *config.Config) *Store {
	gocacheClient := gocache.New(cfg.App.CacheTTL, cfg.App.CacheTTL+5*time.Minute)
	gocacheStore := gocacheStore.NewGoCache(gocacheClient)
	cacheManager := cache.New[[]byte](gocacheStore)
	return &Store{cacheManager: cacheManager}
}
