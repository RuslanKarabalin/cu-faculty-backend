package app

import (
	"errors"
	"time"

	"faculty/internal/cuclient"
	"faculty/internal/handler"
	"faculty/internal/repository"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func (a *App) registerRoutes() {
	a.Fiber.Use(cors.New(cors.Config{
		AllowOrigins: a.Config.AllowedOrigins,
	}))

	a.Fiber.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	publicPaths := map[string]struct{}{
		"/health": {},
	}

	a.Fiber.Use(func(c *fiber.Ctx) error {
		if _, ok := publicPaths[c.Path()]; ok {
			return c.Next()
		}
		cookie := c.Cookies("bff.cookie")
		cuUser, err := a.CuClient.Authorize(c.Context(), cookie)
		if err != nil {
			if errors.Is(err, cuclient.ErrUnauthorized) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
			}
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "upstream error"})
		}
		c.Locals("cuUser", cuUser)
		return c.Next()
	})

	repo := repository.New(a.DB)
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService, a.Logger)

	a.Fiber.Get("/health", func(c *fiber.Ctx) error {
		if err := a.DB.Ping(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
			})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	a.Fiber.Post("/register", userHandler.Register)
	a.Fiber.Get("/users", userHandler.GetUsers)
}
