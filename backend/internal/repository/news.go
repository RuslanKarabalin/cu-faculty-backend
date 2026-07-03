package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateNews(ctx context.Context, params model.CreateNewsParams) (uuid.UUID, error) {
	query := `
	insert into news(author_id, title, content, publish_days, is_draft)
	values($1, $2, $3, $4, $5)
	returning id
	`

	var id uuid.UUID
	if err := r.db.QueryRow(ctx, query, params.AuthorID, params.Title, params.Content, params.PublishDays, params.IsDraft).Scan(&id); err != nil {
		return uuid.UUID{}, wrapPgError(err)
	}
	return id, nil
}

func (r *Repository) UpdateNews(ctx context.Context, params model.UpdateNewsParams) error {
	query := `
	update news
	set title = $3
		, content = $4
		, publish_days = $5
		, is_draft = $6
	where id = $1 and author_id = $2
	`

	tag, err := r.db.Exec(ctx, query, params.ID, params.AuthorID, params.Title, params.Content, params.PublishDays, params.IsDraft)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) UpdateNewsPhoto(ctx context.Context, id, authorID uuid.UUID, key string) (*string, error) {
	query := `
	with old as (select photo_s3_key from news where id = $1 and author_id = $2)
	update news set photo_s3_key = $3 where id = $1 and author_id = $2
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

func (r *Repository) DeleteNews(ctx context.Context, id, authorID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from news where id = $1 and author_id = $2`, id, authorID)
	if err != nil {
		return fmt.Errorf("failed to delete news: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetNewsByID(ctx context.Context, id uuid.UUID) (*model.News, error) {
	query := newsSelectColumns + `
	where n.id = $1
	`

	n := &model.News{Author: &model.User{}}
	err := r.db.QueryRow(ctx, query, id).Scan(newsScanTargets(n)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get news by id: %w", err)
	}
	return n, nil
}

func (r *Repository) GetNews(ctx context.Context, limit, offset int) ([]*model.News, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from news where is_draft = false`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count news: %w", err)
	}

	query := newsSelectColumns + `
	where n.is_draft = false
	order by n.created_at desc, n.id
	limit $1 offset $2
	`

	return r.selectNews(ctx, total, query, limit, offset)
}

func (r *Repository) GetNewsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.News, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from news where author_id = $1`, authorID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count news: %w", err)
	}

	query := newsSelectColumns + `
	where n.author_id = $1
	order by n.created_at desc, n.id
	limit $2 offset $3
	`

	return r.selectNews(ctx, total, query, authorID, limit, offset)
}

func (r *Repository) selectNews(ctx context.Context, total int, query string, args ...any) ([]*model.News, int, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to select news: %w", err)
	}
	defer rows.Close()

	news := make([]*model.News, 0)
	for rows.Next() {
		n := &model.News{Author: &model.User{}}
		if err := rows.Scan(newsScanTargets(n)...); err != nil {
			return nil, 0, fmt.Errorf("failed to scan news: %w", err)
		}
		news = append(news, n)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return news, total, nil
}

const newsSelectColumns = `
	select
		n.id
		, n.title
		, n.content
		, n.photo_s3_key
		, n.publish_days
		, n.is_draft
		, n.created_at
		, u.id
		, u.photo_s3_key
		, u.first_name
		, u.last_name
		, u.bio
		, u.birth_date
		, u.speciality
		, st.content
		, u.role
	from news n
	join users u on u.id = n.author_id
	left join statuses st on st.id = u.status_id
	`

func newsScanTargets(n *model.News) []any {
	return []any{
		&n.ID,
		&n.Title,
		&n.Content,
		&n.PhotoS3Key,
		&n.PublishDays,
		&n.IsDraft,
		&n.CreatedAt,
		&n.Author.ID,
		&n.Author.PhotoS3Key,
		&n.Author.FirstName,
		&n.Author.LastName,
		&n.Author.Bio,
		&n.Author.BirthDate,
		&n.Author.Speciality,
		&n.Author.Status,
		&n.Author.Role,
	}
}
