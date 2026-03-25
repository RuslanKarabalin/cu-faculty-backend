package app

import (
	"faculty/internal/handler"
	"faculty/internal/repository"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3/client"
)

func (a *App) registerRoutes() error {
	repo := repository.Init(a.DB)

	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService, client.New(), a.Config.CuBaseUrl)

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
