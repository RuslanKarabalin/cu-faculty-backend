package model

import "github.com/google/uuid"

type CreateSocialParams struct {
	UserId      uuid.UUID
	Social      string
	Link        string
	IsPreferred bool
}

type UpdateSocialParams struct {
	ID          int
	UserId      uuid.UUID
	Social      string
	Link        string
	IsPreferred bool
}

type SocialRequest struct {
	Social      string `json:"social"`
	Link        string `json:"link"`
	IsPreferred bool   `json:"isPreferred"`
}

type Social struct {
	ID          int    `json:"id"`
	Social      string `json:"social"`
	Link        string `json:"link"`
	IsPreferred bool   `json:"isPreferred"`
}
