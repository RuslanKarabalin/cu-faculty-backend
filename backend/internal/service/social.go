package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type socialRepository interface {
	CreateSocial(ctx context.Context, params model.CreateSocialParams) (int, error)
	UpdateSocial(ctx context.Context, params model.UpdateSocialParams) error
	DeleteSocial(ctx context.Context, id int, userID uuid.UUID) error
	GetSocialByID(ctx context.Context, id int) (*model.Social, error)
	GetSocialsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Social, error)
}

type SocialService struct {
	repo socialRepository
}

func NewSocialService(repo socialRepository) *SocialService {
	return &SocialService{repo: repo}
}

func (s *SocialService) GetSocialsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Social, error) {
	return s.repo.GetSocialsByUserID(ctx, userID)
}

func (s *SocialService) CreateSocial(ctx context.Context, userID uuid.UUID, req model.SocialRequest) (*model.Social, error) {
	id, err := s.repo.CreateSocial(ctx, model.CreateSocialParams{
		UserId:      userID,
		Social:      req.Social,
		Link:        req.Link,
		IsPreferred: req.IsPreferred,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetSocialByID(ctx, id)
}

func (s *SocialService) UpdateSocial(ctx context.Context, userID uuid.UUID, id int, req model.SocialRequest) (*model.Social, error) {
	err := s.repo.UpdateSocial(ctx, model.UpdateSocialParams{
		ID:          id,
		UserId:      userID,
		Social:      req.Social,
		Link:        req.Link,
		IsPreferred: req.IsPreferred,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetSocialByID(ctx, id)
}

func (s *SocialService) DeleteSocial(ctx context.Context, userID uuid.UUID, id int) error {
	return s.repo.DeleteSocial(ctx, id, userID)
}
