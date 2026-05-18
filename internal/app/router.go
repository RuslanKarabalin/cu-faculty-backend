package app

import (
	"faculty/internal/handler"
	"faculty/internal/middleware"
	"faculty/internal/repository"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func (a *App) registerRoutes() {
	a.Fiber.Use(cors.New(cors.Config{
		AllowOrigins:     a.Config.AllowedOrigins,
		AllowCredentials: true,
	}))

	publicPaths := map[string]struct{}{
		"/":       {},
		"/health": {},
	}
	a.Fiber.Use(middleware.Auth(a.CuClient, publicPaths))

	repo := repository.New(a.DB)
	userService := service.NewUserService(repo)
	eduPlaceService := service.NewEduPlaceService(repo)
	registrationService := service.NewRegistrationService(repo, a.CuClient)
	userHandler := handler.NewUserHandler(userService, eduPlaceService, registrationService, a.Logger)

	a.Fiber.Get("/", a.listRoutes)
	a.Fiber.Get("/health", a.health)

	api := a.Fiber.Group("/api")

	students := api.Group("/students")
	students.Post("/register", userHandler.Register)
	students.Get("/", userHandler.GetUsers)
	students.Get("/:id/edu-places", userHandler.GetUserEduPlaces)
}

func (a *App) listRoutes(c fiber.Ctx) error {
	routes := a.Fiber.GetRoutes(true)
	result := make([]map[string]string, 0, len(routes))
	for _, r := range routes {
		if r.Method == "HEAD" {
			continue
		}
		result = append(result, map[string]string{
			"method": r.Method,
			"path":   r.Path,
		})
	}
	return c.JSON(result)
}

func (a *App) health(c fiber.Ctx) error {
	if err := a.DB.Ping(c.Context()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "unhealthy"})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}
