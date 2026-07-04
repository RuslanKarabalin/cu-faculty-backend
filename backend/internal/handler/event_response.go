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

type eventResponseService interface {
	Respond(ctx context.Context, userID, eventID uuid.UUID) (*model.Event, error)
	DeleteResponse(ctx context.Context, userID, eventID uuid.UUID) error
	GetResponders(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]*model.User, int, error)
	GetMyResponses(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Event, int, error)
}

type EventResponseHandler struct {
	service eventResponseService
	storage photoPresigner
	logger  *zap.Logger
}

func NewEventResponseHandler(service eventResponseService, storage photoPresigner, logger *zap.Logger) *EventResponseHandler {
	return &EventResponseHandler{service: service, storage: storage, logger: logger}
}

func (h *EventResponseHandler) attachURLs(ctx context.Context, e *model.Event) {
	e.PhotoURL = presignPhoto(ctx, h.storage, h.logger, e.PhotoS3Key)
	if e.Author != nil {
		e.Author.PhotoURL = presignPhoto(ctx, h.storage, h.logger, e.Author.PhotoS3Key)
	}
}

func (h *EventResponseHandler) RespondToEvent(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid event id")
	}

	event, err := h.service.Respond(c.Context(), cuUser.ID, eventID)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidRefID) || errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "event not found")
		}
		return unexpectedError(c, h.logger, "failed to respond to event", err)
	}
	h.attachURLs(c.Context(), event)
	return c.Status(fiber.StatusCreated).JSON(event)
}

func (h *EventResponseHandler) DeleteMyResponse(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid event id")
	}

	if err := h.service.DeleteResponse(c.Context(), cuUser.ID, eventID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "response not found")
		}
		return unexpectedError(c, h.logger, "failed to delete event response", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *EventResponseHandler) GetResponders(c fiber.Ctx) error {
	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid event id")
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	users, total, err := h.service.GetResponders(c.Context(), eventID, limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get event responders", err)
	}
	for _, u := range users {
		u.PhotoURL = presignPhoto(c.Context(), h.storage, h.logger, u.PhotoS3Key)
	}
	return c.JSON(model.Page[*model.User]{Data: users, Total: total, Limit: limit, Offset: offset})
}

func (h *EventResponseHandler) GetMyResponses(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	events, total, err := h.service.GetMyResponses(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get event responses", err)
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
