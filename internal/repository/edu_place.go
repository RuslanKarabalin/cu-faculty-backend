package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) CreateEduPlace(ctx context.Context, params model.CreateEduPlaceParams) error {
	query := `
	insert into edu_places(user_id, university_id, grade, level, specialization, start_year, end_year, is_studying_now)
	values($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.pgPool.Exec(ctx, query, params.UserId, params.UniversityId, params.Grade, params.Level, params.Specialization, params.StartYear, params.EndYear, params.IsStudyingNow)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == pgUniqueViolation {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (r *Repository) GetEduPlacesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EduPlace, error) {
	query := `
	select
		ep.id
		, u.name
		, ep.grade
		, ep.level
		, ep.specialization
		, ep.start_year
		, ep.end_year
		, ep.is_studying_now
	from edu_places ep
	join universities u on u.id = ep.university_id
	where ep.user_id = $1
	order by ep.start_year desc, ep.id
	`

	rows, err := r.pgPool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select edu places: %w", err)
	}
	defer rows.Close()

	places := make([]*model.EduPlace, 0)
	for rows.Next() {
		ep := &model.EduPlace{}
		if err := rows.Scan(
			&ep.ID,
			&ep.UniversityName,
			&ep.Grade,
			&ep.Level,
			&ep.Specialization,
			&ep.StartYear,
			&ep.EndYear,
			&ep.IsStudyingNow,
		); err != nil {
			return nil, fmt.Errorf("failed to scan edu place: %w", err)
		}
		places = append(places, ep)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return places, nil
}
