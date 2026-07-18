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

type postService interface {
	GetPostByID(ctx context.Context, id, viewerID uuid.UUID) (*model.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]*model.Post, int, error)
	GetPostsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Post, int, error)
	CreatePost(ctx context.Context, authorID uuid.UUID, req model.CreatePostRequest) (*model.Post, error)
	UpdatePost(ctx context.Context, authorID, id uuid.UUID, req model.UpdatePostRequest) (*model.Post, error)
	DeletePost(ctx context.Context, authorID, id uuid.UUID) error
}

type postPhotoStorage interface {
	PresignDownload(ctx context.Context, key string) (string, error)
}

type PostHandler struct {
	service postService
	storage postPhotoStorage
	logger  *zap.Logger
}

func NewPostHandler(service postService, storage postPhotoStorage, logger *zap.Logger) *PostHandler {
	return &PostHandler{service: service, storage: storage, logger: logger}
}

func (h *PostHandler) attachAuthorPhotoURL(ctx context.Context, u *model.User) {
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

func (h *PostHandler) GetPosts(c fiber.Ctx) error {
	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError()
	}
	limit, offset := q.Normalize()

	posts, total, err := h.service.GetPosts(c.Context(), limit, offset)
	if err != nil {
		return unexpectedError(h.logger, "failed to get posts", err)
	}
	for _, a := range posts {
		h.attachAuthorPhotoURL(c.Context(), a.Author)
	}
	return c.JSON(model.Page[*model.Post]{
		Data:   posts,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *PostHandler) GetMyPosts(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError()
	}
	limit, offset := q.Normalize()

	posts, total, err := h.service.GetPostsByAuthorID(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		return unexpectedError(h.logger, "failed to get posts", err)
	}
	for _, a := range posts {
		h.attachAuthorPhotoURL(c.Context(), a.Author)
	}
	return c.JSON(model.Page[*model.Post]{
		Data:   posts,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *PostHandler) GetPostByID(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid post id")
	}

	post, err := h.service.GetPostByID(c.Context(), id, cuUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "post not found")
		}
		return unexpectedError(h.logger, "failed to get post by id", err)
	}
	h.attachAuthorPhotoURL(c.Context(), post.Author)
	return c.JSON(post)
}

func (h *PostHandler) CreatePost(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.CreatePostRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	post, err := h.service.CreatePost(c.Context(), cuUser.ID, req)
	if err != nil {
		return unexpectedError(h.logger, "failed to create post", err)
	}
	h.attachAuthorPhotoURL(c.Context(), post.Author)
	return c.Status(fiber.StatusCreated).JSON(post)
}

func (h *PostHandler) UpdatePost(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid post id")
	}

	var req model.UpdatePostRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	post, err := h.service.UpdatePost(c.Context(), cuUser.ID, id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "post not found")
		}
		return unexpectedError(h.logger, "failed to update post", err)
	}
	h.attachAuthorPhotoURL(c.Context(), post.Author)
	return c.JSON(post)
}

func (h *PostHandler) DeletePost(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid post id")
	}

	if err := h.service.DeletePost(c.Context(), cuUser.ID, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "post not found")
		}
		return unexpectedError(h.logger, "failed to delete post", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
