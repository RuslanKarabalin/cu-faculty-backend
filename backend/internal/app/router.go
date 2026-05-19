package app

import (
	"faculty/internal/handler"
	"faculty/internal/middleware"
	"faculty/internal/repository"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

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

	userHandler := handler.NewUserHandler(
		service.NewUserService(repo),
		service.NewRegistrationService(repo, a.CuClient),
		a.Logger,
	)
	eduPlaceHandler := handler.NewEduPlaceHandler(service.NewEduPlaceService(repo), a.Logger)
	workPlaceHandler := handler.NewWorkPlaceHandler(service.NewWorkPlaceService(repo), a.Logger)
	socialHandler := handler.NewSocialHandler(service.NewSocialService(repo), a.Logger)
	userKeySkillHandler := handler.NewUserKeySkillHandler(service.NewUserKeySkillService(repo), a.Logger)
	userSoftSkillHandler := handler.NewUserSoftSkillHandler(service.NewUserSoftSkillService(repo), a.Logger)
	referenceHandler := handler.NewReferenceHandler(service.NewReferenceService(repo), a.Logger)

	a.Fiber.Get("/", a.listRoutes)
	a.Fiber.Get("/health", a.health)

	api := a.Fiber.Group("/api")

	students := api.Group("/students")

	students.Post("/register", userHandler.Register)
	students.Get("/", userHandler.GetUsers)
	students.Get("/:id/edu-places", eduPlaceHandler.GetUserEduPlaces)
	students.Get("/:id/work-places", workPlaceHandler.GetUserWorkPlaces)
	students.Get("/:id/socials", socialHandler.GetUserSocials)
	students.Get("/:id/key-skills", userKeySkillHandler.GetUserKeySkills)
	students.Get("/:id/soft-skills", userSoftSkillHandler.GetUserSoftSkills)

	me := api.Group("/me")

	me.Get("/", userHandler.GetMe)
	me.Put("/", userHandler.UpdateMe)

	me.Post("/edu-places", eduPlaceHandler.CreateEduPlace)
	me.Get("/edu-places", eduPlaceHandler.GetMyEduPlaces)
	me.Put("/edu-places/:eduId", eduPlaceHandler.UpdateEduPlace)
	me.Delete("/edu-places/:eduId", eduPlaceHandler.DeleteEduPlace)

	me.Post("/work-places", workPlaceHandler.CreateWorkPlace)
	me.Get("/work-places", workPlaceHandler.GetMyWorkPlaces)
	me.Put("/work-places/:workId", workPlaceHandler.UpdateWorkPlace)
	me.Delete("/work-places/:workId", workPlaceHandler.DeleteWorkPlace)

	me.Post("/socials", socialHandler.CreateSocial)
	me.Get("/socials", socialHandler.GetMySocials)
	me.Put("/socials/:socialId", socialHandler.UpdateSocial)
	me.Delete("/socials/:socialId", socialHandler.DeleteSocial)

	me.Get("/key-skills", userKeySkillHandler.GetMyKeySkills)
	me.Post("/key-skills/:skillId", userKeySkillHandler.AddMyKeySkill)
	me.Delete("/key-skills/:skillId", userKeySkillHandler.DeleteMyKeySkill)

	me.Get("/soft-skills", userSoftSkillHandler.GetMySoftSkills)
	me.Post("/soft-skills/:skillId", userSoftSkillHandler.AddMySoftSkill)
	me.Delete("/soft-skills/:skillId", userSoftSkillHandler.DeleteMySoftSkill)

	api.Get("/statuses", referenceHandler.GetStatuses)
	api.Get("/key-skills", referenceHandler.GetKeySkills)
	api.Get("/soft-skills", referenceHandler.GetSoftSkills)
	api.Get("/companies", referenceHandler.GetCompanies)
	api.Get("/work-positions", referenceHandler.GetWorkPositions)
	api.Get("/universities", referenceHandler.GetUniversities)
	api.Get("/faqs", referenceHandler.GetFaqs)
	api.Get("/social-networks", referenceHandler.GetSocialNetworks)
	api.Get("/edu-grades", referenceHandler.GetEduGrades)
	api.Get("/work-grades", referenceHandler.GetWorkGrades)
	api.Get("/event-categories", referenceHandler.GetEventCategories)
}
