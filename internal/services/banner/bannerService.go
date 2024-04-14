package bannerservice

import (
	"context"
	"encoding/json"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/constants"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/services"
	"github.com/go-faster/errors"
	"log/slog"
)

type BannerStore interface {
	SaveBanner(ctx context.Context, banner *models.Banner) (int, error)
	UpdateBanner(ctx context.Context, id int, banner *models.Banner) error
	DeleteBanner(ctx context.Context, bannerID int) error
	GetUserBanner(ctx context.Context, bannerRequest *models.BannerRequest) (map[string]interface{}, error)
	GetAllBanners(ctx context.Context, bannerRequest *models.BannerRequest) ([]models.Banner, error)
}

type CacheStore interface {
	SetRequestCache(ctx context.Context, key, value []byte) error
	GetRequestCache(ctx context.Context, key []byte) ([]byte, error)
}

type TokenService interface {
	NewToken(isAdmin bool) (string, error)
	ParseToken(token string) (bool, error)
}

type Service struct {
	bannerStore  BannerStore
	cacheStore   CacheStore
	tokenService TokenService
	log          *slog.Logger
}

func (c *Service) CreateNewBanner(ctx context.Context, banner *models.Banner, token string) (int, error) {
	log := c.log.With(slog.String("method", "CreateNewBanner"))
	isAdmin, err := c.tokenService.ParseToken(token)
	if err != nil {
		log.Error(err.Error())
		return constants.NoValue, errors.Wrap(err, "validate token failed")
	}
	if !isAdmin {
		log.Error("user is not admin")
		return constants.NoValue, errors.Wrap(services.ErrUserDontHaveAccess, "user is not admin")
	}
	return c.bannerStore.SaveBanner(ctx, banner)
}

func (c *Service) UpdateBannerByID(ctx context.Context, id int, banner *models.Banner, token string) error {
	log := c.log.With(slog.String("method", "UpdateBannerByID"))
	isAdmin, err := c.tokenService.ParseToken(token)
	if err != nil {
		log.Error(err.Error())
		return errors.Wrap(err, "validate token failed")
	}
	if !isAdmin {
		log.Error("user is not admin")
		return errors.Wrap(services.ErrUserDontHaveAccess, "user is not admin")
	}

	return c.bannerStore.UpdateBanner(ctx, id, banner)
}

func (c *Service) DeleteBannerByID(ctx context.Context, bannerID int, token string) error {
	log := c.log.With(slog.String("method", "DeleteBannerByID"))
	isAdmin, err := c.tokenService.ParseToken(token)
	if err != nil {
		log.Error(err.Error())
		return errors.Wrap(err, "validate token failed")
	}
	if !isAdmin {
		log.Error("user is not admin")
		return errors.Wrap(services.ErrUserDontHaveAccess, "user is not admin")
	}

	return c.bannerStore.DeleteBanner(ctx, bannerID)
}

func (c *Service) GetUserBanner(ctx context.Context, bannerRequest *models.BannerRequest, useLastRevision bool, token string) ([]byte, error) {
	log := c.log.With(slog.String("method", "GetUserBanner"))
	_, err := c.tokenService.ParseToken(token)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.Wrap(err, "validate token failed")
	}
	requestBytes, err := json.Marshal(bannerRequest)
	if err != nil {
		return nil, errors.Wrap(err, "marshal request failed")
	}
	if !useLastRevision {
		cache, err := c.cacheStore.GetRequestCache(ctx, requestBytes)
		if err == nil {
			return cache, nil
		}
		log.Error(err.Error())
	}
	banner, err := c.bannerStore.GetUserBanner(ctx, bannerRequest)
	if err != nil {
		return nil, errors.Wrap(err, "get user banner failed")
	}
	bannerBytes, err := json.Marshal(banner)
	if err != nil {
		return nil, errors.Wrap(err, "marshal banner failed")
	}
	err = c.cacheStore.SetRequestCache(context.TODO(), requestBytes, bannerBytes)
	if err != nil {
		return nil, errors.Wrap(err, "set cache failed")
	}
	return bannerBytes, nil
}

func (c *Service) GetAllBanners(ctx context.Context, bannerRequest *models.BannerRequest, token string) ([]models.Banner, error) {
	isAdmin, err := c.tokenService.ParseToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "validate token failed")
	}
	if !isAdmin {
		return nil, errors.Wrap(services.ErrUserDontHaveAccess, "user is not admin")
	}

	return c.bannerStore.GetAllBanners(ctx, bannerRequest)
}

func (c *Service) CreateNewUserToken(isAdmin bool) (string, error) {
	token, err := c.tokenService.NewToken(isAdmin)
	if err != nil {
		return "", errors.Wrap(err, "generate token failed")
	}
	return token, nil
}

func New(bannerStore BannerStore, cacheStore CacheStore, tokenService TokenService, log *slog.Logger) *Service {
	return &Service{bannerStore: bannerStore, cacheStore: cacheStore, tokenService: tokenService, log: log}
}
