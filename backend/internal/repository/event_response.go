package repository

import (
	"context"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) AddEventResponse(ctx context.Context, userID, eventID uuid.UUID) error {
	query := `
	insert into event_responses(user_id, event_id)
	values($1, $2)
	on conflict do nothing
	`

	if _, err := r.db.Exec(ctx, query, userID, eventID); err != nil {
		return wrapPgError(err)
	}
	return nil
}

func (r *Repository) DeleteEventResponse(ctx context.Context, userID, eventID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from event_responses where user_id = $1 and event_id = $2`, userID, eventID)
	if err != nil {
		return fmt.Errorf("failed to delete event response: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetEventResponders(ctx context.Context, eventID uuid.UUID) ([]*model.User, error) {
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
	from event_responses er
	join users u on u.id = er.user_id
	left join statuses st on st.id = u.status_id
	where er.event_id = $1
	order by u.last_name, u.first_name, u.id
	`

	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to select event responders: %w", err)
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

func (r *Repository) GetEventsRespondedByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Event, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from event_responses where user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count event responses: %w", err)
	}

	query := eventSelectColumns + `
	join event_responses er on er.event_id = e.id
	where er.user_id = $1
	order by e.starts_at desc, e.id
	limit $2 offset $3
	`

	return r.selectEvents(ctx, total, query, userID, limit, offset)
}
