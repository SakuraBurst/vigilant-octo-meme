package bannerrouter

import (
	"encoding/json"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/constants"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"github.com/gofiber/fiber/v3"
	"strconv"
)

type BannerController interface {
	CreateNewBanner(banner *models.Banner, token string) error
	UpdateBannerByID(id int, banner *models.Banner, token string) error
	DeleteBannerByID(bannerID int, token string) error
	GetUserBanner(bannerRequest *models.BannerRequest, useLastRevision bool, token string) ([]byte, error)
	GetAllBanners(bannerRequest *models.BannerRequest, token string) ([]models.Banner, error)
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
		return ctx.SendStatus(fiber.StatusUnauthorized)
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
	result, err := r.controller.GetUserBanner(bannerRequest, useLastRevision, token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	ctx.Set("Content-Type", "application/json")
	return ctx.Send(result)
}

func (r *Router) GetAllBanners(ctx fiber.Ctx) error {
	token := ctx.Get("token")
	if token == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
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
	result, err := r.controller.GetAllBanners(bannerRequest, token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(result)
}
func (r *Router) DeleteBannerByID(ctx fiber.Ctx) error {
	token := ctx.Get("token")
	if token == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	err = r.controller.DeleteBannerByID(id, token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(fiber.StatusOK)
}
func (r *Router) UpdateBanner(ctx fiber.Ctx) error {
	token := ctx.Get("token")
	if token == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	banner := new(models.Banner)

	if err := json.Unmarshal(ctx.Body(), banner); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	err = r.controller.UpdateBannerByID(id, banner, token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(fiber.StatusOK)
}
func (r *Router) CreateNewBanner(ctx fiber.Ctx) error {
	token := ctx.Get("token")
	if token == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}
	banner := new(models.Banner)
	if err := json.Unmarshal(ctx.Body(), banner); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	err := r.controller.CreateNewBanner(banner, token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(fiber.StatusOK)
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
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"token": token})
}

func (r *Router) Run() error {
	return r.Listen(":" + r.port)
}

func New(cfg *config.Config, controller BannerController) *Router {
	app := fiber.New()
	r := &Router{App: app, port: cfg.App.Port, controller: controller}
	app.Get("/user_token", r.CreateUserToken)
	app.Get("/user_banner", r.GetUserBanner)
	bannerAPI := app.Group("/banner")
	bannerAPI.Patch("/:id", r.UpdateBanner)
	bannerAPI.Delete("/:id", r.GetUserBanner)
	bannerAPI.Get("", r.GetAllBanners)
	bannerAPI.Post("", r.CreateNewBanner)
	return r

}