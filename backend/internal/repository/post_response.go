package repository

import (
	"context"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) AddPostResponse(ctx context.Context, userID, postID uuid.UUID) error {
	query := `
	insert into post_responses(user_id, post_id)
	select $1, $2
	where exists (
		select 1 from posts a
		join users u on u.id = a.author_id
		where a.id = $2
			and not a.is_archived
			and u.deleted_at is null
			and not exists (
				select 1 from blocked_users bu
				where bu.user_id = a.author_id and bu.blocked_user_id = $1
			)
	)
	on conflict do nothing
	`

	tag, err := r.db.Exec(ctx, query, userID, postID)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		var alreadyResponded bool
		if err := r.db.QueryRow(ctx,
			`select exists(select 1 from post_responses where user_id = $1 and post_id = $2)`,
			userID, postID,
		).Scan(&alreadyResponded); err != nil {
			return fmt.Errorf("failed to check post response: %w", err)
		}
		if !alreadyResponded {
			return ErrNotFound
		}
	}
	return nil
}

func (r *Repository) DeletePostResponse(ctx context.Context, userID, postID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from post_responses where user_id = $1 and post_id = $2`, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post response: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetPostResponders(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*model.User, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from post_responses where post_id = $1`, postID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count post responders: %w", err)
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
	from post_responses ar
	join users u on u.id = ar.user_id
	left join statuses st on st.id = u.status_id
	where ar.post_id = $1
	order by u.last_name, u.first_name, u.id
	limit $2 offset $3
	`

	rows, err := r.db.Query(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to select post responders: %w", err)
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
			return nil, 0, fmt.Errorf("failed to scan responder: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return users, total, nil
}

func (r *Repository) GetPostsRespondedByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Post, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from post_responses where user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count post responses: %w", err)
	}

	query := postSelectColumns + `
	join post_responses ar on ar.post_id = a.id
	where ar.user_id = $1
	order by a.created_at desc, a.id
	limit $2 offset $3
	`

	return r.selectPosts(ctx, total, query, userID, limit, offset)
}
