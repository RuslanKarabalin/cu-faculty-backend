package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type eventResponseRepository interface {
	AddEventResponse(ctx context.Context, userID, eventID uuid.UUID) error
	DeleteEventResponse(ctx context.Context, userID, eventID uuid.UUID) error
	GetEventResponders(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]*model.User, int, error)
	GetEventsRespondedByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Event, int, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (*model.Event, error)
}

type EventResponseService struct {
	repo eventResponseRepository
}

func NewEventResponseService(repo eventResponseRepository) *EventResponseService {
	return &EventResponseService{repo: repo}
}

func (s *EventResponseService) Respond(ctx context.Context, userID, eventID uuid.UUID) (*model.Event, error) {
	if err := s.repo.AddEventResponse(ctx, userID, eventID); err != nil {
		return nil, err
	}
	return s.repo.GetEventByID(ctx, eventID)
}

func (s *EventResponseService) DeleteResponse(ctx context.Context, userID, eventID uuid.UUID) error {
	return s.repo.DeleteEventResponse(ctx, userID, eventID)
}

func (s *EventResponseService) GetResponders(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]*model.User, int, error) {
	return s.repo.GetEventResponders(ctx, eventID, limit, offset)
}

func (s *EventResponseService) GetMyResponses(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Event, int, error) {
	return s.repo.GetEventsRespondedByUser(ctx, userID, limit, offset)
}
