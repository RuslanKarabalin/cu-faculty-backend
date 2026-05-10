package model

import "github.com/google/uuid"

type CuEduPlaceResp struct {
	EducationProgram struct {
		Name      string `json:"name"`
		Level     string `json:"level"`
		StartDate string `json:"startDate"`
	} `json:"educationProgram"`
}

type CreateEduPlaceParams struct {
	UserId         uuid.UUID
	UniversityId   int
	Grade          string
	Level          *string
	Specialization string
	StartYear      int
	EndYear        *int
	IsStudyingNow  bool
}

type EduPlace struct {
	ID             int     `json:"id"`
	UniversityName string  `json:"universityName"`
	Grade          string  `json:"grade"`
	Level          *string `json:"level"`
	Specialization string  `json:"specialization"`
	StartYear      int     `json:"startYear"`
	EndYear        *int    `json:"endYear"`
	IsStudyingNow  bool    `json:"isStudyingNow"`
}
