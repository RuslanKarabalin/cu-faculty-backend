package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateAnnouncement(ctx context.Context, params model.CreateAnnouncementParams) (uuid.UUID, error) {
	query := `
	insert into announcements(author_id, title, content)
	values($1, $2, $3)
	returning id
	`

	var id uuid.UUID
	if err := r.db.QueryRow(ctx, query, params.AuthorID, params.Title, params.Content).Scan(&id); err != nil {
		return uuid.UUID{}, wrapPgError(err)
	}
	return id, nil
}

func (r *Repository) UpdateAnnouncement(ctx context.Context, params model.UpdateAnnouncementParams) error {
	query := `
	update announcements
	set title = $3
		, content = $4
		, is_archived = $5
	where id = $1 and author_id = $2
	`

	tag, err := r.db.Exec(ctx, query, params.ID, params.AuthorID, params.Title, params.Content, params.IsArchived)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteAnnouncement(ctx context.Context, id, authorID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from announcements where id = $1 and author_id = $2`, id, authorID)
	if err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetAnnouncementByID(ctx context.Context, id uuid.UUID) (*model.Announcement, error) {
	query := announcementSelectColumns + `
	where a.id = $1
	`

	a := &model.Announcement{Author: &model.User{}}
	err := r.db.QueryRow(ctx, query, id).Scan(announcementScanTargets(a)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get announcement by id: %w", err)
	}
	return a, nil
}

func (r *Repository) GetAnnouncements(ctx context.Context, limit, offset int) ([]*model.Announcement, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from announcements where is_archived = false`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count announcements: %w", err)
	}

	query := announcementSelectColumns + `
	where a.is_archived = false
	order by a.created_at desc, a.id
	limit $1 offset $2
	`

	return r.selectAnnouncements(ctx, total, query, limit, offset)
}

func (r *Repository) GetAnnouncementsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Announcement, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from announcements where author_id = $1`, authorID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count announcements: %w", err)
	}

	query := announcementSelectColumns + `
	where a.author_id = $1
	order by a.created_at desc, a.id
	limit $2 offset $3
	`

	return r.selectAnnouncements(ctx, total, query, authorID, limit, offset)
}

func (r *Repository) selectAnnouncements(ctx context.Context, total int, query string, args ...any) ([]*model.Announcement, int, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to select announcements: %w", err)
	}
	defer rows.Close()

	announcements := make([]*model.Announcement, 0)
	for rows.Next() {
		a := &model.Announcement{Author: &model.User{}}
		if err := rows.Scan(announcementScanTargets(a)...); err != nil {
			return nil, 0, fmt.Errorf("failed to scan announcement: %w", err)
		}
		announcements = append(announcements, a)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return announcements, total, nil
}

const announcementSelectColumns = `
	select
		a.id
		, a.title
		, a.content
		, a.is_archived
		, a.created_at
		, u.id
		, u.photo_s3_key
		, u.first_name
		, u.last_name
		, u.bio
		, u.birth_date
		, u.speciality
		, st.content
		, u.role
	from announcements a
	join users u on u.id = a.author_id
	left join statuses st on st.id = u.status_id
	`

func announcementScanTargets(a *model.Announcement) []any {
	return []any{
		&a.ID,
		&a.Title,
		&a.Content,
		&a.IsArchived,
		&a.CreatedAt,
		&a.Author.ID,
		&a.Author.PhotoS3Key,
		&a.Author.FirstName,
		&a.Author.LastName,
		&a.Author.Bio,
		&a.Author.BirthDate,
		&a.Author.Speciality,
		&a.Author.Status,
		&a.Author.Role,
	}
}
