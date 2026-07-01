package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	PhotoS3Key *string   `json:"photoS3Key"`
	PhotoURL   *string   `json:"photoUrl"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Bio        *string   `json:"bio"`
	BirthDate  Date      `json:"birthdate"`
	Speciality *string   `json:"speciality"`
	Status     *string   `json:"status"`
	Role       string    `json:"role"`
}

type CuUserResp struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	BirthDate string    `json:"birthdate"`
}

type CreateUserParams struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	BirthDate time.Time
}

type UpdateUserRequest struct {
	PhotoS3Key *string `json:"photoS3Key"`
	Bio        *string `json:"bio"`
	Speciality *string `json:"speciality"`
	StatusID   *int    `json:"statusId"`
}

type UpdateUserParams struct {
	ID         uuid.UUID
	PhotoS3Key *string
	Bio        *string
	Speciality *string
	StatusID   *int
}
