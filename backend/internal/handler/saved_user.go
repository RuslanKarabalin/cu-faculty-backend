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

type savedUserService interface {
	GetSavedUsers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.User, int, error)
	AddSavedUser(ctx context.Context, userID, savedUserID uuid.UUID) (*model.User, error)
	DeleteSavedUser(ctx context.Context, userID, savedUserID uuid.UUID) error
}

type SavedUserHandler struct {
	service savedUserService
	storage photoPresigner
	logger  *zap.Logger
}

func NewSavedUserHandler(service savedUserService, storage photoPresigner, logger *zap.Logger) *SavedUserHandler {
	return &SavedUserHandler{service: service, storage: storage, logger: logger}
}

func (h *SavedUserHandler) attachPhotoURL(ctx context.Context, u *model.User) {
	if u == nil {
		return
	}
	u.PhotoURL = presignPhoto(ctx, h.storage, h.logger, u.PhotoS3Key)
}

func (h *SavedUserHandler) GetMySavedUsers(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError()
	}
	limit, offset := q.Normalize()

	users, total, err := h.service.GetSavedUsers(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		return unexpectedError(h.logger, "failed to get saved users", err)
	}
	for _, u := range users {
		h.attachPhotoURL(c.Context(), u)
	}
	return c.JSON(model.Page[*model.User]{Data: users, Total: total, Limit: limit, Offset: offset})
}

func (h *SavedUserHandler) AddMySavedUser(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	savedUserID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid user id")
	}

	if savedUserID == cuUser.ID {
		return respondError(fiber.StatusBadRequest, "cannot save yourself")
	}

	user, err := h.service.AddSavedUser(c.Context(), cuUser.ID, savedUserID)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidRefID) || errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "user not found")
		}
		return unexpectedError(h.logger, "failed to add saved user", err)
	}
	h.attachPhotoURL(c.Context(), user)
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *SavedUserHandler) DeleteMySavedUser(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	savedUserID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid user id")
	}

	if err := h.service.DeleteSavedUser(c.Context(), cuUser.ID, savedUserID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "saved user not found")
		}
		return unexpectedError(h.logger, "failed to delete saved user", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
