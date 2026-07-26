package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type eventRepository interface {
	CreateEvent(ctx context.Context, params model.CreateEventParams) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, params model.UpdateEventParams) error
	UpdateEventPhoto(ctx context.Context, id, authorID uuid.UUID, key *string) (*string, error)
	DeleteEvent(ctx context.Context, id, authorID uuid.UUID) (*string, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (*model.Event, error)
	GetVisibleEventByID(ctx context.Context, id, viewerID uuid.UUID) (*model.Event, error)
	GetEvents(ctx context.Context, limit, offset int) ([]*model.Event, int, error)
	GetEventsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Event, int, error)
}

type EventService struct {
	repo eventRepository
}

func NewEventService(repo eventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) GetEventByID(ctx context.Context, id, viewerID uuid.UUID) (*model.Event, error) {
	return s.repo.GetVisibleEventByID(ctx, id, viewerID)
}

func (s *EventService) GetEvents(ctx context.Context, limit, offset int) ([]*model.Event, int, error) {
	return s.repo.GetEvents(ctx, limit, offset)
}

func (s *EventService) GetEventsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Event, int, error) {
	return s.repo.GetEventsByAuthorID(ctx, authorID, limit, offset)
}

func (s *EventService) CreateEvent(ctx context.Context, authorID uuid.UUID, req model.CreateEventRequest) (*model.Event, error) {
	id, err := s.repo.CreateEvent(ctx, model.CreateEventParams{
		AuthorID:         authorID,
		Title:            req.Title,
		Content:          req.Content,
		Place:            req.Place,
		Category:         req.Category,
		StartsAt:         req.StartsAt,
		RegistrationLink: req.RegistrationLink,
		IsDraft:          req.IsDraft,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetEventByID(ctx, id)
}

func (s *EventService) UpdateEvent(ctx context.Context, authorID, id uuid.UUID, req model.UpdateEventRequest) (*model.Event, error) {
	current, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		return nil, err
	}

	params := model.UpdateEventParams{
		ID:               id,
		AuthorID:         authorID,
		Title:            current.Title,
		Content:          current.Content,
		Place:            current.Place,
		Category:         current.Category,
		StartsAt:         current.StartsAt,
		RegistrationLink: current.RegistrationLink,
		IsDraft:          current.IsDraft,
	}
	if req.Title != nil {
		params.Title = *req.Title
	}
	if req.Content != nil {
		params.Content = *req.Content
	}
	if req.Place != nil {
		params.Place = *req.Place
	}
	if req.Category != nil {
		params.Category = *req.Category
	}
	if req.StartsAt != nil {
		params.StartsAt = *req.StartsAt
	}
	if req.RegistrationLink != nil {
		params.RegistrationLink = req.RegistrationLink
	}
	if req.IsDraft != nil {
		params.IsDraft = *req.IsDraft
	}

	if err := s.repo.UpdateEvent(ctx, params); err != nil {
		return nil, err
	}
	return s.repo.GetEventByID(ctx, id)
}

func (s *EventService) SetPhoto(ctx context.Context, authorID, id uuid.UUID, key string) (*model.Event, *string, error) {
	oldKey, err := s.repo.UpdateEventPhoto(ctx, id, authorID, &key)
	if err != nil {
		return nil, nil, err
	}
	event, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return event, oldKey, nil
}

// DeletePhoto clears the event photo and returns the key that was stored, if
// any, so the caller can remove the object from storage.
func (s *EventService) DeletePhoto(ctx context.Context, authorID, id uuid.UUID) (*string, error) {
	return s.repo.UpdateEventPhoto(ctx, id, authorID, nil)
}

func (s *EventService) DeleteEvent(ctx context.Context, authorID, id uuid.UUID) (*string, error) {
	return s.repo.DeleteEvent(ctx, id, authorID)
}
