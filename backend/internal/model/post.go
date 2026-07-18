package model

import (
	"time"

	"github.com/google/uuid"
)

type CreatePostParams struct {
	AuthorID uuid.UUID
	Title    string
	Content  string
}

type UpdatePostParams struct {
	ID         uuid.UUID
	AuthorID   uuid.UUID
	Title      string
	Content    string
	IsArchived bool
}

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	IsArchived bool   `json:"isArchived"`
}

type Post struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	IsArchived bool      `json:"isArchived"`
	CreatedAt  time.Time `json:"createdAt"`
	Author     *User     `json:"author"`
}
