package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type savedUserRepository interface {
	AddSavedUser(ctx context.Context, userID, savedUserID uuid.UUID) error
	DeleteSavedUser(ctx context.Context, userID, savedUserID uuid.UUID) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetSavedUsers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.User, int, error)
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
	if err := s.repo.AddSavedUser(ctx, userID, savedUserID); err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(ctx, savedUserID)
}

func (s *SavedUserService) DeleteSavedUser(ctx context.Context, userID, savedUserID uuid.UUID) error {
	return s.repo.DeleteSavedUser(ctx, userID, savedUserID)
}
