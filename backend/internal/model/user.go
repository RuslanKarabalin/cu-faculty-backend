package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	PhotoS3Key *string   `json:"-"`
	PhotoURL   *string   `json:"photoUrl"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Bio        *string   `json:"bio"`
	BirthDate  Date      `json:"birthdate"`
	Speciality *string   `json:"speciality"`
	Status     *string   `json:"status"`
	Role       string    `json:"role"`
}

type UserRelation string

const (
	UserRelationSaved   UserRelation = "saved"
	UserRelationContact UserRelation = "contact"
	UserRelationOther   UserRelation = "other"
)

type UserSearchResult struct {
	*User
	Relation UserRelation `json:"relation"`
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
	Bio        *string `json:"bio"`
	Speciality *string `json:"speciality"`
	StatusID   *int    `json:"statusId"`
}

type UpdateUserParams struct {
	ID         uuid.UUID
	Bio        *string
	Speciality *string
	StatusID   *int
}

type ProfileCompletenessData struct {
	HasPhoto        bool
	HasSpeciality   bool
	HasBio          bool
	HasBirthDate    bool
	SocialsCount    int
	EduPlacesCount  int
	WorkPlacesCount int
	KeySkillsCount  int
	SoftSkillsCount int
}

const (
	completenessPhoto        = 15
	completenessSpeciality   = 15
	completenessBio          = 10
	completenessBirthDate    = 5
	completenessPerSocial    = 5
	completenessPerEduPlace  = 5
	completenessPerSkill     = 5
	completenessPerWorkPlace = 15
	completenessMax          = 100
)

func (d ProfileCompletenessData) Percent() int {
	total := 0
	if d.HasPhoto {
		total += completenessPhoto
	}
	if d.HasSpeciality {
		total += completenessSpeciality
	}
	if d.HasBio {
		total += completenessBio
	}
	if d.HasBirthDate {
		total += completenessBirthDate
	}
	total += completenessPerSocial * d.SocialsCount
	total += completenessPerEduPlace * d.EduPlacesCount
	total += completenessPerSkill * (d.KeySkillsCount + d.SoftSkillsCount)
	total += completenessPerWorkPlace * d.WorkPlacesCount
	if total > completenessMax {
		total = completenessMax
	}
	return total
}

type ProfileCompleteness struct {
	Percent int `json:"percent"`
}
