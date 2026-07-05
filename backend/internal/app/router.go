package app

import (
	"faculty/internal/apierr"
	"faculty/internal/handler"
	"faculty/internal/middleware"
	"faculty/internal/repository"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v3"
)

func (a *App) health(c fiber.Ctx) error {
	if err := a.DB.Ping(c.Context()); err != nil {
		return apierr.WriteCode(c, fiber.StatusServiceUnavailable, apierr.CodeUnavailable, "database is unavailable")
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

func (a *App) registerRoutes() {
	publicPaths := map[string]struct{}{
		"/health": {},
	}
	a.Fiber.Use(middleware.Auth(a.CuClient, publicPaths))

	repo := repository.New(a.DB)

	userHandler := handler.NewUserHandler(
		service.NewUserService(repo),
		service.NewRegistrationService(repo, a.CuClient),
		a.Storage,
		a.Logger,
	)
	eduPlaceHandler := handler.NewEduPlaceHandler(service.NewEduPlaceService(repo), a.Logger)
	workPlaceHandler := handler.NewWorkPlaceHandler(service.NewWorkPlaceService(repo), a.Logger)
	socialHandler := handler.NewSocialHandler(service.NewSocialService(repo), a.Logger)
	contactHandler := handler.NewContactHandler(service.NewContactService(repo), a.Storage, a.Logger)
	userKeySkillHandler := handler.NewUserKeySkillHandler(service.NewUserKeySkillService(repo), a.Logger)
	userSoftSkillHandler := handler.NewUserSoftSkillHandler(service.NewUserSoftSkillService(repo), a.Logger)
	referenceHandler := handler.NewReferenceHandler(service.NewReferenceService(repo), a.Logger)
	announcementHandler := handler.NewAnnouncementHandler(service.NewAnnouncementService(repo), a.Storage, a.Logger)
	newsHandler := handler.NewNewsHandler(service.NewNewsService(repo), a.Storage, a.Logger)
	eventHandler := handler.NewEventHandler(service.NewEventService(repo), a.Storage, a.Logger)
	announcementResponseHandler := handler.NewAnnouncementResponseHandler(service.NewAnnouncementResponseService(repo), a.Storage, a.Logger)
	eventResponseHandler := handler.NewEventResponseHandler(service.NewEventResponseService(repo), a.Storage, a.Logger)
	savedUserHandler := handler.NewSavedUserHandler(service.NewSavedUserService(repo), a.Storage, a.Logger)

	a.Fiber.Get("/health", a.health)

	api := a.Fiber.Group("/api")

	students := api.Group("/students")

	students.Post("/register", userHandler.Register)
	students.Get("/", userHandler.GetUsers)
	students.Get("/:id", userHandler.GetStudentByID)
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

	me.Post("/contacts", contactHandler.CreateContact)
	me.Get("/contacts", contactHandler.GetMyContacts)
	me.Put("/contacts/:contactId", contactHandler.UpdateContact)
	me.Delete("/contacts/:contactId", contactHandler.DeleteContact)

	me.Get("/announcements", announcementHandler.GetMyAnnouncements)
	me.Get("/news", newsHandler.GetMyNews)
	me.Get("/events", eventHandler.GetMyEvents)
	me.Get("/announcement-responses", announcementResponseHandler.GetMyResponses)
	me.Get("/event-responses", eventResponseHandler.GetMyResponses)

	me.Get("/key-skills", userKeySkillHandler.GetMyKeySkills)
	me.Post("/key-skills/:skillId", userKeySkillHandler.AddMyKeySkill)
	me.Delete("/key-skills/:skillId", userKeySkillHandler.DeleteMyKeySkill)

	me.Get("/soft-skills", userSoftSkillHandler.GetMySoftSkills)
	me.Post("/soft-skills/:skillId", userSoftSkillHandler.AddMySoftSkill)
	me.Delete("/soft-skills/:skillId", userSoftSkillHandler.DeleteMySoftSkill)

	me.Get("/saved-users", savedUserHandler.GetMySavedUsers)
	me.Post("/saved-users/:userId", savedUserHandler.AddMySavedUser)
	me.Delete("/saved-users/:userId", savedUserHandler.DeleteMySavedUser)

	announcements := api.Group("/announcements")
	announcements.Get("/", announcementHandler.GetAnnouncements)
	announcements.Post("/", announcementHandler.CreateAnnouncement)
	announcements.Get("/:id", announcementHandler.GetAnnouncementByID)
	announcements.Put("/:id", announcementHandler.UpdateAnnouncement)
	announcements.Delete("/:id", announcementHandler.DeleteAnnouncement)
	announcements.Get("/:id/responses", announcementResponseHandler.GetResponders)
	announcements.Post("/:id/responses", announcementResponseHandler.RespondToAnnouncement)
	announcements.Delete("/:id/responses", announcementResponseHandler.DeleteMyResponse)

	news := api.Group("/news")
	news.Get("/", newsHandler.GetNews)
	news.Post("/", newsHandler.CreateNews)
	news.Get("/:id", newsHandler.GetNewsByID)
	news.Put("/:id", newsHandler.UpdateNews)
	news.Delete("/:id", newsHandler.DeleteNews)

	events := api.Group("/events")
	events.Get("/", eventHandler.GetEvents)
	events.Post("/", eventHandler.CreateEvent)
	events.Get("/:id", eventHandler.GetEventByID)
	events.Put("/:id", eventHandler.UpdateEvent)
	events.Delete("/:id", eventHandler.DeleteEvent)
	events.Get("/:id/responses", eventResponseHandler.GetResponders)
	events.Post("/:id/responses", eventResponseHandler.RespondToEvent)
	events.Delete("/:id/responses", eventResponseHandler.DeleteMyResponse)

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
