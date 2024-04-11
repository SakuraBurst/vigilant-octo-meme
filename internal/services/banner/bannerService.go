package bannerservice

import (
	"context"
	"encoding/json"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"github.com/go-faster/errors"
)

type BannerStore interface {
	SaveBanner(banner *models.Banner) error
	UpdateBanner(id int, banner *models.Banner) error
	DeleteBanner(bannerId int) error
	GetUserBanner(bannerRequest *models.BannerRequest) (*models.Banner, error)
	GetAllBanners(bannerRequest *models.BannerRequest) ([]models.Banner, error)
}

type CacheStore interface {
	SetRequestCache(ctx context.Context, key, value []byte) error
	GetRequestCache(ctx context.Context, key []byte) ([]byte, error)
}

type TokenService interface {
	NewToken(isAdmin bool) (string, error)
	ParseToken(token string) (bool, error)
}

type Controller struct {
	bannerStore  BannerStore
	cacheStore   CacheStore
	tokenService TokenService
}

func (c *Controller) CreateNewBanner(banner *models.Banner, token string) error {
	isAdmin, err := c.tokenService.ParseToken(token)
	if err != nil {
		return errors.Wrap(err, "validate token failed")
	}
	if !isAdmin {
		return errors.New("user is not admin")
	}

	return c.bannerStore.SaveBanner(banner)
}

func (c *Controller) UpdateBannerById(id int, banner *models.Banner, token string) error {
	isAdmin, err := c.tokenService.ParseToken(token)
	if err != nil {
		return errors.Wrap(err, "validate token failed")
	}
	if !isAdmin {
		return errors.New("user is not admin")
	}

	return c.bannerStore.UpdateBanner(id, banner)
}

func (c *Controller) DeleteBannerById(bannerId int, token string) error {
	isAdmin, err := c.tokenService.ParseToken(token)
	if err != nil {
		return errors.Wrap(err, "validate token failed")
	}
	if !isAdmin {
		return errors.New("user is not admin")
	}

	return c.bannerStore.DeleteBanner(bannerId)
}

func (c *Controller) GetUserBanner(bannerRequest *models.BannerRequest, useLastRevision bool) ([]byte, error) {
	requestBytes, err := json.Marshal(bannerRequest)
	if !useLastRevision {
		if err != nil {
			return nil, errors.Wrap(err, "marshal request failed")
		}
		cache, err := c.cacheStore.GetRequestCache(context.TODO(), requestBytes)
		if err == nil {
			return cache, nil
		}
	}
	banner, err := c.bannerStore.GetUserBanner(bannerRequest)
	if err != nil {
		return nil, errors.Wrap(err, "get user banner failed")
	}
	bannerBytes, err := json.Marshal(banner)
	if err != nil {
		return nil, errors.Wrap(err, "marshal banner failed")

	}
	c.cacheStore.SetRequestCache(context.TODO(), requestBytes, bannerBytes)
	return bannerBytes, nil
}

func (c *Controller) GetAllBanners(bannerRequest *models.BannerRequest, token string) ([]models.Banner, error) {
	isAdmin, err := c.tokenService.ParseToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "validate token failed")
	}
	if !isAdmin {
		return nil, errors.New("user is not admin")
	}

	return c.bannerStore.GetAllBanners(bannerRequest)
}

func (c *Controller) CreateNewUserToken(isAdmin bool) (string, error) {
	token, err := c.tokenService.NewToken(isAdmin)
	if err != nil {
		return "", errors.Wrap(err, "generate token failed")
	}
	return token, nil
}

func NewBannerController(cfg *config.Config) *Controller {
	return &Controller{}
}
