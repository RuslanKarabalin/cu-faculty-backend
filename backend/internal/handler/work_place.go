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
	GetWorkPlacesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.WorkPlace, error)
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

	places, err := h.service.GetWorkPlacesByUserID(c.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get work places", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(places)
}

func (h *WorkPlaceHandler) GetMyWorkPlaces(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	places, err := h.service.GetWorkPlacesByUserID(c.Context(), cuUser.ID)
	if err != nil {
		h.logger.Error("failed to get work places", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(places)
}

func (h *WorkPlaceHandler) CreateWorkPlace(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.WorkPlaceRequest
	if err := c.Bind().JSON(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	place, err := h.service.CreateWorkPlace(c.Context(), cuUser.ID, req)
	if err != nil {
		h.logger.Error("failed to create work place", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
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
	if err := c.Bind().JSON(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	place, err := h.service.UpdateWorkPlace(c.Context(), cuUser.ID, id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "work place not found")
		}
		h.logger.Error("failed to update work place", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
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
		h.logger.Error("failed to delete work place", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
