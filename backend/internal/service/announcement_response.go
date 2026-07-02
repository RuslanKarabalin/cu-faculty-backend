package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type announcementResponseRepository interface {
	AddAnnouncementResponse(ctx context.Context, userID, announcementID uuid.UUID) error
	DeleteAnnouncementResponse(ctx context.Context, userID, announcementID uuid.UUID) error
	GetAnnouncementResponders(ctx context.Context, announcementID uuid.UUID) ([]*model.User, error)
	GetAnnouncementsRespondedByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Announcement, int, error)
	GetAnnouncementByID(ctx context.Context, id uuid.UUID) (*model.Announcement, error)
}

type AnnouncementResponseService struct {
	repo announcementResponseRepository
}

func NewAnnouncementResponseService(repo announcementResponseRepository) *AnnouncementResponseService {
	return &AnnouncementResponseService{repo: repo}
}

func (s *AnnouncementResponseService) Respond(ctx context.Context, userID, announcementID uuid.UUID) (*model.Announcement, error) {
	if err := s.repo.AddAnnouncementResponse(ctx, userID, announcementID); err != nil {
		return nil, err
	}
	return s.repo.GetAnnouncementByID(ctx, announcementID)
}

func (s *AnnouncementResponseService) DeleteResponse(ctx context.Context, userID, announcementID uuid.UUID) error {
	return s.repo.DeleteAnnouncementResponse(ctx, userID, announcementID)
}

func (s *AnnouncementResponseService) GetResponders(ctx context.Context, announcementID uuid.UUID) ([]*model.User, error) {
	return s.repo.GetAnnouncementResponders(ctx, announcementID)
}

func (s *AnnouncementResponseService) GetMyResponses(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Announcement, int, error) {
	return s.repo.GetAnnouncementsRespondedByUser(ctx, userID, limit, offset)
}
