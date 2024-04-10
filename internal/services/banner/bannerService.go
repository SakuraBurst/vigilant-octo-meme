package bannerservice

import (
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"time"
)

type BannerRepository interface {
	SaveBanner(*models.Banner) error
	GetBanner() (*models.Banner, error)
	GetAllBanners() ([]*models.Banner, error)
	DeleteBanner() error
}

type Controller struct {
	cacheTTL time.Duration
}

func NewBannerController(cfg *config.Config) *Controller {
	return &Controller{cacheTTL: cfg.App.CacheTTL}
}
