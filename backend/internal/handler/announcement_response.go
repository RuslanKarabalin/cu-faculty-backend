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

type announcementResponseService interface {
	Respond(ctx context.Context, userID, announcementID uuid.UUID) (*model.Announcement, error)
	DeleteResponse(ctx context.Context, userID, announcementID uuid.UUID) error
	GetResponders(ctx context.Context, announcementID uuid.UUID) ([]*model.User, error)
	GetMyResponses(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Announcement, int, error)
}

type AnnouncementResponseHandler struct {
	service announcementResponseService
	storage photoPresigner
	logger  *zap.Logger
}

func NewAnnouncementResponseHandler(service announcementResponseService, storage photoPresigner, logger *zap.Logger) *AnnouncementResponseHandler {
	return &AnnouncementResponseHandler{service: service, storage: storage, logger: logger}
}

func (h *AnnouncementResponseHandler) RespondToAnnouncement(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	announcementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid announcement id")
	}

	announcement, err := h.service.Respond(c.Context(), cuUser.ID, announcementID)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidRefID) || errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "announcement not found")
		}
		h.logger.Error("failed to respond to announcement", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	if announcement.Author != nil {
		announcement.Author.PhotoURL = presignPhoto(c.Context(), h.storage, h.logger, announcement.Author.PhotoS3Key)
	}
	return c.Status(fiber.StatusCreated).JSON(announcement)
}

func (h *AnnouncementResponseHandler) DeleteMyResponse(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	announcementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid announcement id")
	}

	if err := h.service.DeleteResponse(c.Context(), cuUser.ID, announcementID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "response not found")
		}
		h.logger.Error("failed to delete announcement response", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AnnouncementResponseHandler) GetResponders(c fiber.Ctx) error {
	announcementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid announcement id")
	}

	users, err := h.service.GetResponders(c.Context(), announcementID)
	if err != nil {
		h.logger.Error("failed to get announcement responders", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	for _, u := range users {
		u.PhotoURL = presignPhoto(c.Context(), h.storage, h.logger, u.PhotoS3Key)
	}
	return c.JSON(users)
}

func (h *AnnouncementResponseHandler) GetMyResponses(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}
	limit, offset := q.Normalize()

	announcements, total, err := h.service.GetMyResponses(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get announcement responses", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	for _, a := range announcements {
		if a.Author != nil {
			a.Author.PhotoURL = presignPhoto(c.Context(), h.storage, h.logger, a.Author.PhotoS3Key)
		}
	}
	return c.JSON(model.Page[*model.Announcement]{
		Data:   announcements,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}
