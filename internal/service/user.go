package service

import (
	"context"
	"errors"
	"faculty/internal/model"
	"faculty/internal/repository"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAlreadyRegistered = errors.New("user already registered")
	ErrInvalidBirthDate  = errors.New("invalid birth_date format, expected YYYY-MM-DD")
)

type userRepository interface {
	CreateUser(ctx context.Context, params model.CreateUserParams) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error)
}

type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, cuUser model.CuUserResp) error {
	birthDate, err := time.Parse("2006-01-02", cuUser.BirthDate)
	if err != nil {
		return ErrInvalidBirthDate
	}

	err = s.repo.CreateUser(ctx, model.CreateUserParams{
		ID:        cuUser.Id,
		FirstName: cuUser.FirstName,
		LastName:  cuUser.LastName,
		BirthDate: birthDate,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return ErrAlreadyRegistered
		}
		return err
	}
	return nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error) {
	return s.repo.GetAllUsers(ctx, limit, offset)
}
