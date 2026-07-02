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
	GetContactsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Contact, error)
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

	contacts, err := h.service.GetContactsByUserID(c.Context(), cuUser.ID)
	if err != nil {
		h.logger.Error("failed to get contacts", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	for _, contact := range contacts {
		h.attachPhotoURL(c.Context(), contact.User)
	}
	return c.JSON(contacts)
}

func (h *ContactHandler) CreateContact(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.CreateContactRequest
	if err := c.Bind().JSON(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	contact, err := h.service.CreateContact(c.Context(), cuUser.ID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSelfContact):
			return respondError(c, fiber.StatusBadRequest, err.Error())
		case errors.Is(err, repository.ErrDuplicate):
			return respondError(c, fiber.StatusConflict, "contact already exists")
		case errors.Is(err, repository.ErrInvalidRefID):
			return respondError(c, fiber.StatusNotFound, "user not found")
		}
		h.logger.Error("failed to create contact", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
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
		return respondError(c, fiber.StatusBadRequest, "invalid contact id")
	}

	var req model.UpdateContactRequest
	if err := c.Bind().JSON(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	contact, err := h.service.UpdateContact(c.Context(), cuUser.ID, contactID, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "contact not found")
		}
		h.logger.Error("failed to update contact", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
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
		return respondError(c, fiber.StatusBadRequest, "invalid contact id")
	}

	if err := h.service.DeleteContact(c.Context(), cuUser.ID, contactID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "contact not found")
		}
		h.logger.Error("failed to delete contact", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
