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

type postResponseService interface {
	Respond(ctx context.Context, userID, postID uuid.UUID) (*model.Post, error)
	DeleteResponse(ctx context.Context, userID, postID uuid.UUID) error
	GetResponders(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*model.User, int, error)
	GetMyResponses(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Post, int, error)
}

type PostResponseHandler struct {
	service postResponseService
	storage photoPresigner
	logger  *zap.Logger
}

func NewPostResponseHandler(service postResponseService, storage photoPresigner, logger *zap.Logger) *PostResponseHandler {
	return &PostResponseHandler{service: service, storage: storage, logger: logger}
}

func (h *PostResponseHandler) RespondToPost(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid post id")
	}

	post, err := h.service.Respond(c.Context(), cuUser.ID, postID)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidRefID) || errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "post not found")
		}
		return unexpectedError(h.logger, "failed to respond to post", err)
	}
	if post.Author != nil {
		post.Author.PhotoURL = presignPhoto(c.Context(), h.storage, h.logger, post.Author.PhotoS3Key)
	}
	return c.Status(fiber.StatusCreated).JSON(post)
}

func (h *PostResponseHandler) DeleteMyResponse(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid post id")
	}

	if err := h.service.DeleteResponse(c.Context(), cuUser.ID, postID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "response not found")
		}
		return unexpectedError(h.logger, "failed to delete post response", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *PostResponseHandler) GetResponders(c fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid post id")
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError()
	}
	limit, offset := q.Normalize()

	users, total, err := h.service.GetResponders(c.Context(), postID, limit, offset)
	if err != nil {
		return unexpectedError(h.logger, "failed to get post responders", err)
	}
	for _, u := range users {
		u.PhotoURL = presignPhoto(c.Context(), h.storage, h.logger, u.PhotoS3Key)
	}
	return c.JSON(model.Page[*model.User]{Data: users, Total: total, Limit: limit, Offset: offset})
}

func (h *PostResponseHandler) GetMyResponses(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError()
	}
	limit, offset := q.Normalize()

	posts, total, err := h.service.GetMyResponses(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		return unexpectedError(h.logger, "failed to get post responses", err)
	}
	for _, a := range posts {
		if a.Author != nil {
			a.Author.PhotoURL = presignPhoto(c.Context(), h.storage, h.logger, a.Author.PhotoS3Key)
		}
	}
	return c.JSON(model.Page[*model.Post]{
		Data:   posts,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}
