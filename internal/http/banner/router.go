package bannerrouter

import (
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"github.com/gofiber/fiber/v3"
)

type BannerController interface {
	GetBanner() (models.Banner, error)
	GetAllBanners() ([]models.Banner, error)
	CreateNewBanner() error
}

type Router struct {
	*fiber.App
	port string
}

func (r *Router) GetUserBanner(c fiber.Ctx) error {
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

func New(cfg *config.Config) *Router {
	app := fiber.New()
	r := &Router{App: app, port: cfg.App.Port}
	app.Get("/user_banner", r.GetUserBanner)
	bannerAPI := app.Group("/banner")
	bannerAPI.Patch("/:id", r.UpdateBanner)
	bannerAPI.Delete("/:id", r.GetUserBanner)
	bannerAPI.Get("", r.GetAllBanners)
	bannerAPI.Post("", r.CreateNewBanner)
	return r

}
