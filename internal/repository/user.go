package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) CreateUser(ctx context.Context, cuUserResp model.CuUserResp) error {
	birthDate, err := time.Parse("2006-01-02", cuUserResp.BirthDate)
	if err != nil {
		return ErrInvalidDate
	}

	query := `insert into users(id, first_name, last_name, birth_date, role) values($1, $2, $3, $4, 'user')`

	_, err = r.pgPool.Exec(
		ctx,
		query,
		cuUserResp.Id,
		cuUserResp.FirstName,
		cuUserResp.LastName,
		birthDate,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (r *Repository) GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error) {
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
		, count(*) over() as total
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
	var total int
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
			&total,
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
