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

type eventService interface {
	GetEventByID(ctx context.Context, id, viewerID uuid.UUID) (*model.Event, error)
	GetEvents(ctx context.Context, limit, offset int) ([]*model.Event, int, error)
	GetEventsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Event, int, error)
	CreateEvent(ctx context.Context, authorID uuid.UUID, req model.CreateEventRequest) (*model.Event, error)
	UpdateEvent(ctx context.Context, authorID, id uuid.UUID, req model.UpdateEventRequest) (*model.Event, error)
	SetPhoto(ctx context.Context, authorID, id uuid.UUID, key string) (*model.Event, *string, error)
	DeleteEvent(ctx context.Context, authorID, id uuid.UUID) (*string, error)
}

type eventPhotoStorage interface {
	Upload(ctx context.Context, key, contentType string, body io.Reader, size int64) error
	PresignDownload(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type EventHandler struct {
	service eventService
	storage eventPhotoStorage
	logger  *zap.Logger
}

func NewEventHandler(service eventService, storage eventPhotoStorage, logger *zap.Logger) *EventHandler {
	return &EventHandler{service: service, storage: storage, logger: logger}
}

func (h *EventHandler) attachURLs(ctx context.Context, e *model.Event) {
	e.PhotoURL = presignPhoto(ctx, h.storage, h.logger, e.PhotoS3Key)
	if e.Author != nil {
		e.Author.PhotoURL = presignPhoto(ctx, h.storage, h.logger, e.Author.PhotoS3Key)
	}
}

func (h *EventHandler) GetEvents(c fiber.Ctx) error {
	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}
	limit, offset := q.Normalize()

	events, total, err := h.service.GetEvents(c.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to get events", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	for _, e := range events {
		h.attachURLs(c.Context(), e)
	}
	return c.JSON(model.Page[*model.Event]{
		Data:   events,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *EventHandler) GetMyEvents(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}
	limit, offset := q.Normalize()

	events, total, err := h.service.GetEventsByAuthorID(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get events", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	for _, e := range events {
		h.attachURLs(c.Context(), e)
	}
	return c.JSON(model.Page[*model.Event]{
		Data:   events,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *EventHandler) GetEventByID(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid event id")
	}

	event, err := h.service.GetEventByID(c.Context(), id, cuUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "event not found")
		}
		h.logger.Error("failed to get event by id", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	h.attachURLs(c.Context(), event)
	return c.JSON(event)
}

func (h *EventHandler) CreateEvent(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.CreateEventRequest
	if err := bindMultipartData(c, &req); err != nil {
		return err
	}

	event, err := h.service.CreateEvent(c.Context(), cuUser.ID, req)
	if err != nil {
		h.logger.Error("failed to create event", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}

	event, err = h.applyPhoto(c, cuUser.ID, event)
	if err != nil {
		return err
	}

	h.attachURLs(c.Context(), event)
	return c.Status(fiber.StatusCreated).JSON(event)
}

func (h *EventHandler) UpdateEvent(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid event id")
	}

	var req model.UpdateEventRequest
	present, err := bindOptionalMultipartData(c, &req)
	if err != nil {
		return err
	}

	var event *model.Event
	if present {
		event, err = h.service.UpdateEvent(c.Context(), cuUser.ID, id, req)
	} else {
		event, err = h.service.GetEventByID(c.Context(), id, cuUser.ID)
	}
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "event not found")
		}
		h.logger.Error("failed to update event", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}

	event, err = h.applyPhoto(c, cuUser.ID, event)
	if err != nil {
		return err
	}

	h.attachURLs(c.Context(), event)
	return c.JSON(event)
}

func (h *EventHandler) applyPhoto(c fiber.Ctx, authorID uuid.UUID, event *model.Event) (*model.Event, error) {
	key, err := uploadOptionalPhoto(c, h.storage, h.logger, "events/"+event.ID.String()+"/")
	if err != nil {
		return nil, err
	}
	if key == "" {
		return event, nil
	}

	updated, oldKey, err := h.service.SetPhoto(c.Context(), authorID, event.ID, key)
	if err != nil {
		deletePhoto(c.Context(), h.storage, h.logger, key)
		if errors.Is(err, repository.ErrNotFound) {
			return nil, respondError(c, fiber.StatusNotFound, "event not found")
		}
		h.logger.Error("failed to set event photo", zap.Error(err))
		return nil, respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	deleteReplacedPhoto(c.Context(), h.storage, h.logger, oldKey, key)
	return updated, nil
}

func (h *EventHandler) DeleteEvent(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid event id")
	}

	photoKey, err := h.service.DeleteEvent(c.Context(), cuUser.ID, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "event not found")
		}
		h.logger.Error("failed to delete event", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	if photoKey != nil {
		deletePhoto(c.Context(), h.storage, h.logger, *photoKey)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
