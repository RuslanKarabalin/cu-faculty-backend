package service

import (
	"context"
	"errors"

	"faculty/internal/model"

	"github.com/google/uuid"
)

var ErrInvalidBirthDate = errors.New("invalid birth_date format, expected YYYY-MM-DD")

type userRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error)
}

type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error) {
	return s.repo.GetAllUsers(ctx, limit, offset)
}
