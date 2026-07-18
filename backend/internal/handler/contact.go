package handler

import (
	"context"
	"errors"

	"faculty/internal/model"
	"faculty/internal/repository"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contactService interface {
	GetContactsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Contact, int, error)
	CreateContact(ctx context.Context, userID uuid.UUID, req model.CreateContactRequest) (*model.Contact, error)
	UpdateContact(ctx context.Context, userID, contactID uuid.UUID, req model.UpdateContactRequest) (*model.Contact, error)
	DeleteContact(ctx context.Context, userID, contactID uuid.UUID) error
}

type contactPhotoStorage interface {
	PresignDownload(ctx context.Context, key string) (string, error)
}

type ContactHandler struct {
	service contactService
	storage contactPhotoStorage
	logger  *zap.Logger
}

func NewContactHandler(service contactService, storage contactPhotoStorage, logger *zap.Logger) *ContactHandler {
	return &ContactHandler{service: service, storage: storage, logger: logger}
}

func (h *ContactHandler) attachPhotoURL(ctx context.Context, u *model.User) {
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

func (h *ContactHandler) GetMyContacts(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError()
	}
	limit, offset := q.Normalize()

	contacts, total, err := h.service.GetContactsByUserID(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		return unexpectedError(h.logger, "failed to get contacts", err)
	}
	for _, contact := range contacts {
		h.attachPhotoURL(c.Context(), contact.User)
	}
	return c.JSON(model.Page[*model.Contact]{Data: contacts, Total: total, Limit: limit, Offset: offset})
}

func (h *ContactHandler) CreateContact(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.CreateContactRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	contact, err := h.service.CreateContact(c.Context(), cuUser.ID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSelfContact):
			return respondValidation(err.Error())
		case errors.Is(err, repository.ErrDuplicate):
			return respondError(fiber.StatusConflict, "contact already exists")
		case errors.Is(err, repository.ErrInvalidRefID):
			return respondError(fiber.StatusNotFound, "user not found")
		}
		return unexpectedError(h.logger, "failed to create contact", err)
	}
	h.attachPhotoURL(c.Context(), contact.User)
	return c.Status(fiber.StatusCreated).JSON(contact)
}

func (h *ContactHandler) UpdateContact(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	contactID, err := uuid.Parse(c.Params("contactId"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid contact id")
	}

	var req model.UpdateContactRequest
	if err := bindJSON(c, &req); err != nil {
		return err
	}

	contact, err := h.service.UpdateContact(c.Context(), cuUser.ID, contactID, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "contact not found")
		}
		return unexpectedError(h.logger, "failed to update contact", err)
	}
	h.attachPhotoURL(c.Context(), contact.User)
	return c.JSON(contact)
}

func (h *ContactHandler) DeleteContact(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	contactID, err := uuid.Parse(c.Params("contactId"))
	if err != nil {
		return respondError(fiber.StatusBadRequest, "invalid contact id")
	}

	if err := h.service.DeleteContact(c.Context(), cuUser.ID, contactID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(fiber.StatusNotFound, "contact not found")
		}
		return unexpectedError(h.logger, "failed to delete contact", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
