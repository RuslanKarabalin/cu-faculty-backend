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
	Title            *string    `json:"title"`
	Content          *string    `json:"content"`
	Place            *string    `json:"place"`
	Category         *string    `json:"category"`
	StartsAt         *time.Time `json:"startsAt"`
	RegistrationLink *string    `json:"registrationLink"`
	IsDraft          *bool      `json:"isDraft"`
}

type UpsertExternalEventParams struct {
	ExternalID       int64
	AuthorID         uuid.UUID
	Title            string
	Content          string
	Place            string
	Category         string
	StartsAt         time.Time
	RegistrationLink *string
}

type CuEventListRequest struct {
	Paging CuEventPaging `json:"paging"`
	Filter CuEventFilter `json:"filter"`
}

type CuEventPaging struct {
	Limit   int           `json:"limit"`
	Offset  int           `json:"offset"`
	Sorting []CuEventSort `json:"sorting"`
}

type CuEventSort struct {
	By    string `json:"by"`
	IsAsc bool   `json:"isAsc"`
}

type CuEventFilter struct {
	EndDateGreaterThanOrEqualTo string `json:"endDateGreaterThanOrEqualTo,omitempty"`
}

type CuEventListResponse struct {
	Paging struct {
		Limit      int `json:"limit"`
		Offset     int `json:"offset"`
		TotalCount int `json:"totalCount"`
	} `json:"paging"`
	Items []CuEvent `json:"items"`
}

type CuEvent struct {
	ID              int64            `json:"id"`
	ShortTitle      string           `json:"shortTitle"`
	Title           *string          `json:"title"`
	Slug            string           `json:"slug"`
	Summary         *string          `json:"summary"`
	StartDate       time.Time        `json:"startDate"`
	Location        *CuEventLocation `json:"location"`
	Format          []CuEventTag     `json:"format"`
	RegistrationUri *string          `json:"registrationUri"`
	Link            *string          `json:"link"`
	LandingLink     *string          `json:"landingLink"`
}

type CuEventLocation struct {
	City  string  `json:"city"`
	Label string  `json:"label"`
	Uri   *string `json:"uri"`
}

type CuEventTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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
