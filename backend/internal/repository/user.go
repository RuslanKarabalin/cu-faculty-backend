package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateUser(ctx context.Context, params model.CreateUserParams) error {
	query := `insert into users(id, first_name, last_name, birth_date, role) values($1, $2, $3, $4, 'user')`

	if _, err := r.db.Exec(ctx, query, params.ID, params.FirstName, params.LastName, params.BirthDate); err != nil {
		return wrapPgError(err)
	}
	return nil
}

func (r *Repository) UpdateUserPhoto(ctx context.Context, id uuid.UUID, key string) error {
	tag, err := r.db.Exec(ctx, `update users set photo_s3_key = $2 where id = $1`, id, key)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from users where role = 'user'`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	query := `
	select
		u.id
		, u.photo_s3_key
		, u.first_name
		, u.last_name
		, u.bio
		, u.birth_date
		, u.speciality
		, st.content
		, u.role
	from users u
	left join statuses st on st.id = u.status_id
	where u.role = 'user'
	order by u.last_name, u.first_name, u.id
	limit $1 offset $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to select users: %w", err)
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(
			&u.ID,
			&u.PhotoS3Key,
			&u.FirstName,
			&u.LastName,
			&u.Bio,
			&u.BirthDate,
			&u.Speciality,
			&u.Status,
			&u.Role,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan users: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return users, total, nil
}

func (r *Repository) UpdateUser(ctx context.Context, params model.UpdateUserParams) error {
	query := `
	update users
	set photo_s3_key = $2
		, bio = $3
		, speciality = $4
		, status_id = $5
	where id = $1
	`

	tag, err := r.db.Exec(ctx, query, params.ID, params.PhotoS3Key, params.Bio, params.Speciality, params.StatusID)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
	select
		u.id
		, u.photo_s3_key
		, u.first_name
		, u.last_name
		, u.bio
		, u.birth_date
		, u.speciality
		, st.content
		, u.role
	from users u
	left join statuses st on st.id = u.status_id
	where u.id = $1
	`

	u := &model.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.PhotoS3Key,
		&u.FirstName,
		&u.LastName,
		&u.Bio,
		&u.BirthDate,
		&u.Speciality,
		&u.Status,
		&u.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return u, nil
}
