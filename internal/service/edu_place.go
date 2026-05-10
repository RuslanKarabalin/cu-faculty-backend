package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type eduPlaceRepository interface {
	GetEduPlacesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EduPlace, error)
}

type EduPlaceService struct {
	repo eduPlaceRepository
}

func NewEduPlaceService(repo eduPlaceRepository) *EduPlaceService {
	return &EduPlaceService{repo: repo}
}

func (s *EduPlaceService) GetEduPlacesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EduPlace, error) {
	return s.repo.GetEduPlacesByUserID(ctx, userID)
}
