package service

import (
	"context"
	"errors"

	"faculty/internal/model"

	"github.com/google/uuid"
)

var (
	ErrInvalidBirthDate    = errors.New("invalid birth_date format, expected YYYY-MM-DD")
	ErrInvalidUpstreamData = errors.New("invalid data from upstream")
)

type userRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error)
	UpdateUser(ctx context.Context, params model.UpdateUserParams) error
	UpdateUserPhoto(ctx context.Context, id uuid.UUID, key string) error
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

func (s *UserService) SetPhoto(ctx context.Context, id uuid.UUID, key string) (*model.User, error) {
	if err := s.repo.UpdateUserPhoto(ctx, id, key); err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req model.UpdateUserRequest) (*model.User, error) {
	if err := s.repo.UpdateUser(ctx, model.UpdateUserParams{
		ID:         id,
		PhotoS3Key: req.PhotoS3Key,
		Bio:        req.Bio,
		Speciality: req.Speciality,
		StatusID:   req.StatusID,
	}); err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(ctx, id)
}
