package model

import (
	"strings"
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

const maxSearchItems = 100

type SearchUsersRequest struct {
	Query        *string  `json:"q"`
	StatusID     *int     `json:"statusId"`
	Speciality   *string  `json:"speciality"`
	Companies    []string `json:"companies"`
	KeySkillIDs  []int    `json:"keySkillIds"`
	SoftSkillIDs []int    `json:"softSkillIds"`
	Limit        int      `json:"limit"`
	Offset       int      `json:"offset"`
}

type SearchUsersParams struct {
	Query        *string
	StatusID     *int
	Speciality   *string
	Companies    []string
	KeySkillIDs  []int
	SoftSkillIDs []int
	Limit        int
	Offset       int
}

func (r SearchUsersRequest) Normalize() SearchUsersParams {
	p := SearchUsersParams{
		Query:        trimToPtr(r.Query),
		Speciality:   trimToPtr(r.Speciality),
		Companies:    normalizeCompanies(r.Companies),
		KeySkillIDs:  cleanIDs(r.KeySkillIDs),
		SoftSkillIDs: cleanIDs(r.SoftSkillIDs),
	}
	if r.StatusID != nil && *r.StatusID > 0 {
		p.StatusID = r.StatusID
	}
	p.Limit, p.Offset = PageQuery{Limit: r.Limit, Offset: r.Offset}.Normalize()
	return p
}

func trimToPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

// normalizeCompanies trims, lowercases and dedups company names for
// case-insensitive matching against work_places.company_name.
func normalizeCompanies(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.ToLower(strings.TrimSpace(v))
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
		if len(out) >= maxSearchItems {
			break
		}
	}
	return out
}

func cleanIDs(ids []int) []int {
	seen := make(map[int]struct{}, len(ids))
	out := make([]int, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
		if len(out) >= maxSearchItems {
			break
		}
	}
	return out
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
