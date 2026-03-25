package service

import (
	"context"
	"faculty/internal/model"
	"faculty/internal/repository"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user model.CuUserResp) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.GetAllUsers(ctx)
}
