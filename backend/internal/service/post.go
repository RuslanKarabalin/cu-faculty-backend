package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type postRepository interface {
	CreatePost(ctx context.Context, params model.CreatePostParams) (uuid.UUID, error)
	UpdatePost(ctx context.Context, params model.UpdatePostParams) error
	DeletePost(ctx context.Context, id, authorID uuid.UUID) error
	GetPostByID(ctx context.Context, id uuid.UUID) (*model.Post, error)
	GetVisiblePostByID(ctx context.Context, id, viewerID uuid.UUID) (*model.Post, error)
	GetPosts(ctx context.Context, viewerID uuid.UUID, limit, offset int) ([]*model.Post, int, error)
	GetPostsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Post, int, error)
}

type PostService struct {
	repo postRepository
}

func NewPostService(repo postRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) GetPostByID(ctx context.Context, id, viewerID uuid.UUID) (*model.Post, error) {
	return s.repo.GetVisiblePostByID(ctx, id, viewerID)
}

func (s *PostService) GetPosts(ctx context.Context, viewerID uuid.UUID, limit, offset int) ([]*model.Post, int, error) {
	return s.repo.GetPosts(ctx, viewerID, limit, offset)
}

func (s *PostService) GetPostsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Post, int, error) {
	return s.repo.GetPostsByAuthorID(ctx, authorID, limit, offset)
}

func (s *PostService) CreatePost(ctx context.Context, authorID uuid.UUID, req model.CreatePostRequest) (*model.Post, error) {
	id, err := s.repo.CreatePost(ctx, model.CreatePostParams{
		AuthorID: authorID,
		Title:    req.Title,
		Content:  req.Content,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetPostByID(ctx, id)
}

func (s *PostService) UpdatePost(ctx context.Context, authorID, id uuid.UUID, req model.UpdatePostRequest) (*model.Post, error) {
	err := s.repo.UpdatePost(ctx, model.UpdatePostParams{
		ID:         id,
		AuthorID:   authorID,
		Title:      req.Title,
		Content:    req.Content,
		IsArchived: req.IsArchived,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetPostByID(ctx, id)
}

func (s *PostService) DeletePost(ctx context.Context, authorID, id uuid.UUID) error {
	return s.repo.DeletePost(ctx, id, authorID)
}
