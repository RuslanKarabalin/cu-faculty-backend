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
	SearchUsers(ctx context.Context, userID uuid.UUID, search string, limit, offset int) ([]*model.UserSearchResult, int, error)
	UpdateUser(ctx context.Context, params model.UpdateUserParams) error
	UpdateUserPhoto(ctx context.Context, id uuid.UUID, key string) (*string, error)
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

func (s *UserService) SearchUsers(ctx context.Context, userID uuid.UUID, search string, limit, offset int) ([]*model.UserSearchResult, int, error) {
	return s.repo.SearchUsers(ctx, userID, search, limit, offset)
}

func (s *UserService) SetPhoto(ctx context.Context, id uuid.UUID, key string) (*model.User, *string, error) {
	oldKey, err := s.repo.UpdateUserPhoto(ctx, id, key)
	if err != nil {
		return nil, nil, err
	}
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return user, oldKey, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req model.UpdateUserRequest) (*model.User, error) {
	if err := s.repo.UpdateUser(ctx, model.UpdateUserParams{
		ID:         id,
		Bio:        req.Bio,
		Speciality: req.Speciality,
		StatusID:   req.StatusID,
	}); err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(ctx, id)
}
