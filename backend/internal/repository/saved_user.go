package repository

import (
	"context"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) AddSavedUser(ctx context.Context, userID, savedUserID uuid.UUID) error {
	query := `
	insert into saved_users(user_id, saved_user_id)
	values($1, $2)
	on conflict do nothing
	`

	if _, err := r.db.Exec(ctx, query, userID, savedUserID); err != nil {
		return wrapPgError(err)
	}
	return nil
}

func (r *Repository) DeleteSavedUser(ctx context.Context, userID, savedUserID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from saved_users where user_id = $1 and saved_user_id = $2`, userID, savedUserID)
	if err != nil {
		return fmt.Errorf("failed to delete saved user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetSavedUsers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.User, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from saved_users where user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count saved users: %w", err)
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
	from saved_users su
	join users u on u.id = su.saved_user_id
	left join statuses st on st.id = u.status_id
	where su.user_id = $1
	order by u.last_name, u.first_name, u.id
	limit $2 offset $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to select saved users: %w", err)
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
			return nil, 0, fmt.Errorf("failed to scan saved user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return users, total, nil
}
