package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateEventParams struct {
	AuthorID         uuid.UUID
	Title            string
	Content          string
	Place            string
	Category         string
	StartsAt         time.Time
	RegistrationLink *string
	IsDraft          bool
}

type UpdateEventParams struct {
	ID               uuid.UUID
	AuthorID         uuid.UUID
	Title            string
	Content          string
	Place            string
	Category         string
	StartsAt         time.Time
	RegistrationLink *string
	IsDraft          bool
}

type CreateEventRequest struct {
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	Place            string    `json:"place"`
	Category         string    `json:"category"`
	StartsAt         time.Time `json:"startsAt"`
	RegistrationLink *string   `json:"registrationLink"`
	IsDraft          bool      `json:"isDraft"`
}

type UpdateEventRequest struct {
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	Place            string    `json:"place"`
	Category         string    `json:"category"`
	StartsAt         time.Time `json:"startsAt"`
	RegistrationLink *string   `json:"registrationLink"`
	IsDraft          bool      `json:"isDraft"`
}

type Event struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	PhotoS3Key       *string   `json:"-"`
	PhotoURL         *string   `json:"photoUrl"`
	Place            string    `json:"place"`
	Category         string    `json:"category"`
	StartsAt         time.Time `json:"startsAt"`
	RegistrationLink *string   `json:"registrationLink"`
	IsDraft          bool      `json:"isDraft"`
	CreatedAt        time.Time `json:"createdAt"`
	Author           *User     `json:"author"`
}
