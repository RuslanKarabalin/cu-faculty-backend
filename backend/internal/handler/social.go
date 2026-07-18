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
	GetSocialsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Social, int, error)
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
		return respondError(fiber.StatusBadRequest, "invalid user id")
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError()
	}
	limit, offset := q.Normalize()

	socials, total, err := h.service.GetSocialsByUserID(c.Context(), userID, limit, offset)
	if err != nil {
		return unexpectedError(h.logger, "failed to get socials", err)
	}
	return c.JSON(model.Page[*model.Social]{Data: socials, Total: total, Limit: limit, Offset: offset})
}

func (h *SocialHandler) GetMySocials(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError()
	}
	limit, offset := q.Normalize()

	socials, total, err := h.service.GetSocialsByUserID(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		return unexpectedError(h.logger, "failed to get socials", err)
	}
	return c.JSON(model.Page[*model.Social]{Data: socials, Total: total, Limit: limit, Offset: offset})
}

func (h *SocialHandler) CreateSocial(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.SocialRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	social, err := h.service.CreateSocial(c.Context(), cuUser.ID, req)
	if err != nil {
		return unexpectedError(h.logger, "failed to create social", err)
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
		return respondError(fiber.StatusBadRequest, "invalid social id")
	}

	var req model.SocialRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	social, err := h.service.UpdateSocial(c.Context(), cuUser.ID, id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "social not found")
		}
		return unexpectedError(h.logger, "failed to update social", err)
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
		return respondError(fiber.StatusBadRequest, "invalid social id")
	}

	if err := h.service.DeleteSocial(c.Context(), cuUser.ID, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "social not found")
		}
		return unexpectedError(h.logger, "failed to delete social", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
