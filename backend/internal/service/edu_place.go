package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type eduPlaceRepository interface {
	CreateEduPlace(ctx context.Context, params model.CreateEduPlaceParams) (int, error)
	UpdateEduPlace(ctx context.Context, params model.UpdateEduPlaceParams) error
	DeleteEduPlace(ctx context.Context, id int, userID uuid.UUID) error
	GetEduPlaceByID(ctx context.Context, id int) (*model.EduPlace, error)
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

func (s *EduPlaceService) CreateEduPlace(ctx context.Context, userID uuid.UUID, req model.EduPlaceRequest) (*model.EduPlace, error) {
	id, err := s.repo.CreateEduPlace(ctx, model.CreateEduPlaceParams{
		UserId:         userID,
		UniversityId:   req.UniversityId,
		Grade:          req.Grade,
		Level:          req.Level,
		Specialization: req.Specialization,
		StartYear:      req.StartYear,
		EndYear:        req.EndYear,
		IsStudyingNow:  req.IsStudyingNow,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetEduPlaceByID(ctx, id)
}

func (s *EduPlaceService) UpdateEduPlace(ctx context.Context, userID uuid.UUID, id int, req model.EduPlaceRequest) (*model.EduPlace, error) {
	err := s.repo.UpdateEduPlace(ctx, model.UpdateEduPlaceParams{
		ID:             id,
		UserId:         userID,
		UniversityId:   req.UniversityId,
		Grade:          req.Grade,
		Level:          req.Level,
		Specialization: req.Specialization,
		StartYear:      req.StartYear,
		EndYear:        req.EndYear,
		IsStudyingNow:  req.IsStudyingNow,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetEduPlaceByID(ctx, id)
}

func (s *EduPlaceService) DeleteEduPlace(ctx context.Context, userID uuid.UUID, id int) error {
	return s.repo.DeleteEduPlace(ctx, id, userID)
}
