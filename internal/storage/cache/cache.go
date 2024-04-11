package cache

import (
	"context"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	gocache "github.com/patrickmn/go-cache"
	"time"
)

type CacheStore struct {
	cacheManager *cache.Cache[[]byte]
}

func (s CacheStore) SetRequestCache(ctx context.Context, key, value []byte) error {
	return s.cacheManager.Set(ctx, key, value)
}

func (s CacheStore) GetRequestCache(ctx context.Context, key []byte) ([]byte, error) {
	return s.cacheManager.Get(ctx, key)
}

func New(cfg *config.Config) *CacheStore {
	gocacheClient := gocache.New(cfg.App.CacheTTL, cfg.App.CacheTTL+5*time.Minute)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)
	cacheManager := cache.New[[]byte](gocacheStore)
	return &CacheStore{cacheManager: cacheManager}
}
