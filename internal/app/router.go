package app

import (
	"time"

	"faculty/internal/handler"
	"faculty/internal/repository"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func (a *App) registerRoutes() error {
	a.Fiber.Use(cors.New(cors.Config{
		AllowOrigins: a.Config.AllowedOrigins,
	}))

	a.Fiber.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	repo := repository.Init(a.DB)
	userService := service.NewUserService(repo)

	userHandler := handler.NewUserHandler(userService, a.CuClient, a.Logger)

	a.Fiber.Get("/health", func(c *fiber.Ctx) error {
		if err := a.DB.Ping(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
			})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	a.Fiber.Get("/me", userHandler.Me)
	a.Fiber.Get("/users", userHandler.GetUsers)

	return nil
}
