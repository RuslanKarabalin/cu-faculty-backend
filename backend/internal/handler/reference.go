package handler

import (
	"context"

	"faculty/internal/model"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type referenceService interface {
	GetStatuses(ctx context.Context, limit, offset int) ([]*model.Status, int, error)
	GetKeySkills(ctx context.Context, limit, offset int) ([]*model.Skill, int, error)
	GetSoftSkills(ctx context.Context, limit, offset int) ([]*model.Skill, int, error)
	GetCompanies(ctx context.Context, limit, offset int) ([]*model.Company, int, error)
	GetWorkPositions(ctx context.Context, limit, offset int) ([]*model.WorkPosition, int, error)
	GetUniversities(ctx context.Context, limit, offset int) ([]*model.University, int, error)
	GetFaqs(ctx context.Context, limit, offset int) ([]*model.Faq, int, error)
	GetSocialNetworks(ctx context.Context, limit, offset int) ([]string, int, error)
	GetEduGrades(ctx context.Context, limit, offset int) ([]string, int, error)
	GetWorkGrades(ctx context.Context, limit, offset int) ([]string, int, error)
	GetEventCategories(ctx context.Context, limit, offset int) ([]string, int, error)
}

type ReferenceHandler struct {
	service referenceService
	logger  *zap.Logger
}

func NewReferenceHandler(service referenceService, logger *zap.Logger) *ReferenceHandler {
	return &ReferenceHandler{service: service, logger: logger}
}

func listReference[T any](c fiber.Ctx, h *ReferenceHandler, name string, fn func(ctx context.Context, limit, offset int) ([]T, int, error)) error {
	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	items, total, err := fn(c.Context(), limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get "+name, err)
	}
	return c.JSON(model.Page[T]{Data: items, Total: total, Limit: limit, Offset: offset})
}

func (h *ReferenceHandler) GetStatuses(c fiber.Ctx) error {
	return listReference(c, h, "statuses", h.service.GetStatuses)
}

func (h *ReferenceHandler) GetKeySkills(c fiber.Ctx) error {
	return listReference(c, h, "key skills", h.service.GetKeySkills)
}

func (h *ReferenceHandler) GetSoftSkills(c fiber.Ctx) error {
	return listReference(c, h, "soft skills", h.service.GetSoftSkills)
}

func (h *ReferenceHandler) GetCompanies(c fiber.Ctx) error {
	return listReference(c, h, "companies", h.service.GetCompanies)
}

func (h *ReferenceHandler) GetWorkPositions(c fiber.Ctx) error {
	return listReference(c, h, "work positions", h.service.GetWorkPositions)
}

func (h *ReferenceHandler) GetUniversities(c fiber.Ctx) error {
	return listReference(c, h, "universities", h.service.GetUniversities)
}

func (h *ReferenceHandler) GetFaqs(c fiber.Ctx) error {
	return listReference(c, h, "faqs", h.service.GetFaqs)
}

func (h *ReferenceHandler) GetSocialNetworks(c fiber.Ctx) error {
	return listReference(c, h, "social networks", h.service.GetSocialNetworks)
}

func (h *ReferenceHandler) GetEduGrades(c fiber.Ctx) error {
	return listReference(c, h, "edu grades", h.service.GetEduGrades)
}

func (h *ReferenceHandler) GetWorkGrades(c fiber.Ctx) error {
	return listReference(c, h, "work grades", h.service.GetWorkGrades)
}

func (h *ReferenceHandler) GetEventCategories(c fiber.Ctx) error {
	return listReference(c, h, "event categories", h.service.GetEventCategories)
}
