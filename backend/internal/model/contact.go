package model

import "github.com/google/uuid"

type CreateContactParams struct {
	UserID    uuid.UUID
	ContactID uuid.UUID
	Note      *string
}

type UpdateContactParams struct {
	UserID    uuid.UUID
	ContactID uuid.UUID
	Note      *string
}

type CreateContactRequest struct {
	ContactID uuid.UUID `json:"contactId"`
	Note      *string   `json:"note"`
}

type UpdateContactRequest struct {
	Note *string `json:"note"`
}

type Contact struct {
	Note *string `json:"note"`
	User *User   `json:"user"`
}
