package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

const defaultPublishDays = 7

type newsRepository interface {
	CreateNews(ctx context.Context, params model.CreateNewsParams) (uuid.UUID, error)
	UpdateNews(ctx context.Context, params model.UpdateNewsParams) error
	UpdateNewsPhoto(ctx context.Context, id, authorID uuid.UUID, key string) (*string, error)
	DeleteNews(ctx context.Context, id, authorID uuid.UUID) error
	GetNewsByID(ctx context.Context, id uuid.UUID) (*model.News, error)
	GetNews(ctx context.Context, limit, offset int) ([]*model.News, int, error)
	GetNewsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.News, int, error)
}

type NewsService struct {
	repo newsRepository
}

func NewNewsService(repo newsRepository) *NewsService {
	return &NewsService{repo: repo}
}

func (s *NewsService) GetNewsByID(ctx context.Context, id uuid.UUID) (*model.News, error) {
	return s.repo.GetNewsByID(ctx, id)
}

func (s *NewsService) GetNews(ctx context.Context, limit, offset int) ([]*model.News, int, error) {
	return s.repo.GetNews(ctx, limit, offset)
}

func (s *NewsService) GetNewsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.News, int, error) {
	return s.repo.GetNewsByAuthorID(ctx, authorID, limit, offset)
}

func (s *NewsService) CreateNews(ctx context.Context, authorID uuid.UUID, req model.CreateNewsRequest) (*model.News, error) {
	id, err := s.repo.CreateNews(ctx, model.CreateNewsParams{
		AuthorID:    authorID,
		Title:       req.Title,
		Content:     req.Content,
		PublishDays: normalizePublishDays(req.PublishDays),
		IsDraft:     req.IsDraft,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetNewsByID(ctx, id)
}

func (s *NewsService) UpdateNews(ctx context.Context, authorID, id uuid.UUID, req model.UpdateNewsRequest) (*model.News, error) {
	err := s.repo.UpdateNews(ctx, model.UpdateNewsParams{
		ID:          id,
		AuthorID:    authorID,
		Title:       req.Title,
		Content:     req.Content,
		PublishDays: normalizePublishDays(req.PublishDays),
		IsDraft:     req.IsDraft,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetNewsByID(ctx, id)
}

// SetPhoto updates the news photo and returns the updated news along with the
// previous photo key (if any) so the caller can delete the replaced object.
func (s *NewsService) SetPhoto(ctx context.Context, authorID, id uuid.UUID, key string) (*model.News, *string, error) {
	oldKey, err := s.repo.UpdateNewsPhoto(ctx, id, authorID, key)
	if err != nil {
		return nil, nil, err
	}
	news, err := s.repo.GetNewsByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return news, oldKey, nil
}

func (s *NewsService) DeleteNews(ctx context.Context, authorID, id uuid.UUID) error {
	return s.repo.DeleteNews(ctx, id, authorID)
}

func normalizePublishDays(days int) int {
	if days <= 0 {
		return defaultPublishDays
	}
	return days
}
