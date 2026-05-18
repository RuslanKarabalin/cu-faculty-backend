package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateSocial(ctx context.Context, params model.CreateSocialParams) (int, error) {
	query := `
	insert into socials(user_id, social, link, is_preferred)
	values($1, $2, $3, $4)
	returning id
	`

	var id int
	if err := r.db.QueryRow(ctx, query, params.UserId, params.Social, params.Link, params.IsPreferred).Scan(&id); err != nil {
		return 0, wrapPgError(err)
	}
	return id, nil
}

func (r *Repository) UpdateSocial(ctx context.Context, params model.UpdateSocialParams) error {
	query := `
	update socials
	set social = $3
		, link = $4
		, is_preferred = $5
	where id = $1 and user_id = $2
	`

	tag, err := r.db.Exec(ctx, query, params.ID, params.UserId, params.Social, params.Link, params.IsPreferred)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteSocial(ctx context.Context, id int, userID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from socials where id = $1 and user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete social: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetSocialByID(ctx context.Context, id int) (*model.Social, error) {
	query := `
	select
		s.id
		, s.social
		, s.link
		, s.is_preferred
	from socials s
	where s.id = $1
	`

	s := &model.Social{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID,
		&s.Social,
		&s.Link,
		&s.IsPreferred,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get social by id: %w", err)
	}
	return s, nil
}

func (r *Repository) GetSocialsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Social, error) {
	query := `
	select
		s.id
		, s.social
		, s.link
		, s.is_preferred
	from socials s
	where s.user_id = $1
	order by s.is_preferred desc, s.id
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select socials: %w", err)
	}
	defer rows.Close()

	socials := make([]*model.Social, 0)
	for rows.Next() {
		s := &model.Social{}
		if err := rows.Scan(
			&s.ID,
			&s.Social,
			&s.Link,
			&s.IsPreferred,
		); err != nil {
			return nil, fmt.Errorf("failed to scan social: %w", err)
		}
		socials = append(socials, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return socials, nil
}
