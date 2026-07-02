package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateAnnouncementParams struct {
	AuthorID uuid.UUID
	Title    string
	Content  string
}

type UpdateAnnouncementParams struct {
	ID         uuid.UUID
	AuthorID   uuid.UUID
	Title      string
	Content    string
	IsArchived bool
}

type CreateAnnouncementRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateAnnouncementRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	IsArchived bool   `json:"isArchived"`
}

type Announcement struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	IsArchived bool      `json:"isArchived"`
	CreatedAt  time.Time `json:"createdAt"`
	Author     *User     `json:"author"`
}
