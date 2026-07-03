package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateEvent(ctx context.Context, params model.CreateEventParams) (uuid.UUID, error) {
	query := `
	insert into events(author_id, title, content, place, category, starts_at, registration_link, is_draft)
	values($1, $2, $3, $4, $5, $6, $7, $8)
	returning id
	`

	var id uuid.UUID
	if err := r.db.QueryRow(ctx, query, params.AuthorID, params.Title, params.Content, params.Place, params.Category, params.StartsAt, params.RegistrationLink, params.IsDraft).Scan(&id); err != nil {
		return uuid.UUID{}, wrapPgError(err)
	}
	return id, nil
}

func (r *Repository) UpdateEvent(ctx context.Context, params model.UpdateEventParams) error {
	query := `
	update events
	set title = $3
		, content = $4
		, place = $5
		, category = $6
		, starts_at = $7
		, registration_link = $8
		, is_draft = $9
	where id = $1 and author_id = $2
	`

	tag, err := r.db.Exec(ctx, query, params.ID, params.AuthorID, params.Title, params.Content, params.Place, params.Category, params.StartsAt, params.RegistrationLink, params.IsDraft)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// UpdateEventPhoto sets the event's photo key and returns the previous key (if
// any) so the caller can clean up the replaced object.
func (r *Repository) UpdateEventPhoto(ctx context.Context, id, authorID uuid.UUID, key string) (*string, error) {
	query := `
	with old as (select photo_s3_key from events where id = $1 and author_id = $2)
	update events set photo_s3_key = $3 where id = $1 and author_id = $2
	returning (select photo_s3_key from old)
	`
	var oldKey *string
	if err := r.db.QueryRow(ctx, query, id, authorID, key).Scan(&oldKey); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, wrapPgError(err)
	}
	return oldKey, nil
}

func (r *Repository) DeleteEvent(ctx context.Context, id, authorID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from events where id = $1 and author_id = $2`, id, authorID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetEventByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	query := eventSelectColumns + `
	where e.id = $1
	`

	e := &model.Event{Author: &model.User{}}
	err := r.db.QueryRow(ctx, query, id).Scan(eventScanTargets(e)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get event by id: %w", err)
	}
	return e, nil
}

func (r *Repository) GetEvents(ctx context.Context, limit, offset int) ([]*model.Event, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from events where is_draft = false`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	query := eventSelectColumns + `
	where e.is_draft = false
	order by e.starts_at desc, e.id
	limit $1 offset $2
	`

	return r.selectEvents(ctx, total, query, limit, offset)
}

func (r *Repository) GetEventsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Event, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from events where author_id = $1`, authorID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	query := eventSelectColumns + `
	where e.author_id = $1
	order by e.starts_at desc, e.id
	limit $2 offset $3
	`

	return r.selectEvents(ctx, total, query, authorID, limit, offset)
}

func (r *Repository) selectEvents(ctx context.Context, total int, query string, args ...any) ([]*model.Event, int, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to select events: %w", err)
	}
	defer rows.Close()

	events := make([]*model.Event, 0)
	for rows.Next() {
		e := &model.Event{Author: &model.User{}}
		if err := rows.Scan(eventScanTargets(e)...); err != nil {
			return nil, 0, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return events, total, nil
}

const eventSelectColumns = `
	select
		e.id
		, e.title
		, e.content
		, e.photo_s3_key
		, e.place
		, e.category
		, e.starts_at
		, e.registration_link
		, e.is_draft
		, e.created_at
		, u.id
		, u.photo_s3_key
		, u.first_name
		, u.last_name
		, u.bio
		, u.birth_date
		, u.speciality
		, st.content
		, u.role
	from events e
	join users u on u.id = e.author_id
	left join statuses st on st.id = u.status_id
	`

func eventScanTargets(e *model.Event) []any {
	return []any{
		&e.ID,
		&e.Title,
		&e.Content,
		&e.PhotoS3Key,
		&e.Place,
		&e.Category,
		&e.StartsAt,
		&e.RegistrationLink,
		&e.IsDraft,
		&e.CreatedAt,
		&e.Author.ID,
		&e.Author.PhotoS3Key,
		&e.Author.FirstName,
		&e.Author.LastName,
		&e.Author.Bio,
		&e.Author.BirthDate,
		&e.Author.Speciality,
		&e.Author.Status,
		&e.Author.Role,
	}
}
