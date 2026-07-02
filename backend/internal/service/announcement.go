package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type announcementRepository interface {
	CreateAnnouncement(ctx context.Context, params model.CreateAnnouncementParams) (uuid.UUID, error)
	UpdateAnnouncement(ctx context.Context, params model.UpdateAnnouncementParams) error
	DeleteAnnouncement(ctx context.Context, id, authorID uuid.UUID) error
	GetAnnouncementByID(ctx context.Context, id uuid.UUID) (*model.Announcement, error)
	GetAnnouncements(ctx context.Context, limit, offset int) ([]*model.Announcement, int, error)
	GetAnnouncementsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Announcement, int, error)
}

type AnnouncementService struct {
	repo announcementRepository
}

func NewAnnouncementService(repo announcementRepository) *AnnouncementService {
	return &AnnouncementService{repo: repo}
}

func (s *AnnouncementService) GetAnnouncementByID(ctx context.Context, id uuid.UUID) (*model.Announcement, error) {
	return s.repo.GetAnnouncementByID(ctx, id)
}

func (s *AnnouncementService) GetAnnouncements(ctx context.Context, limit, offset int) ([]*model.Announcement, int, error) {
	return s.repo.GetAnnouncements(ctx, limit, offset)
}

func (s *AnnouncementService) GetAnnouncementsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Announcement, int, error) {
	return s.repo.GetAnnouncementsByAuthorID(ctx, authorID, limit, offset)
}

func (s *AnnouncementService) CreateAnnouncement(ctx context.Context, authorID uuid.UUID, req model.CreateAnnouncementRequest) (*model.Announcement, error) {
	id, err := s.repo.CreateAnnouncement(ctx, model.CreateAnnouncementParams{
		AuthorID: authorID,
		Title:    req.Title,
		Content:  req.Content,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetAnnouncementByID(ctx, id)
}

func (s *AnnouncementService) UpdateAnnouncement(ctx context.Context, authorID, id uuid.UUID, req model.UpdateAnnouncementRequest) (*model.Announcement, error) {
	err := s.repo.UpdateAnnouncement(ctx, model.UpdateAnnouncementParams{
		ID:         id,
		AuthorID:   authorID,
		Title:      req.Title,
		Content:    req.Content,
		IsArchived: req.IsArchived,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetAnnouncementByID(ctx, id)
}

func (s *AnnouncementService) DeleteAnnouncement(ctx context.Context, authorID, id uuid.UUID) error {
	return s.repo.DeleteAnnouncement(ctx, id, authorID)
}
