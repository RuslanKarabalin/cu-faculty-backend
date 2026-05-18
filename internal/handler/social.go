package handler

import (
	"context"
	"errors"
	"strconv"

	"faculty/internal/model"
	"faculty/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type socialService interface {
	GetSocialsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Social, error)
	CreateSocial(ctx context.Context, userID uuid.UUID, req model.SocialRequest) (*model.Social, error)
	UpdateSocial(ctx context.Context, userID uuid.UUID, id int, req model.SocialRequest) (*model.Social, error)
	DeleteSocial(ctx context.Context, userID uuid.UUID, id int) error
}

type SocialHandler struct {
	service socialService
	logger  *zap.Logger
}

func NewSocialHandler(service socialService, logger *zap.Logger) *SocialHandler {
	return &SocialHandler{service: service, logger: logger}
}

func (h *SocialHandler) GetUserSocials(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid user id")
	}

	socials, err := h.service.GetSocialsByUserID(c.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get socials", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(socials)
}

func (h *SocialHandler) GetMySocials(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	socials, err := h.service.GetSocialsByUserID(c.Context(), cuUser.ID)
	if err != nil {
		h.logger.Error("failed to get socials", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(socials)
}

func (h *SocialHandler) CreateSocial(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.SocialRequest
	if err := c.Bind().JSON(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	social, err := h.service.CreateSocial(c.Context(), cuUser.ID, req)
	if err != nil {
		h.logger.Error("failed to create social", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.Status(fiber.StatusCreated).JSON(social)
}

func (h *SocialHandler) UpdateSocial(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Params("socialId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid social id")
	}

	var req model.SocialRequest
	if err := c.Bind().JSON(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	social, err := h.service.UpdateSocial(c.Context(), cuUser.ID, id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "social not found")
		}
		h.logger.Error("failed to update social", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(social)
}

func (h *SocialHandler) DeleteSocial(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Params("socialId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid social id")
	}

	if err := h.service.DeleteSocial(c.Context(), cuUser.ID, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "social not found")
		}
		h.logger.Error("failed to delete social", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
