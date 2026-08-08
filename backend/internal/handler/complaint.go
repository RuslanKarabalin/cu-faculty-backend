package handler

import (
	"context"
	"errors"

	"faculty/internal/model"
	"faculty/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type complaintService interface {
	ComplainAboutUser(ctx context.Context, complainant *model.CuUserResp, targetID uuid.UUID, req model.ComplaintRequest) error
	ComplainAboutPost(ctx context.Context, complainant *model.CuUserResp, targetID uuid.UUID, req model.ComplaintRequest) error
}

type ComplaintHandler struct {
	service    complaintService
	configured bool
	logger     *zap.Logger
}

func NewComplaintHandler(service complaintService, configured bool, logger *zap.Logger) *ComplaintHandler {
	return &ComplaintHandler{service: service, configured: configured, logger: logger}
}

func (h *ComplaintHandler) ComplainAboutUser(c fiber.Ctx) error {
	if !h.configured {
		return respondError(fiber.StatusServiceUnavailable, "complaints are not configured")
	}

	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	targetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid user id")
	}

	var req model.ComplaintRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	if targetID == cuUser.ID {
		return respondError(fiber.StatusBadRequest, "cannot report yourself")
	}

	if err := h.service.ComplainAboutUser(c.Context(), cuUser, targetID, req); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "user not found")
		}
		return unexpectedError(h.logger, "failed to send complaint", err)
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func (h *ComplaintHandler) ComplainAboutPost(c fiber.Ctx) error {
	if !h.configured {
		return respondError(fiber.StatusServiceUnavailable, "complaints are not configured")
	}

	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	targetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid post id")
	}

	var req model.ComplaintRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	if err := h.service.ComplainAboutPost(c.Context(), cuUser, targetID, req); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "post not found")
		}
		return unexpectedError(h.logger, "failed to send complaint", err)
	}
	return c.SendStatus(fiber.StatusAccepted)
}
