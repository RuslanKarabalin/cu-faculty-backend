package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) CreateUser(ctx context.Context, params model.CreateUserParams) error {
	query := `insert into users(id, first_name, last_name, birth_date, role) values($1, $2, $3, $4, 'user')`

	_, err := r.pgPool.Exec(ctx, query, params.ID, params.FirstName, params.LastName, params.BirthDate)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == pgUniqueViolation {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (r *Repository) GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error) {
	var total int
	err := r.pgPool.QueryRow(ctx, `select count(*) from users where role = 'user'`).Scan(&total)
	if err != nil {
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
	limit $1 offset $2
	`

	rows, err := r.pgPool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to select users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		t := &model.User{}
		err := rows.Scan(
			&t.Id,
			&t.PhotoS3Key,
			&t.FirstName,
			&t.LastName,
			&t.Bio,
			&t.BirthDate,
			&t.Speciality,
			&t.Status,
			&t.Role,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan users: %w", err)
		}
		users = append(users, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return users, total, nil
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
	err := r.pgPool.QueryRow(ctx, query, id).Scan(
		&u.Id,
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
