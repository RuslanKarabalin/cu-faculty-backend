package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateUser(ctx context.Context, params model.CreateUserParams) error {
	query := `insert into users(id, first_name, last_name, birth_date, role) values($1, $2, $3, $4, 'user') on conflict (id) do nothing`

	tag, err := r.db.Exec(ctx, query, params.ID, params.FirstName, params.LastName, params.BirthDate)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrDuplicate
	}
	return nil
}

func (r *Repository) UpdateUserPhoto(ctx context.Context, id uuid.UUID, key string) (*string, error) {
	query := `
	with old as (select photo_s3_key from users where id = $1)
	update users set photo_s3_key = $2 where id = $1
	returning (select photo_s3_key from old)
	`
	var oldKey *string
	if err := r.db.QueryRow(ctx, query, id, key).Scan(&oldKey); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, wrapPgError(err)
	}
	return oldKey, nil
}

func (r *Repository) GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from users where role = 'user'`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
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
	from users u
	left join statuses st on st.id = u.status_id
	where u.role = 'user'
	order by u.last_name, u.first_name, u.id
	limit $1 offset $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to select users: %w", err)
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
			return nil, 0, fmt.Errorf("failed to scan users: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return users, total, nil
}

func (r *Repository) SearchUsers(ctx context.Context, userID uuid.UUID, search string, limit, offset int) ([]*model.UserSearchResult, int, error) {
	pattern := "%" + search + "%"

	filter := `
	where u.role = 'user'
		and u.id <> $1
		and (
			u.first_name ilike $2
			or u.last_name ilike $2
			or (u.first_name || ' ' || u.last_name) ilike $2
		)
	`

	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from users u `+filter, userID, pattern).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count user search: %w", err)
	}

	selectQuery := `
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
		, case
			when su.saved_user_id is not null then 0
			when c.contact_id is not null then 1
			else 2
		end as rank
	from users u
	left join statuses st on st.id = u.status_id
	left join saved_users su on su.user_id = $1 and su.saved_user_id = u.id
	left join contacts c on c.user_id = $1 and c.contact_id = u.id
	` + filter + `
	order by rank, u.last_name, u.first_name, u.id
	limit $3 offset $4
	`

	rows, err := r.db.Query(ctx, selectQuery, userID, pattern, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	results := make([]*model.UserSearchResult, 0)
	for rows.Next() {
		u := &model.User{}
		var rank int
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
			&rank,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user search result: %w", err)
		}
		results = append(results, &model.UserSearchResult{User: u, Relation: relationForRank(rank)})
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return results, total, nil
}

func (r *Repository) GetProfileCompletenessData(ctx context.Context, id uuid.UUID) (*model.ProfileCompletenessData, error) {
	query := `
	select
		u.photo_s3_key is not null
		, coalesce(btrim(u.speciality), '') <> ''
		, coalesce(btrim(u.bio), '') <> ''
		, u.birth_date is not null
		, (select count(*) from socials where user_id = u.id)
		, (select count(*) from edu_places where user_id = u.id)
		, (select count(*) from work_places where user_id = u.id)
		, (select count(*) from user_key_skills where user_id = u.id)
		, (select count(*) from user_soft_skills where user_id = u.id)
	from users u
	where u.id = $1
	`

	d := &model.ProfileCompletenessData{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&d.HasPhoto,
		&d.HasSpeciality,
		&d.HasBio,
		&d.HasBirthDate,
		&d.SocialsCount,
		&d.EduPlacesCount,
		&d.WorkPlacesCount,
		&d.KeySkillsCount,
		&d.SoftSkillsCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get profile completeness data: %w", err)
	}
	return d, nil
}

func relationForRank(rank int) model.UserRelation {
	switch rank {
	case 0:
		return model.UserRelationSaved
	case 1:
		return model.UserRelationContact
	default:
		return model.UserRelationOther
	}
}

func (r *Repository) UpdateUser(ctx context.Context, params model.UpdateUserParams) error {
	query := `
	update users
	set bio = $2
		, speciality = $3
		, status_id = $4
	where id = $1
	`

	tag, err := r.db.Exec(ctx, query, params.ID, params.Bio, params.Speciality, params.StatusID)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
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
	from users u
	left join statuses st on st.id = u.status_id
	where u.id = $1
	`

	u := &model.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.PhotoS3Key,
		&u.FirstName,
		&u.LastName,
		&u.Bio,
		&u.BirthDate,
		&u.Speciality,
		&u.Status,
		&u.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return u, nil
}
