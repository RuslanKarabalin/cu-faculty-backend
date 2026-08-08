package service

import (
	"context"

	"faculty/internal/model"
	"faculty/internal/repository"

	"github.com/google/uuid"
)

type savedUserRepository interface {
	AddSavedUser(ctx context.Context, userID, savedUserID uuid.UUID) error
	DeleteSavedUser(ctx context.Context, userID, savedUserID uuid.UUID) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetSavedUsers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.User, int, error)
	HasBlockBetween(ctx context.Context, a, b uuid.UUID) (bool, error)
}

type SavedUserService struct {
	repo savedUserRepository
}

func NewSavedUserService(repo savedUserRepository) *SavedUserService {
	return &SavedUserService{repo: repo}
}

func (s *SavedUserService) GetSavedUsers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.User, int, error) {
	return s.repo.GetSavedUsers(ctx, userID, limit, offset)
}

func (s *SavedUserService) AddSavedUser(ctx context.Context, userID, savedUserID uuid.UUID) (*model.User, error) {
	target, err := s.repo.GetUserByID(ctx, savedUserID)
	if err != nil {
		return nil, err
	}
	if target.IsDeleted() {
		return nil, repository.ErrNotFound
	}
	blocked, err := s.repo.HasBlockBetween(ctx, userID, savedUserID)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, ErrBlocked
	}
	if err := s.repo.AddSavedUser(ctx, userID, savedUserID); err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(ctx, savedUserID)
}

func (s *SavedUserService) DeleteSavedUser(ctx context.Context, userID, savedUserID uuid.UUID) error {
	return s.repo.DeleteSavedUser(ctx, userID, savedUserID)
}
