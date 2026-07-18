package handler

import (
	"context"

	"faculty/internal/cuclient"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type eventSyncService interface {
	Sync(ctx context.Context, cookie string) (int, error)
}

type EventSyncHandler struct {
	service eventSyncService
	logger  *zap.Logger
}

func NewEventSyncHandler(service eventSyncService, logger *zap.Logger) *EventSyncHandler {
	return &EventSyncHandler{service: service, logger: logger}
}

func (h *EventSyncHandler) TriggerSync(c fiber.Ctx) error {
	processed, err := h.service.Sync(c.Context(), c.Cookies(cuclient.CookieName))
	if err != nil {
		return unexpectedError(c, h.logger, "failed to sync events", err)
	}
	return c.JSON(fiber.Map{"processed": processed})
}
