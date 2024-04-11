package bannerrouter

import (
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"github.com/gofiber/fiber/v3"
)

type BannerController interface {
	CreateNewBanner(banner *models.Banner, token string) error
	UpdateBannerById(id int, banner *models.Banner, token string) error
	DeleteBannerById(bannerId int, token string) error
	GetUserBanner(bannerRequest *models.BannerRequest, useLastRevision bool) ([]byte, error)
	GetAllBanners(bannerRequest *models.BannerRequest, token string) ([]models.Banner, error)
	CreateNewUserToken(isAdmin bool) (string, error)
}

type Router struct {
	*fiber.App
	port       string
	controller BannerController
}

func (r *Router) GetUserBanner(c fiber.Ctx) error {
	token := c.Get("token")
	if token == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.SendString("zdarova")
}

func (r *Router) GetAllBanners(c fiber.Ctx) error {
	return c.SendString("zdarova")
}
func (r *Router) DeleteBannerById(c fiber.Ctx) error {
	return c.SendString("zdarova")
}
func (r *Router) UpdateBanner(c fiber.Ctx) error {
	return c.SendString("zdarova")
}
func (r *Router) CreateNewBanner(c fiber.Ctx) error {
	return c.SendString("zdarova")
}

func (r *Router) Run() error {
	return r.Listen(":" + r.port)
}

func New(cfg *config.Config, controller BannerController) *Router {
	app := fiber.New()
	r := &Router{App: app, port: cfg.App.Port, controller: controller}
	app.Get("/user_banner", r.GetUserBanner)
	bannerAPI := app.Group("/banner")
	bannerAPI.Patch("/:id", r.UpdateBanner)
	bannerAPI.Delete("/:id", r.GetUserBanner)
	bannerAPI.Get("", r.GetAllBanners)
	bannerAPI.Post("", r.CreateNewBanner)
	return r

}
