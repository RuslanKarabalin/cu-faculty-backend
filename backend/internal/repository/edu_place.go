package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateEduPlace(ctx context.Context, params model.CreateEduPlaceParams) (int, error) {
	query := `
	insert into edu_places(user_id, university_id, grade, level, specialization, start_year, end_year, is_studying_now)
	values($1, $2, $3, $4, $5, $6, $7, $8)
	returning id
	`

	var id int
	if err := r.db.QueryRow(ctx, query, params.UserId, params.UniversityId, params.Grade, params.Level, params.Specialization, params.StartYear, params.EndYear, params.IsStudyingNow).Scan(&id); err != nil {
		return 0, wrapPgError(err)
	}
	return id, nil
}

func (r *Repository) UpdateEduPlace(ctx context.Context, params model.UpdateEduPlaceParams) error {
	query := `
	update edu_places
	set university_id = $3
		, grade = $4
		, level = $5
		, specialization = $6
		, start_year = $7
		, end_year = $8
		, is_studying_now = $9
	where id = $1 and user_id = $2
	`

	tag, err := r.db.Exec(ctx, query, params.ID, params.UserId, params.UniversityId, params.Grade, params.Level, params.Specialization, params.StartYear, params.EndYear, params.IsStudyingNow)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteEduPlace(ctx context.Context, id int, userID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from edu_places where id = $1 and user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete edu place: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetEduPlaceByID(ctx context.Context, id int) (*model.EduPlace, error) {
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
	where ep.id = $1
	`

	ep := &model.EduPlace{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&ep.ID,
		&ep.UniversityName,
		&ep.Grade,
		&ep.Level,
		&ep.Specialization,
		&ep.StartYear,
		&ep.EndYear,
		&ep.IsStudyingNow,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get edu place by id: %w", err)
	}
	return ep, nil
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

	rows, err := r.db.Query(ctx, query, userID)
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
