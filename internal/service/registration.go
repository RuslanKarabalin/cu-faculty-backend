package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"faculty/internal/cuclient"
	"faculty/internal/model"
	"faculty/internal/repository"
)

const cuUniversityID = 207

type RegistrationService struct {
	repo     *repository.Repository
	cuClient *cuclient.Client
}

func NewRegistrationService(repo *repository.Repository, cuClient *cuclient.Client) *RegistrationService {
	return &RegistrationService{repo: repo, cuClient: cuClient}
}

func (s *RegistrationService) Register(ctx context.Context, cuUser model.CuUserResp, cookie string) (*model.User, bool, error) {
	if existing, err := s.repo.GetUserByID(ctx, cuUser.ID); err == nil {
		return existing, false, nil
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, false, err
	}

	birthDate, err := time.Parse(model.DateLayout, cuUser.BirthDate)
	if err != nil {
		return nil, false, ErrInvalidBirthDate
	}

	cuEduPlaces, err := s.cuClient.StudentEduInfo(ctx, cookie)
	if err != nil {
		return nil, false, fmt.Errorf("fetch student edu info: %w", err)
	}

	eduPlaceParams, err := buildEduPlaceParams(cuEduPlaces)
	if err != nil {
		return nil, false, err
	}

	isNewUser := true
	err = s.repo.RunInTx(ctx, func(r *repository.Repository) error {
		err := r.CreateUser(ctx, model.CreateUserParams{
			ID:        cuUser.ID,
			FirstName: cuUser.FirstName,
			LastName:  cuUser.LastName,
			BirthDate: birthDate,
		})
		if err != nil {
			if errors.Is(err, repository.ErrDuplicate) {
				isNewUser = false
				return nil
			}
			return err
		}

		for _, p := range eduPlaceParams {
			p.UserId = cuUser.ID
			if _, err := r.CreateEduPlace(ctx, p); err != nil {
				if errors.Is(err, repository.ErrDuplicate) {
					continue
				}
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	user, err := s.repo.GetUserByID(ctx, cuUser.ID)
	if err != nil {
		return nil, false, err
	}
	return user, isNewUser, nil
}

func buildEduPlaceParams(cuEduPlaces []model.CuEduPlaceResp) ([]model.CreateEduPlaceParams, error) {
	params := make([]model.CreateEduPlaceParams, 0, len(cuEduPlaces))
	for _, e := range cuEduPlaces {
		t, err := time.Parse(model.DateLayout, e.EducationProgram.StartDate)
		if err != nil {
			return nil, ErrInvalidUpstreamData
		}
		params = append(params, model.CreateEduPlaceParams{
			UniversityId:   cuUniversityID,
			Grade:          strings.ToLower(e.EducationProgram.Level),
			Specialization: e.EducationProgram.Name,
			StartYear:      t.Year(),
			IsStudyingNow:  true,
		})
	}
	return params, nil
}
