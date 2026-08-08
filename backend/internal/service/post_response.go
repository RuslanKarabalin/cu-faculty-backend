package service

import (
	"context"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type postResponseRepository interface {
	AddPostResponse(ctx context.Context, userID, postID uuid.UUID) error
	DeletePostResponse(ctx context.Context, userID, postID uuid.UUID) error
	GetPostResponders(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*model.User, int, error)
	GetPostsRespondedByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Post, int, error)
	GetVisiblePostByID(ctx context.Context, id, viewerID uuid.UUID) (*model.Post, error)
}

type PostResponseService struct {
	repo postResponseRepository
}

func NewPostResponseService(repo postResponseRepository) *PostResponseService {
	return &PostResponseService{repo: repo}
}

func (s *PostResponseService) Respond(ctx context.Context, userID, postID uuid.UUID) (*model.Post, error) {
	post, err := s.repo.GetVisiblePostByID(ctx, postID, userID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.AddPostResponse(ctx, userID, postID); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostResponseService) DeleteResponse(ctx context.Context, userID, postID uuid.UUID) error {
	return s.repo.DeletePostResponse(ctx, userID, postID)
}

func (s *PostResponseService) GetResponders(ctx context.Context, viewerID, postID uuid.UUID, limit, offset int) ([]*model.User, int, error) {
	if _, err := s.repo.GetVisiblePostByID(ctx, postID, viewerID); err != nil {
		return nil, 0, err
	}
	return s.repo.GetPostResponders(ctx, postID, limit, offset)
}

func (s *PostResponseService) GetMyResponses(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Post, int, error) {
	return s.repo.GetPostsRespondedByUser(ctx, userID, limit, offset)
}
