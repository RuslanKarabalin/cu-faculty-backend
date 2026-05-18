package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type userSoftSkillRepository interface {
	AddUserSoftSkill(ctx context.Context, userID uuid.UUID, skillID int) error
	DeleteUserSoftSkill(ctx context.Context, userID uuid.UUID, skillID int) error
	GetSoftSkillByID(ctx context.Context, id int) (*model.Skill, error)
	GetUserSoftSkills(ctx context.Context, userID uuid.UUID) ([]*model.Skill, error)
}

type UserSoftSkillService struct {
	repo userSoftSkillRepository
}

func NewUserSoftSkillService(repo userSoftSkillRepository) *UserSoftSkillService {
	return &UserSoftSkillService{repo: repo}
}

func (s *UserSoftSkillService) GetUserSoftSkills(ctx context.Context, userID uuid.UUID) ([]*model.Skill, error) {
	return s.repo.GetUserSoftSkills(ctx, userID)
}

func (s *UserSoftSkillService) AddUserSoftSkill(ctx context.Context, userID uuid.UUID, skillID int) (*model.Skill, error) {
	if err := s.repo.AddUserSoftSkill(ctx, userID, skillID); err != nil {
		return nil, err
	}
	return s.repo.GetSoftSkillByID(ctx, skillID)
}

func (s *UserSoftSkillService) DeleteUserSoftSkill(ctx context.Context, userID uuid.UUID, skillID int) error {
	return s.repo.DeleteUserSoftSkill(ctx, userID, skillID)
}
