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
	GetEduPlacesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EduPlace, error)
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

	places, err := h.service.GetEduPlacesByUserID(c.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get edu places", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(places)
}

func (h *EduPlaceHandler) GetMyEduPlaces(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	places, err := h.service.GetEduPlacesByUserID(c.Context(), cuUser.ID)
	if err != nil {
		h.logger.Error("failed to get edu places", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(places)
}

func (h *EduPlaceHandler) CreateEduPlace(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.EduPlaceRequest
	if err := c.Bind().JSON(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	place, err := h.service.CreateEduPlace(c.Context(), cuUser.ID, req)
	if err != nil {
		h.logger.Error("failed to create edu place", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
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
	if err := c.Bind().JSON(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	place, err := h.service.UpdateEduPlace(c.Context(), cuUser.ID, id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "edu place not found")
		}
		h.logger.Error("failed to update edu place", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
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
		h.logger.Error("failed to delete edu place", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
