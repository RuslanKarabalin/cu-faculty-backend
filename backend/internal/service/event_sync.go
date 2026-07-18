package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"faculty/internal/model"

	"github.com/google/uuid"
)

const SystemAuthorID = "00000000-0000-0000-0000-000000000001"

const (
	syncPageSize = 50

	maxTitleLen = 127
	maxBodyLen  = 255
	maxPlaceLen = 63

	categoryOnline  = "online"
	categoryOffline = "offline"
)

type eventSyncRepository interface {
	SyncExternalEvents(ctx context.Context, events []model.UpsertExternalEventParams) (int, error)
}

type eventSyncClient interface {
	ListPublicEvents(ctx context.Context, cookie string, limit, offset int, endDateGTE time.Time) (*model.CuEventListResponse, error)
}

type EventSyncService struct {
	repo eventSyncRepository
	cu   eventSyncClient
}

func NewEventSyncService(repo eventSyncRepository, cu eventSyncClient) *EventSyncService {
	return &EventSyncService{repo: repo, cu: cu}
}

func (s *EventSyncService) Sync(ctx context.Context, cookie string) (int, error) {
	authorID, err := uuid.Parse(SystemAuthorID)
	if err != nil {
		return 0, fmt.Errorf("parse system author id: %w", err)
	}

	now := time.Now()
	offset := 0
	var events []model.UpsertExternalEventParams

	for {
		resp, err := s.cu.ListPublicEvents(ctx, cookie, syncPageSize, offset, now)
		if err != nil {
			return 0, fmt.Errorf("fetch cu events: %w", err)
		}

		for _, item := range resp.Items {
			events = append(events, mapCuEvent(item, authorID))
		}

		offset += len(resp.Items)
		if len(resp.Items) == 0 || offset >= resp.Paging.TotalCount {
			break
		}
	}

	processed, err := s.repo.SyncExternalEvents(ctx, events)
	if err != nil {
		return 0, fmt.Errorf("sync external events: %w", err)
	}
	return processed, nil
}

func mapCuEvent(item model.CuEvent, authorID uuid.UUID) model.UpsertExternalEventParams {
	title := firstNonEmpty(deref(item.Title), item.ShortTitle, item.Slug, "Мероприятие")
	content := firstNonEmpty(deref(item.Summary), title)

	place := ""
	if item.Location != nil {
		place = firstNonEmpty(item.Location.Label, item.Location.City)
	}

	link := firstNonEmptyPtr(item.RegistrationUri, item.Link, item.LandingLink)
	if link != nil {
		truncated := truncate(*link, maxBodyLen)
		link = &truncated
	}

	return model.UpsertExternalEventParams{
		ExternalID:       item.ID,
		AuthorID:         authorID,
		Title:            truncate(title, maxTitleLen),
		Content:          truncate(content, maxBodyLen),
		Place:            truncate(place, maxPlaceLen),
		Category:         categoryFromFormat(item.Format),
		StartsAt:         item.StartDate,
		RegistrationLink: link,
	}
}

func categoryFromFormat(format []model.CuEventTag) string {
	for _, f := range format {
		if strings.EqualFold(f.Key, "Online") || strings.EqualFold(f.Value, "Онлайн") {
			return categoryOnline
		}
	}
	return categoryOffline
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func firstNonEmptyPtr(values ...*string) *string {
	for _, v := range values {
		if v != nil && strings.TrimSpace(*v) != "" {
			return v
		}
	}
	return nil
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}
