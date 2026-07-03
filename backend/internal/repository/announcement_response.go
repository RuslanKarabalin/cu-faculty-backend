package repository

import (
	"context"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) AddAnnouncementResponse(ctx context.Context, userID, announcementID uuid.UUID) error {
	query := `
	insert into announcement_responses(user_id, announcement_id)
	select $1, $2
	where exists (
		select 1 from announcements where id = $2 and not is_archived
	)
	on conflict do nothing
	`

	tag, err := r.db.Exec(ctx, query, userID, announcementID)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		// No row inserted: either the announcement is missing/archived, or the
		// user has already responded. Only the former is an error.
		var alreadyResponded bool
		if err := r.db.QueryRow(ctx,
			`select exists(select 1 from announcement_responses where user_id = $1 and announcement_id = $2)`,
			userID, announcementID,
		).Scan(&alreadyResponded); err != nil {
			return fmt.Errorf("failed to check announcement response: %w", err)
		}
		if !alreadyResponded {
			return ErrNotFound
		}
	}
	return nil
}

func (r *Repository) DeleteAnnouncementResponse(ctx context.Context, userID, announcementID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from announcement_responses where user_id = $1 and announcement_id = $2`, userID, announcementID)
	if err != nil {
		return fmt.Errorf("failed to delete announcement response: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetAnnouncementResponders(ctx context.Context, announcementID uuid.UUID) ([]*model.User, error) {
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
	from announcement_responses ar
	join users u on u.id = ar.user_id
	left join statuses st on st.id = u.status_id
	where ar.announcement_id = $1
	order by u.last_name, u.first_name, u.id
	`

	rows, err := r.db.Query(ctx, query, announcementID)
	if err != nil {
		return nil, fmt.Errorf("failed to select announcement responders: %w", err)
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
			return nil, fmt.Errorf("failed to scan responder: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return users, nil
}

func (r *Repository) GetAnnouncementsRespondedByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Announcement, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from announcement_responses where user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count announcement responses: %w", err)
	}

	query := announcementSelectColumns + `
	join announcement_responses ar on ar.announcement_id = a.id
	where ar.user_id = $1
	order by a.created_at desc, a.id
	limit $2 offset $3
	`

	return r.selectAnnouncements(ctx, total, query, userID, limit, offset)
}
