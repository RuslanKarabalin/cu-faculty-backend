package handler

import (
	"context"
	"errors"
	"io"

	"faculty/internal/model"
	"faculty/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type newsService interface {
	GetNewsByID(ctx context.Context, id, viewerID uuid.UUID) (*model.News, error)
	GetNews(ctx context.Context, limit, offset int) ([]*model.News, int, error)
	GetNewsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.News, int, error)
	CreateNews(ctx context.Context, authorID uuid.UUID, req model.CreateNewsRequest) (*model.News, error)
	UpdateNews(ctx context.Context, authorID, id uuid.UUID, req model.UpdateNewsRequest) (*model.News, error)
	SetPhoto(ctx context.Context, authorID, id uuid.UUID, key string) (*model.News, *string, error)
	DeleteNews(ctx context.Context, authorID, id uuid.UUID) error
}

type newsPhotoStorage interface {
	Upload(ctx context.Context, key, contentType string, body io.Reader, size int64) error
	PresignDownload(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type NewsHandler struct {
	service newsService
	storage newsPhotoStorage
	logger  *zap.Logger
}

func NewNewsHandler(service newsService, storage newsPhotoStorage, logger *zap.Logger) *NewsHandler {
	return &NewsHandler{service: service, storage: storage, logger: logger}
}

func (h *NewsHandler) attachURLs(ctx context.Context, n *model.News) {
	n.PhotoURL = presignPhoto(ctx, h.storage, h.logger, n.PhotoS3Key)
	if n.Author != nil {
		n.Author.PhotoURL = presignPhoto(ctx, h.storage, h.logger, n.Author.PhotoS3Key)
	}
}

func (h *NewsHandler) GetNews(c fiber.Ctx) error {
	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}
	limit, offset := q.Normalize()

	news, total, err := h.service.GetNews(c.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to get news", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	for _, n := range news {
		h.attachURLs(c.Context(), n)
	}
	return c.JSON(model.Page[*model.News]{
		Data:   news,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *NewsHandler) GetMyNews(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}
	limit, offset := q.Normalize()

	news, total, err := h.service.GetNewsByAuthorID(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get news", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	for _, n := range news {
		h.attachURLs(c.Context(), n)
	}
	return c.JSON(model.Page[*model.News]{
		Data:   news,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *NewsHandler) GetNewsByID(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid news id")
	}

	news, err := h.service.GetNewsByID(c.Context(), id, cuUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "news not found")
		}
		h.logger.Error("failed to get news by id", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	h.attachURLs(c.Context(), news)
	return c.JSON(news)
}

func (h *NewsHandler) CreateNews(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.CreateNewsRequest
	if err := bindMultipartData(c, &req); err != nil {
		return err
	}

	news, err := h.service.CreateNews(c.Context(), cuUser.ID, req)
	if err != nil {
		h.logger.Error("failed to create news", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}

	news, err = h.applyPhoto(c, cuUser.ID, news)
	if err != nil {
		return err
	}

	h.attachURLs(c.Context(), news)
	return c.Status(fiber.StatusCreated).JSON(news)
}

func (h *NewsHandler) UpdateNews(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid news id")
	}

	var req model.UpdateNewsRequest
	present, err := bindOptionalMultipartData(c, &req)
	if err != nil {
		return err
	}

	var news *model.News
	if present {
		news, err = h.service.UpdateNews(c.Context(), cuUser.ID, id, req)
	} else {
		news, err = h.service.GetNewsByID(c.Context(), id, cuUser.ID)
	}
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "news not found")
		}
		h.logger.Error("failed to update news", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}

	news, err = h.applyPhoto(c, cuUser.ID, news)
	if err != nil {
		return err
	}

	h.attachURLs(c.Context(), news)
	return c.JSON(news)
}

func (h *NewsHandler) applyPhoto(c fiber.Ctx, authorID uuid.UUID, news *model.News) (*model.News, error) {
	key, err := uploadOptionalPhoto(c, h.storage, h.logger, "news/"+news.ID.String()+"/")
	if err != nil {
		return nil, err
	}
	if key == "" {
		return news, nil
	}

	updated, oldKey, err := h.service.SetPhoto(c.Context(), authorID, news.ID, key)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, respondError(c, fiber.StatusNotFound, "news not found")
		}
		h.logger.Error("failed to set news photo", zap.Error(err))
		return nil, respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	deleteReplacedPhoto(c.Context(), h.storage, h.logger, oldKey, key)
	return updated, nil
}

func (h *NewsHandler) DeleteNews(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid news id")
	}

	if err := h.service.DeleteNews(c.Context(), cuUser.ID, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "news not found")
		}
		h.logger.Error("failed to delete news", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
