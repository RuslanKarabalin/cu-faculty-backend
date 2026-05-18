package model

import "github.com/google/uuid"

type CreateWorkPlaceParams struct {
	UserId       uuid.UUID
	CompanyName  string
	Grade        string
	Position     string
	StartYear    int
	EndYear      *int
	IsWorkingNow bool
}

type UpdateWorkPlaceParams struct {
	ID           int
	UserId       uuid.UUID
	CompanyName  string
	Grade        string
	Position     string
	StartYear    int
	EndYear      *int
	IsWorkingNow bool
}

type WorkPlaceRequest struct {
	CompanyName  string `json:"companyName"`
	Grade        string `json:"grade"`
	Position     string `json:"position"`
	StartYear    int    `json:"startYear"`
	EndYear      *int   `json:"endYear"`
	IsWorkingNow bool   `json:"isWorkingNow"`
}

type WorkPlace struct {
	ID           int    `json:"id"`
	CompanyName  string `json:"companyName"`
	Grade        string `json:"grade"`
	Position     string `json:"position"`
	StartYear    int    `json:"startYear"`
	EndYear      *int   `json:"endYear"`
	IsWorkingNow bool   `json:"isWorkingNow"`
}
