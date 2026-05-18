package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type userKeySkillRepository interface {
	AddUserKeySkill(ctx context.Context, userID uuid.UUID, skillID int) error
	DeleteUserKeySkill(ctx context.Context, userID uuid.UUID, skillID int) error
	GetKeySkillByID(ctx context.Context, id int) (*model.Skill, error)
	GetUserKeySkills(ctx context.Context, userID uuid.UUID) ([]*model.Skill, error)
}

type UserKeySkillService struct {
	repo userKeySkillRepository
}

func NewUserKeySkillService(repo userKeySkillRepository) *UserKeySkillService {
	return &UserKeySkillService{repo: repo}
}

func (s *UserKeySkillService) GetUserKeySkills(ctx context.Context, userID uuid.UUID) ([]*model.Skill, error) {
	return s.repo.GetUserKeySkills(ctx, userID)
}

func (s *UserKeySkillService) AddUserKeySkill(ctx context.Context, userID uuid.UUID, skillID int) (*model.Skill, error) {
	if err := s.repo.AddUserKeySkill(ctx, userID, skillID); err != nil {
		return nil, err
	}
	return s.repo.GetKeySkillByID(ctx, skillID)
}

func (s *UserKeySkillService) DeleteUserKeySkill(ctx context.Context, userID uuid.UUID, skillID int) error {
	return s.repo.DeleteUserKeySkill(ctx, userID, skillID)
}
