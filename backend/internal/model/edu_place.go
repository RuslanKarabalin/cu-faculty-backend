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

type UpdateEduPlaceParams struct {
	ID             int
	UserId         uuid.UUID
	UniversityId   int
	Grade          string
	Level          *string
	Specialization string
	StartYear      int
	EndYear        *int
	IsStudyingNow  bool
}

type EduPlaceRequest struct {
	UniversityId   int     `json:"universityId"`
	Grade          string  `json:"grade"`
	Level          *string `json:"level"`
	Specialization string  `json:"specialization"`
	StartYear      int     `json:"startYear"`
	EndYear        *int    `json:"endYear"`
	IsStudyingNow  bool    `json:"isStudyingNow"`
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
