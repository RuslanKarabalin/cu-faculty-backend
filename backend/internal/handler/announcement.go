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

type announcementService interface {
	GetAnnouncementByID(ctx context.Context, id, viewerID uuid.UUID) (*model.Announcement, error)
	GetAnnouncements(ctx context.Context, limit, offset int) ([]*model.Announcement, int, error)
	GetAnnouncementsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Announcement, int, error)
	CreateAnnouncement(ctx context.Context, authorID uuid.UUID, req model.CreateAnnouncementRequest) (*model.Announcement, error)
	UpdateAnnouncement(ctx context.Context, authorID, id uuid.UUID, req model.UpdateAnnouncementRequest) (*model.Announcement, error)
	DeleteAnnouncement(ctx context.Context, authorID, id uuid.UUID) error
}

type announcementPhotoStorage interface {
	PresignDownload(ctx context.Context, key string) (string, error)
}

type AnnouncementHandler struct {
	service announcementService
	storage announcementPhotoStorage
	logger  *zap.Logger
}

func NewAnnouncementHandler(service announcementService, storage announcementPhotoStorage, logger *zap.Logger) *AnnouncementHandler {
	return &AnnouncementHandler{service: service, storage: storage, logger: logger}
}

func (h *AnnouncementHandler) attachAuthorPhotoURL(ctx context.Context, u *model.User) {
	if u == nil || u.PhotoS3Key == nil {
		return
	}
	url, err := h.storage.PresignDownload(ctx, *u.PhotoS3Key)
	if err != nil {
		h.logger.Warn("failed to presign photo url", zap.Error(err))
		return
	}
	u.PhotoURL = &url
}

func (h *AnnouncementHandler) GetAnnouncements(c fiber.Ctx) error {
	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	announcements, total, err := h.service.GetAnnouncements(c.Context(), limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get announcements", err)
	}
	for _, a := range announcements {
		h.attachAuthorPhotoURL(c.Context(), a.Author)
	}
	return c.JSON(model.Page[*model.Announcement]{
		Data:   announcements,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *AnnouncementHandler) GetMyAnnouncements(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	announcements, total, err := h.service.GetAnnouncementsByAuthorID(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get announcements", err)
	}
	for _, a := range announcements {
		h.attachAuthorPhotoURL(c.Context(), a.Author)
	}
	return c.JSON(model.Page[*model.Announcement]{
		Data:   announcements,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *AnnouncementHandler) GetAnnouncementByID(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid announcement id")
	}

	announcement, err := h.service.GetAnnouncementByID(c.Context(), id, cuUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "announcement not found")
		}
		return unexpectedError(c, h.logger, "failed to get announcement by id", err)
	}
	h.attachAuthorPhotoURL(c.Context(), announcement.Author)
	return c.JSON(announcement)
}

func (h *AnnouncementHandler) CreateAnnouncement(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.CreateAnnouncementRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	announcement, err := h.service.CreateAnnouncement(c.Context(), cuUser.ID, req)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to create announcement", err)
	}
	h.attachAuthorPhotoURL(c.Context(), announcement.Author)
	return c.Status(fiber.StatusCreated).JSON(announcement)
}

func (h *AnnouncementHandler) UpdateAnnouncement(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid announcement id")
	}

	var req model.UpdateAnnouncementRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	announcement, err := h.service.UpdateAnnouncement(c.Context(), cuUser.ID, id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "announcement not found")
		}
		return unexpectedError(c, h.logger, "failed to update announcement", err)
	}
	h.attachAuthorPhotoURL(c.Context(), announcement.Author)
	return c.JSON(announcement)
}

func (h *AnnouncementHandler) DeleteAnnouncement(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid announcement id")
	}

	if err := h.service.DeleteAnnouncement(c.Context(), cuUser.ID, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "announcement not found")
		}
		return unexpectedError(c, h.logger, "failed to delete announcement", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
