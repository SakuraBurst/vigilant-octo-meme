package bannerrouter

import (
	"context"
	"encoding/json"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/constants"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/services"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/storage"
	"github.com/go-faster/errors"
	"github.com/gofiber/fiber/v3"
	recoverMiddleware "github.com/gofiber/fiber/v3/middleware/recover"
	"strconv"
)

type BannerController interface {
	CreateNewBanner(ctx context.Context, banner *models.Banner, token string) (int, error)
	UpdateBannerByID(ctx context.Context, id int, banner *models.Banner, token string) error
	DeleteBannerByID(ctx context.Context, bannerID int, token string) error
	GetUserBanner(ctx context.Context, bannerRequest *models.BannerRequest, useLastRevision bool, token string) ([]byte, error)
	GetAllBanners(ctx context.Context, bannerRequest *models.BannerRequest, token string) ([]models.Banner, error)
	CreateNewUserToken(isAdmin bool) (string, error)
}

type Router struct {
	*fiber.App
	port       string
	controller BannerController
}

func (r *Router) GetUserBanner(ctx fiber.Ctx) error {
	token := ctx.Get("token")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token is required"})
	}
	tagID := constants.NoValue
	featureID := constants.NoValue
	useLastRevision := false
	var err error
	tagIDQuery := ctx.Query("tag_id")
	featureIDQuery := ctx.Query("feature_id")
	useLastRevisionQuery := ctx.Query("use_last_revision")
	if tagIDQuery == "" && featureIDQuery == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	if tagIDQuery != "" {
		tagID, err = strconv.Atoi(tagIDQuery)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
	}
	if featureIDQuery != "" {
		featureID, err = strconv.Atoi(featureIDQuery)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
	}

	if useLastRevisionQuery != "" {
		useLastRevision, err = strconv.ParseBool(useLastRevisionQuery)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
	}
	bannerRequest := &models.BannerRequest{
		TagID:     tagID,
		FeatureID: featureID,
	}
	result, err := r.controller.GetUserBanner(ctx.Context(), bannerRequest, useLastRevision, token)
	if err != nil {
		return returnError(ctx, err)
	}
	ctx.Set("Content-Type", "application/json")
	return ctx.Send(result)
}

func (r *Router) GetAllBanners(ctx fiber.Ctx) error {
	token := ctx.Get("token")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token is required"})
	}
	tagID := constants.NoValue
	featureID := constants.NoValue
	limit := constants.DefaultLimit
	offset := constants.DefaultOffset
	var err error
	tagIDQuery := ctx.Query("tag_id")
	featureIDQuery := ctx.Query("feature_id")
	limitQuery := ctx.Query("limit")
	offsetQuery := ctx.Query("offset")
	if tagIDQuery == "" && featureIDQuery == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	if tagIDQuery != "" {
		tagID, err = strconv.Atoi(tagIDQuery)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
	}
	if featureIDQuery != "" {
		featureID, err = strconv.Atoi(featureIDQuery)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
	}

	if limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
	}

	if offsetQuery != "" {
		offset, err = strconv.Atoi(offsetQuery)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
	}

	bannerRequest := &models.BannerRequest{
		TagID:     tagID,
		FeatureID: featureID,
		Limit:     limit,
		Offset:    offset,
	}
	result, err := r.controller.GetAllBanners(ctx.Context(), bannerRequest, token)
	if err != nil {
		return returnError(ctx, err)
	}
	return ctx.JSON(result)
}
func (r *Router) DeleteBannerByID(ctx fiber.Ctx) error {
	token := ctx.Get("token")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token is required"})
	}
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	err = r.controller.DeleteBannerByID(ctx.Context(), id, token)
	if err != nil {
		return returnError(ctx, err)
	}
	return ctx.SendStatus(fiber.StatusOK)
}
func (r *Router) UpdateBanner(ctx fiber.Ctx) error {
	token := ctx.Get("token")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token is required"})
	}
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	banner := new(models.Banner)

	if err := json.Unmarshal(ctx.Body(), banner); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	err = r.controller.UpdateBannerByID(ctx.Context(), id, banner, token)
	if err != nil {
		return returnError(ctx, err)
	}
	return ctx.SendStatus(fiber.StatusOK)
}
func (r *Router) CreateNewBanner(ctx fiber.Ctx) error {
	token := ctx.Get("token")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token is required"})
	}
	banner := new(models.Banner)
	if err := json.Unmarshal(ctx.Body(), banner); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	id, err := r.controller.CreateNewBanner(ctx.Context(), banner, token)
	if err != nil {
		return returnError(ctx, err)
	}
	ctx.Status(fiber.StatusCreated)
	return ctx.JSON(fiber.Map{"id": id})
}

func (r *Router) CreateUserToken(ctx fiber.Ctx) error {
	isAdmin := false
	var err error
	isAdminQuery := ctx.Query("isAdmin")
	if isAdminQuery != "" {
		isAdmin, err = strconv.ParseBool(ctx.Query("isAdmin"))
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
	}
	token, err := r.controller.CreateNewUserToken(isAdmin)
	if err != nil {
		return returnError(ctx, err)
	}
	return ctx.JSON(fiber.Map{"token": token})
}

func returnError(ctx fiber.Ctx, err error) error {
	if errors.Is(err, services.ErrTokenInvalid) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if errors.Is(err, services.ErrUserDontHaveAccess) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}
	if errors.Is(err, storage.BannerNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}

func (r *Router) Run() error {
	return r.Listen(":" + r.port)
}

func (r *Router) Close() error {
	return r.Shutdown()
}

func New(cfg *config.Config, controller BannerController) *Router {
	app := fiber.New()
	app.Use(recoverMiddleware.New())
	r := &Router{App: app, port: cfg.App.Port, controller: controller}
	app.Get("/user_token", r.CreateUserToken)
	app.Get("/user_banner", r.GetUserBanner)
	bannerAPI := app.Group("/banner")
	bannerAPI.Patch("/:id", r.UpdateBanner)
	bannerAPI.Delete("/:id", r.DeleteBannerByID)
	bannerAPI.Get("", r.GetAllBanners)
	bannerAPI.Post("", r.CreateNewBanner)
	return r

}
