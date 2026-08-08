package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type blockedUserRepository interface {
	AddBlockedUser(ctx context.Context, userID, blockedUserID uuid.UUID) error
	DeleteBlockedUser(ctx context.Context, userID, blockedUserID uuid.UUID) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetBlockedUsers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.User, int, error)
}

type BlockedUserService struct {
	repo blockedUserRepository
}

func NewBlockedUserService(repo blockedUserRepository) *BlockedUserService {
	return &BlockedUserService{repo: repo}
}

func (s *BlockedUserService) GetBlockedUsers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.User, int, error) {
	return s.repo.GetBlockedUsers(ctx, userID, limit, offset)
}

func (s *BlockedUserService) AddBlockedUser(ctx context.Context, userID, blockedUserID uuid.UUID) (*model.User, error) {
	if err := s.repo.AddBlockedUser(ctx, userID, blockedUserID); err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(ctx, blockedUserID)
}

func (s *BlockedUserService) DeleteBlockedUser(ctx context.Context, userID, blockedUserID uuid.UUID) error {
	return s.repo.DeleteBlockedUser(ctx, userID, blockedUserID)
}
