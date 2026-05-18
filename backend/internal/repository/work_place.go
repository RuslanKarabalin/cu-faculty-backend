package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateWorkPlace(ctx context.Context, params model.CreateWorkPlaceParams) (int, error) {
	query := `
	insert into work_places(user_id, company_name, grade, position, start_year, end_year, is_working_now)
	values($1, $2, $3, $4, $5, $6, $7)
	returning id
	`

	var id int
	if err := r.db.QueryRow(ctx, query, params.UserId, params.CompanyName, params.Grade, params.Position, params.StartYear, params.EndYear, params.IsWorkingNow).Scan(&id); err != nil {
		return 0, wrapPgError(err)
	}
	return id, nil
}

func (r *Repository) UpdateWorkPlace(ctx context.Context, params model.UpdateWorkPlaceParams) error {
	query := `
	update work_places
	set company_name = $3
		, grade = $4
		, position = $5
		, start_year = $6
		, end_year = $7
		, is_working_now = $8
	where id = $1 and user_id = $2
	`

	tag, err := r.db.Exec(ctx, query, params.ID, params.UserId, params.CompanyName, params.Grade, params.Position, params.StartYear, params.EndYear, params.IsWorkingNow)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteWorkPlace(ctx context.Context, id int, userID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from work_places where id = $1 and user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete work place: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetWorkPlaceByID(ctx context.Context, id int) (*model.WorkPlace, error) {
	query := `
	select
		wp.id
		, wp.company_name
		, wp.grade
		, wp.position
		, wp.start_year
		, wp.end_year
		, wp.is_working_now
	from work_places wp
	where wp.id = $1
	`

	wp := &model.WorkPlace{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&wp.ID,
		&wp.CompanyName,
		&wp.Grade,
		&wp.Position,
		&wp.StartYear,
		&wp.EndYear,
		&wp.IsWorkingNow,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get work place by id: %w", err)
	}
	return wp, nil
}

func (r *Repository) GetWorkPlacesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.WorkPlace, error) {
	query := `
	select
		wp.id
		, wp.company_name
		, wp.grade
		, wp.position
		, wp.start_year
		, wp.end_year
		, wp.is_working_now
	from work_places wp
	where wp.user_id = $1
	order by wp.start_year desc, wp.id
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select work places: %w", err)
	}
	defer rows.Close()

	places := make([]*model.WorkPlace, 0)
	for rows.Next() {
		wp := &model.WorkPlace{}
		if err := rows.Scan(
			&wp.ID,
			&wp.CompanyName,
			&wp.Grade,
			&wp.Position,
			&wp.StartYear,
			&wp.EndYear,
			&wp.IsWorkingNow,
		); err != nil {
			return nil, fmt.Errorf("failed to scan work place: %w", err)
		}
		places = append(places, wp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return places, nil
}
