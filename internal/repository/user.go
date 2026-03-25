package repository

import (
	"context"
	"faculty/internal/model"
	"fmt"
)

func (r *Repository) CreateUser(ctx context.Context, cuUserResp model.CuUserResp) error {
	query := `insert into users(id, first_name, last_name, birth_date, role) values($1, $2, $3, $4, 'user')`

	_, err := r.pgPool.Exec(
		ctx,
		query,
		cuUserResp.Id,
		cuUserResp.FirstName,
		cuUserResp.LastName,
		cuUserResp.BirthDate,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	query := `
	select
		u.id
		, u.photo_s3_key
		, u.first_name
		, u.last_name
		, u.bio
		, u.birth_date
		, st.content
		, u.role
	from users u
	left join statuses st on st.id = u.status_id
	where u.role = 'user'
	`

	rows, err := r.pgPool.Query(
		ctx,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
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
			&t.Status,
			&t.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan users: %w", err)
		}
		users = append(users, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return users, nil
}
