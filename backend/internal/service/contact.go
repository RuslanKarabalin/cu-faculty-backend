package service

import (
	"context"
	"errors"

	"faculty/internal/model"
	"faculty/internal/repository"

	"github.com/google/uuid"
)

var ErrSelfContact = errors.New("cannot add yourself as a contact")

type contactRepository interface {
	CreateContact(ctx context.Context, params model.CreateContactParams) error
	UpdateContact(ctx context.Context, params model.UpdateContactParams) error
	DeleteContact(ctx context.Context, userID, contactID uuid.UUID) error
	GetContact(ctx context.Context, userID, contactID uuid.UUID) (*model.Contact, error)
	GetContactsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Contact, int, error)
	HasBlockBetween(ctx context.Context, a, b uuid.UUID) (bool, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type ContactService struct {
	repo contactRepository
}

func NewContactService(repo contactRepository) *ContactService {
	return &ContactService{repo: repo}
}

func (s *ContactService) GetContactsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Contact, int, error) {
	return s.repo.GetContactsByUserID(ctx, userID, limit, offset)
}

func (s *ContactService) CreateContact(ctx context.Context, userID uuid.UUID, req model.CreateContactRequest) (*model.Contact, error) {
	if req.ContactID == userID {
		return nil, ErrSelfContact
	}
	target, err := s.repo.GetUserByID(ctx, req.ContactID)
	if err != nil {
		return nil, err
	}
	if target.IsDeleted() {
		return nil, repository.ErrNotFound
	}
	blocked, err := s.repo.HasBlockBetween(ctx, userID, req.ContactID)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, ErrBlocked
	}
	if err := s.repo.CreateContact(ctx, model.CreateContactParams{
		UserID:    userID,
		ContactID: req.ContactID,
		Note:      req.Note,
	}); err != nil {
		return nil, err
	}
	return s.repo.GetContact(ctx, userID, req.ContactID)
}

func (s *ContactService) UpdateContact(ctx context.Context, userID, contactID uuid.UUID, req model.UpdateContactRequest) (*model.Contact, error) {
	if err := s.repo.UpdateContact(ctx, model.UpdateContactParams{
		UserID:    userID,
		ContactID: contactID,
		Note:      req.Note,
	}); err != nil {
		return nil, err
	}
	return s.repo.GetContact(ctx, userID, contactID)
}

func (s *ContactService) DeleteContact(ctx context.Context, userID, contactID uuid.UUID) error {
	return s.repo.DeleteContact(ctx, userID, contactID)
}
