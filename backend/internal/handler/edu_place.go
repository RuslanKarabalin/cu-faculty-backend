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

type eduPlaceService interface {
	GetEduPlacesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.EduPlace, int, error)
	CreateEduPlace(ctx context.Context, userID uuid.UUID, req model.EduPlaceRequest) (*model.EduPlace, error)
	UpdateEduPlace(ctx context.Context, userID uuid.UUID, id int, req model.EduPlaceRequest) (*model.EduPlace, error)
	DeleteEduPlace(ctx context.Context, userID uuid.UUID, id int) error
}

type EduPlaceHandler struct {
	service eduPlaceService
	logger  *zap.Logger
}

func NewEduPlaceHandler(service eduPlaceService, logger *zap.Logger) *EduPlaceHandler {
	return &EduPlaceHandler{service: service, logger: logger}
}

func (h *EduPlaceHandler) GetUserEduPlaces(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid user id")
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	places, total, err := h.service.GetEduPlacesByUserID(c.Context(), userID, limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get edu places", err)
	}
	return c.JSON(model.Page[*model.EduPlace]{Data: places, Total: total, Limit: limit, Offset: offset})
}

func (h *EduPlaceHandler) GetMyEduPlaces(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	places, total, err := h.service.GetEduPlacesByUserID(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get edu places", err)
	}
	return c.JSON(model.Page[*model.EduPlace]{Data: places, Total: total, Limit: limit, Offset: offset})
}

func (h *EduPlaceHandler) CreateEduPlace(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.EduPlaceRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	place, err := h.service.CreateEduPlace(c.Context(), cuUser.ID, req)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to create edu place", err)
	}
	return c.Status(fiber.StatusCreated).JSON(place)
}

func (h *EduPlaceHandler) UpdateEduPlace(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Params("eduId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid edu place id")
	}

	var req model.EduPlaceRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	place, err := h.service.UpdateEduPlace(c.Context(), cuUser.ID, id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "edu place not found")
		}
		return unexpectedError(c, h.logger, "failed to update edu place", err)
	}
	return c.JSON(place)
}

func (h *EduPlaceHandler) DeleteEduPlace(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Params("eduId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid edu place id")
	}

	if err := h.service.DeleteEduPlace(c.Context(), cuUser.ID, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "edu place not found")
		}
		return unexpectedError(c, h.logger, "failed to delete edu place", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
