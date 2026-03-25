package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `json:"id"`
	PhotoS3Key *string   `json:"photoS3Key"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Bio        *string   `json:"bio"`
	BirthDate  time.Time `json:"birthdate"`
	Status     *string   `json:"status"`
	Role       string    `json:"role"`
}

type CuUserResp struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	BirthDate string    `json:"birthdate"`
}
