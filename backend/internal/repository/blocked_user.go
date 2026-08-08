package repository

import (
	"context"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
)

// IsBlocked reports whether blockerID has blocked blockedUserID.
func (r *Repository) IsBlocked(ctx context.Context, blockerID, blockedUserID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		select exists(
			select 1 from blocked_users
			where user_id = $1 and blocked_user_id = $2
		)
	`, blockerID, blockedUserID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check block: %w", err)
	}
	return exists, nil
}

// HasBlockBetween reports whether either user has blocked the other.
func (r *Repository) HasBlockBetween(ctx context.Context, a, b uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		select exists(
			select 1 from blocked_users
			where (user_id = $1 and blocked_user_id = $2)
			   or (user_id = $2 and blocked_user_id = $1)
		)
	`, a, b).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check block between users: %w", err)
	}
	return exists, nil
}

func (r *Repository) AddBlockedUser(ctx context.Context, userID, blockedUserID uuid.UUID) error {
	return r.RunInTx(ctx, func(tx *Repository) error {
		query := `
		insert into blocked_users(user_id, blocked_user_id)
		values($1, $2)
		on conflict do nothing
		`
		if _, err := tx.db.Exec(ctx, query, userID, blockedUserID); err != nil {
			return wrapPgError(err)
		}

		// Drop mutual contact / saved links so the block takes effect immediately.
		if _, err := tx.db.Exec(ctx, `
			delete from contacts
			where (user_id = $1 and contact_id = $2)
			   or (user_id = $2 and contact_id = $1)
		`, userID, blockedUserID); err != nil {
			return fmt.Errorf("failed to clear contacts on block: %w", err)
		}
		if _, err := tx.db.Exec(ctx, `
			delete from saved_users
			where (user_id = $1 and saved_user_id = $2)
			   or (user_id = $2 and saved_user_id = $1)
		`, userID, blockedUserID); err != nil {
			return fmt.Errorf("failed to clear saved users on block: %w", err)
		}
		return nil
	})
}

func (r *Repository) DeleteBlockedUser(ctx context.Context, userID, blockedUserID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from blocked_users where user_id = $1 and blocked_user_id = $2`, userID, blockedUserID)
	if err != nil {
		return fmt.Errorf("failed to delete blocked user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetBlockedUsers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.User, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from blocked_users where user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count blocked users: %w", err)
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
	from blocked_users bu
	join users u on u.id = bu.blocked_user_id
	left join statuses st on st.id = u.status_id
	where bu.user_id = $1
	order by u.last_name, u.first_name, u.id
	limit $2 offset $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to select blocked users: %w", err)
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
			return nil, 0, fmt.Errorf("failed to scan blocked user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return users, total, nil
}
