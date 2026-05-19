package handler

import (
	"context"

	"faculty/internal/model"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type referenceService interface {
	GetStatuses(ctx context.Context) ([]*model.Status, error)
	GetKeySkills(ctx context.Context) ([]*model.Skill, error)
	GetSoftSkills(ctx context.Context) ([]*model.Skill, error)
	GetCompanies(ctx context.Context) ([]*model.Company, error)
	GetWorkPositions(ctx context.Context) ([]*model.WorkPosition, error)
	GetUniversities(ctx context.Context) ([]*model.University, error)
	GetFaqs(ctx context.Context) ([]*model.Faq, error)
	GetSocialNetworks(ctx context.Context) ([]string, error)
	GetEduGrades(ctx context.Context) ([]string, error)
	GetWorkGrades(ctx context.Context) ([]string, error)
	GetEventCategories(ctx context.Context) ([]string, error)
}

type ReferenceHandler struct {
	service referenceService
	logger  *zap.Logger
}

func NewReferenceHandler(service referenceService, logger *zap.Logger) *ReferenceHandler {
	return &ReferenceHandler{service: service, logger: logger}
}

func (h *ReferenceHandler) GetStatuses(c fiber.Ctx) error {
	items, err := h.service.GetStatuses(c.Context())
	if err != nil {
		h.logger.Error("failed to get statuses", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(items)
}

func (h *ReferenceHandler) GetKeySkills(c fiber.Ctx) error {
	items, err := h.service.GetKeySkills(c.Context())
	if err != nil {
		h.logger.Error("failed to get key skills", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(items)
}

func (h *ReferenceHandler) GetSoftSkills(c fiber.Ctx) error {
	items, err := h.service.GetSoftSkills(c.Context())
	if err != nil {
		h.logger.Error("failed to get soft skills", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(items)
}

func (h *ReferenceHandler) GetCompanies(c fiber.Ctx) error {
	items, err := h.service.GetCompanies(c.Context())
	if err != nil {
		h.logger.Error("failed to get companies", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(items)
}

func (h *ReferenceHandler) GetWorkPositions(c fiber.Ctx) error {
	items, err := h.service.GetWorkPositions(c.Context())
	if err != nil {
		h.logger.Error("failed to get work positions", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(items)
}

func (h *ReferenceHandler) GetUniversities(c fiber.Ctx) error {
	items, err := h.service.GetUniversities(c.Context())
	if err != nil {
		h.logger.Error("failed to get universities", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(items)
}

func (h *ReferenceHandler) GetFaqs(c fiber.Ctx) error {
	items, err := h.service.GetFaqs(c.Context())
	if err != nil {
		h.logger.Error("failed to get faqs", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(items)
}

func (h *ReferenceHandler) GetSocialNetworks(c fiber.Ctx) error {
	values, err := h.service.GetSocialNetworks(c.Context())
	if err != nil {
		h.logger.Error("failed to get social networks", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(values)
}

func (h *ReferenceHandler) GetEduGrades(c fiber.Ctx) error {
	values, err := h.service.GetEduGrades(c.Context())
	if err != nil {
		h.logger.Error("failed to get edu grades", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(values)
}

func (h *ReferenceHandler) GetWorkGrades(c fiber.Ctx) error {
	values, err := h.service.GetWorkGrades(c.Context())
	if err != nil {
		h.logger.Error("failed to get work grades", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(values)
}

func (h *ReferenceHandler) GetEventCategories(c fiber.Ctx) error {
	values, err := h.service.GetEventCategories(c.Context())
	if err != nil {
		h.logger.Error("failed to get event categories", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(values)
}
