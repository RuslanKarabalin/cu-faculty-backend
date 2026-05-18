package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type workPlaceRepository interface {
	CreateWorkPlace(ctx context.Context, params model.CreateWorkPlaceParams) (int, error)
	UpdateWorkPlace(ctx context.Context, params model.UpdateWorkPlaceParams) error
	DeleteWorkPlace(ctx context.Context, id int, userID uuid.UUID) error
	GetWorkPlaceByID(ctx context.Context, id int) (*model.WorkPlace, error)
	GetWorkPlacesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.WorkPlace, error)
}

type WorkPlaceService struct {
	repo workPlaceRepository
}

func NewWorkPlaceService(repo workPlaceRepository) *WorkPlaceService {
	return &WorkPlaceService{repo: repo}
}

func (s *WorkPlaceService) GetWorkPlacesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.WorkPlace, error) {
	return s.repo.GetWorkPlacesByUserID(ctx, userID)
}

func (s *WorkPlaceService) CreateWorkPlace(ctx context.Context, userID uuid.UUID, req model.WorkPlaceRequest) (*model.WorkPlace, error) {
	id, err := s.repo.CreateWorkPlace(ctx, model.CreateWorkPlaceParams{
		UserId:       userID,
		CompanyName:  req.CompanyName,
		Grade:        req.Grade,
		Position:     req.Position,
		StartYear:    req.StartYear,
		EndYear:      req.EndYear,
		IsWorkingNow: req.IsWorkingNow,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetWorkPlaceByID(ctx, id)
}

func (s *WorkPlaceService) UpdateWorkPlace(ctx context.Context, userID uuid.UUID, id int, req model.WorkPlaceRequest) (*model.WorkPlace, error) {
	err := s.repo.UpdateWorkPlace(ctx, model.UpdateWorkPlaceParams{
		ID:           id,
		UserId:       userID,
		CompanyName:  req.CompanyName,
		Grade:        req.Grade,
		Position:     req.Position,
		StartYear:    req.StartYear,
		EndYear:      req.EndYear,
		IsWorkingNow: req.IsWorkingNow,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetWorkPlaceByID(ctx, id)
}

func (s *WorkPlaceService) DeleteWorkPlace(ctx context.Context, userID uuid.UUID, id int) error {
	return s.repo.DeleteWorkPlace(ctx, id, userID)
}
