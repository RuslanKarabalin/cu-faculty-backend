package service

import (
	"context"
	"errors"
	"time"

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
	SearchUsers(ctx context.Context, viewerID uuid.UUID, params model.SearchUsersParams) ([]*model.UserSearchResult, int, error)
	UpdateUser(ctx context.Context, params model.UpdateUserParams) error
	UpdateUserPhoto(ctx context.Context, id uuid.UUID, key *string) (*string, error)
	GetProfileCompletenessData(ctx context.Context, id uuid.UUID) (*model.ProfileCompletenessData, error)
	SetUserCompletedAt(ctx context.Context, id uuid.UUID, completedAt *time.Time) error
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

func (s *UserService) SearchUsers(ctx context.Context, viewerID uuid.UUID, params model.SearchUsersParams) ([]*model.UserSearchResult, int, error) {
	return s.repo.SearchUsers(ctx, viewerID, params)
}

func (s *UserService) SetPhoto(ctx context.Context, id uuid.UUID, key string) (*model.User, *string, error) {
	oldKey, err := s.repo.UpdateUserPhoto(ctx, id, &key)
	if err != nil {
		return nil, nil, err
	}
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return user, oldKey, nil
}

func (s *UserService) DeletePhoto(ctx context.Context, id uuid.UUID) (*string, error) {
	return s.repo.UpdateUserPhoto(ctx, id, nil)
}

func (s *UserService) GetProfileCompleteness(ctx context.Context, id uuid.UUID) (*model.ProfileCompleteness, error) {
	data, err := s.repo.GetProfileCompletenessData(ctx, id)
	if err != nil {
		return nil, err
	}

	percent := data.Percent()
	switch {
	case percent == model.CompletenessMax && data.CompletedAt == nil:
		now := time.Now()
		if err := s.repo.SetUserCompletedAt(ctx, id, &now); err != nil {
			return nil, err
		}
	case percent != model.CompletenessMax && data.CompletedAt != nil:
		if err := s.repo.SetUserCompletedAt(ctx, id, nil); err != nil {
			return nil, err
		}
	}
	return &model.ProfileCompleteness{Percent: percent}, nil
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
