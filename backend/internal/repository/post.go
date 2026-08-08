package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreatePost(ctx context.Context, params model.CreatePostParams) (uuid.UUID, error) {
	query := `
	insert into posts(author_id, title, content)
	values($1, $2, $3)
	returning id
	`

	var id uuid.UUID
	if err := r.db.QueryRow(ctx, query, params.AuthorID, params.Title, params.Content).Scan(&id); err != nil {
		return uuid.UUID{}, wrapPgError(err)
	}
	return id, nil
}

func (r *Repository) UpdatePost(ctx context.Context, params model.UpdatePostParams) error {
	query := `
	update posts
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

func (r *Repository) DeletePost(ctx context.Context, id, authorID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from posts where id = $1 and author_id = $2`, id, authorID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetPostByID(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	query := postSelectColumns + `
	where a.id = $1
	`

	a := &model.Post{Author: &model.User{}}
	err := r.db.QueryRow(ctx, query, id).Scan(postScanTargets(a)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get post by id: %w", err)
	}
	return a, nil
}

func (r *Repository) GetVisiblePostByID(ctx context.Context, id, viewerID uuid.UUID) (*model.Post, error) {
	query := postSelectColumns + `
	where a.id = $1
		and (a.is_archived = false or a.author_id = $2)
		and u.deleted_at is null
		and not exists (
			select 1 from blocked_users bu
			where bu.user_id = a.author_id and bu.blocked_user_id = $2
		)
	`

	a := &model.Post{Author: &model.User{}}
	err := r.db.QueryRow(ctx, query, id, viewerID).Scan(postScanTargets(a)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get post by id: %w", err)
	}
	return a, nil
}

func (r *Repository) GetPosts(ctx context.Context, viewerID uuid.UUID, limit, offset int) ([]*model.Post, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `
		select count(*) from posts a
		join users u on u.id = a.author_id
		where a.is_archived = false
			and u.deleted_at is null
			and not exists (
				select 1 from blocked_users bu
				where bu.user_id = a.author_id and bu.blocked_user_id = $1
			)
	`, viewerID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	query := postSelectColumns + `
	where a.is_archived = false
		and u.deleted_at is null
		and not exists (
			select 1 from blocked_users bu
			where bu.user_id = a.author_id and bu.blocked_user_id = $1
		)
	order by a.created_at desc, a.id
	limit $2 offset $3
	`

	return r.selectPosts(ctx, total, query, viewerID, limit, offset)
}

func (r *Repository) GetPostsByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Post, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from posts where author_id = $1`, authorID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	query := postSelectColumns + `
	where a.author_id = $1
	order by a.created_at desc, a.id
	limit $2 offset $3
	`

	return r.selectPosts(ctx, total, query, authorID, limit, offset)
}

func (r *Repository) selectPosts(ctx context.Context, total int, query string, args ...any) ([]*model.Post, int, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to select posts: %w", err)
	}
	defer rows.Close()

	posts := make([]*model.Post, 0)
	for rows.Next() {
		a := &model.Post{Author: &model.User{}}
		if err := rows.Scan(postScanTargets(a)...); err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, a)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return posts, total, nil
}

const postSelectColumns = `
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
	from posts a
	join users u on u.id = a.author_id
	left join statuses st on st.id = u.status_id
	`

func postScanTargets(a *model.Post) []any {
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
