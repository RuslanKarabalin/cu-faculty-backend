package service

import (
	"context"

	"faculty/internal/model"
)

type referenceRepository interface {
	GetStatuses(ctx context.Context, limit, offset int) ([]*model.Status, int, error)
	GetKeySkills(ctx context.Context, limit, offset int) ([]*model.Skill, int, error)
	GetSoftSkills(ctx context.Context, limit, offset int) ([]*model.Skill, int, error)
	GetCompanies(ctx context.Context, limit, offset int) ([]*model.Company, int, error)
	GetWorkPositions(ctx context.Context, limit, offset int) ([]*model.WorkPosition, int, error)
	GetUniversities(ctx context.Context, limit, offset int) ([]*model.University, int, error)
	GetFaqs(ctx context.Context, limit, offset int) ([]*model.Faq, int, error)
	GetEnumValues(ctx context.Context, typeName string, limit, offset int) ([]string, int, error)
}

type ReferenceService struct {
	repo referenceRepository
}

func NewReferenceService(repo referenceRepository) *ReferenceService {
	return &ReferenceService{repo: repo}
}

func (s *ReferenceService) GetStatuses(ctx context.Context, limit, offset int) ([]*model.Status, int, error) {
	return s.repo.GetStatuses(ctx, limit, offset)
}

func (s *ReferenceService) GetKeySkills(ctx context.Context, limit, offset int) ([]*model.Skill, int, error) {
	return s.repo.GetKeySkills(ctx, limit, offset)
}

func (s *ReferenceService) GetSoftSkills(ctx context.Context, limit, offset int) ([]*model.Skill, int, error) {
	return s.repo.GetSoftSkills(ctx, limit, offset)
}

func (s *ReferenceService) GetCompanies(ctx context.Context, limit, offset int) ([]*model.Company, int, error) {
	return s.repo.GetCompanies(ctx, limit, offset)
}

func (s *ReferenceService) GetWorkPositions(ctx context.Context, limit, offset int) ([]*model.WorkPosition, int, error) {
	return s.repo.GetWorkPositions(ctx, limit, offset)
}

func (s *ReferenceService) GetUniversities(ctx context.Context, limit, offset int) ([]*model.University, int, error) {
	return s.repo.GetUniversities(ctx, limit, offset)
}

func (s *ReferenceService) GetFaqs(ctx context.Context, limit, offset int) ([]*model.Faq, int, error) {
	return s.repo.GetFaqs(ctx, limit, offset)
}

func (s *ReferenceService) GetSocialNetworks(ctx context.Context, limit, offset int) ([]string, int, error) {
	return s.repo.GetEnumValues(ctx, "social_network", limit, offset)
}

func (s *ReferenceService) GetEduGrades(ctx context.Context, limit, offset int) ([]string, int, error) {
	return s.repo.GetEnumValues(ctx, "edu_grade", limit, offset)
}

func (s *ReferenceService) GetWorkGrades(ctx context.Context, limit, offset int) ([]string, int, error) {
	return s.repo.GetEnumValues(ctx, "work_grade", limit, offset)
}

func (s *ReferenceService) GetEventCategories(ctx context.Context, limit, offset int) ([]string, int, error) {
	return s.repo.GetEnumValues(ctx, "event_category", limit, offset)
}
