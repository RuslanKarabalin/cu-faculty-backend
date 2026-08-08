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

type blockedUserService interface {
	GetBlockedUsers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.User, int, error)
	AddBlockedUser(ctx context.Context, userID, blockedUserID uuid.UUID) (*model.User, error)
	DeleteBlockedUser(ctx context.Context, userID, blockedUserID uuid.UUID) error
}

type BlockedUserHandler struct {
	service blockedUserService
	storage photoPresigner
	logger  *zap.Logger
}

func NewBlockedUserHandler(service blockedUserService, storage photoPresigner, logger *zap.Logger) *BlockedUserHandler {
	return &BlockedUserHandler{service: service, storage: storage, logger: logger}
}

func (h *BlockedUserHandler) attachPhotoURL(ctx context.Context, u *model.User) {
	if u == nil {
		return
	}
	u.PhotoURL = presignPhoto(ctx, h.storage, h.logger, u.PhotoS3Key)
}

func (h *BlockedUserHandler) GetMyBlockedUsers(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError()
	}
	limit, offset := q.Normalize()

	users, total, err := h.service.GetBlockedUsers(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		return unexpectedError(h.logger, "failed to get blocked users", err)
	}
	for _, u := range users {
		h.attachPhotoURL(c.Context(), u)
	}
	return c.JSON(model.Page[*model.User]{Data: users, Total: total, Limit: limit, Offset: offset})
}

func (h *BlockedUserHandler) AddMyBlockedUser(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	blockedUserID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid user id")
	}

	if blockedUserID == cuUser.ID {
		return respondError(fiber.StatusBadRequest, "cannot block yourself")
	}

	user, err := h.service.AddBlockedUser(c.Context(), cuUser.ID, blockedUserID)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidRefID) || errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "user not found")
		}
		return unexpectedError(h.logger, "failed to add blocked user", err)
	}
	h.attachPhotoURL(c.Context(), user)
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *BlockedUserHandler) DeleteMyBlockedUser(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	blockedUserID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid user id")
	}

	if err := h.service.DeleteBlockedUser(c.Context(), cuUser.ID, blockedUserID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "blocked user not found")
		}
		return unexpectedError(h.logger, "failed to delete blocked user", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
