package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateNewsParams struct {
	AuthorID    uuid.UUID
	Title       string
	Content     string
	PublishDays int
	IsDraft     bool
}

type UpdateNewsParams struct {
	ID          uuid.UUID
	AuthorID    uuid.UUID
	Title       string
	Content     string
	PublishDays int
	IsDraft     bool
}

type CreateNewsRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	PublishDays int    `json:"publishDays"`
	IsDraft     bool   `json:"isDraft"`
}

type UpdateNewsRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	PublishDays int    `json:"publishDays"`
	IsDraft     bool   `json:"isDraft"`
}

type News struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PhotoS3Key  *string   `json:"photoS3Key"`
	PhotoURL    *string   `json:"photoUrl"`
	PublishDays int       `json:"publishDays"`
	IsDraft     bool      `json:"isDraft"`
	CreatedAt   time.Time `json:"createdAt"`
	Author      *User     `json:"author"`
}
