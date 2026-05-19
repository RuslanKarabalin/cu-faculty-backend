package service

import (
	"context"

	"faculty/internal/model"
)

type referenceRepository interface {
	GetStatuses(ctx context.Context) ([]*model.Status, error)
	GetKeySkills(ctx context.Context) ([]*model.Skill, error)
	GetSoftSkills(ctx context.Context) ([]*model.Skill, error)
	GetCompanies(ctx context.Context) ([]*model.Company, error)
	GetWorkPositions(ctx context.Context) ([]*model.WorkPosition, error)
	GetUniversities(ctx context.Context) ([]*model.University, error)
	GetFaqs(ctx context.Context) ([]*model.Faq, error)
	GetEnumValues(ctx context.Context, typeName string) ([]string, error)
}

type ReferenceService struct {
	repo referenceRepository
}

func NewReferenceService(repo referenceRepository) *ReferenceService {
	return &ReferenceService{repo: repo}
}

func (s *ReferenceService) GetStatuses(ctx context.Context) ([]*model.Status, error) {
	return s.repo.GetStatuses(ctx)
}

func (s *ReferenceService) GetKeySkills(ctx context.Context) ([]*model.Skill, error) {
	return s.repo.GetKeySkills(ctx)
}

func (s *ReferenceService) GetSoftSkills(ctx context.Context) ([]*model.Skill, error) {
	return s.repo.GetSoftSkills(ctx)
}

func (s *ReferenceService) GetCompanies(ctx context.Context) ([]*model.Company, error) {
	return s.repo.GetCompanies(ctx)
}

func (s *ReferenceService) GetWorkPositions(ctx context.Context) ([]*model.WorkPosition, error) {
	return s.repo.GetWorkPositions(ctx)
}

func (s *ReferenceService) GetUniversities(ctx context.Context) ([]*model.University, error) {
	return s.repo.GetUniversities(ctx)
}

func (s *ReferenceService) GetFaqs(ctx context.Context) ([]*model.Faq, error) {
	return s.repo.GetFaqs(ctx)
}

func (s *ReferenceService) GetSocialNetworks(ctx context.Context) ([]string, error) {
	return s.repo.GetEnumValues(ctx, "social_network")
}

func (s *ReferenceService) GetEduGrades(ctx context.Context) ([]string, error) {
	return s.repo.GetEnumValues(ctx, "edu_grade")
}

func (s *ReferenceService) GetWorkGrades(ctx context.Context) ([]string, error) {
	return s.repo.GetEnumValues(ctx, "work_grade")
}

func (s *ReferenceService) GetEventCategories(ctx context.Context) ([]string, error) {
	return s.repo.GetEnumValues(ctx, "event_category")
}
