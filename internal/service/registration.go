package service

import (
	"context"
	"errors"
	"time"

	"faculty/internal/model"
	"faculty/internal/repository"
)

type RegistrationService struct {
	repo *repository.Repository
}

func NewRegistrationService(repo *repository.Repository) *RegistrationService {
	return &RegistrationService{repo: repo}
}

func (s *RegistrationService) Register(ctx context.Context, cuUser model.CuUserResp, eduPlaces []model.CreateEduPlaceParams) (*model.User, bool, error) {
	birthDate, err := time.Parse("2006-01-02", cuUser.BirthDate)
	if err != nil {
		return nil, false, ErrInvalidBirthDate
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

		for _, p := range eduPlaces {
			p.UserId = cuUser.ID
			if err := r.CreateEduPlace(ctx, p); err != nil {
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
