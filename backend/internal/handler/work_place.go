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

type workPlaceService interface {
	GetWorkPlacesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.WorkPlace, int, error)
	CreateWorkPlace(ctx context.Context, userID uuid.UUID, req model.WorkPlaceRequest) (*model.WorkPlace, error)
	UpdateWorkPlace(ctx context.Context, userID uuid.UUID, id int, req model.WorkPlaceRequest) (*model.WorkPlace, error)
	DeleteWorkPlace(ctx context.Context, userID uuid.UUID, id int) error
}

type WorkPlaceHandler struct {
	service workPlaceService
	logger  *zap.Logger
}

func NewWorkPlaceHandler(service workPlaceService, logger *zap.Logger) *WorkPlaceHandler {
	return &WorkPlaceHandler{service: service, logger: logger}
}

func (h *WorkPlaceHandler) GetUserWorkPlaces(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid user id")
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	places, total, err := h.service.GetWorkPlacesByUserID(c.Context(), userID, limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get work places", err)
	}
	return c.JSON(model.Page[*model.WorkPlace]{Data: places, Total: total, Limit: limit, Offset: offset})
}

func (h *WorkPlaceHandler) GetMyWorkPlaces(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	places, total, err := h.service.GetWorkPlacesByUserID(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get work places", err)
	}
	return c.JSON(model.Page[*model.WorkPlace]{Data: places, Total: total, Limit: limit, Offset: offset})
}

func (h *WorkPlaceHandler) CreateWorkPlace(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.WorkPlaceRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	place, err := h.service.CreateWorkPlace(c.Context(), cuUser.ID, req)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to create work place", err)
	}
	return c.Status(fiber.StatusCreated).JSON(place)
}

func (h *WorkPlaceHandler) UpdateWorkPlace(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Params("workId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid work place id")
	}

	var req model.WorkPlaceRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	place, err := h.service.UpdateWorkPlace(c.Context(), cuUser.ID, id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "work place not found")
		}
		return unexpectedError(c, h.logger, "failed to update work place", err)
	}
	return c.JSON(place)
}

func (h *WorkPlaceHandler) DeleteWorkPlace(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Params("workId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid work place id")
	}

	if err := h.service.DeleteWorkPlace(c.Context(), cuUser.ID, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "work place not found")
		}
		return unexpectedError(c, h.logger, "failed to delete work place", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
