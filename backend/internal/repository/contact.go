package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateContact(ctx context.Context, params model.CreateContactParams) error {
	return r.RunInTx(ctx, func(tx *Repository) error {
		forward := `
		insert into contacts(user_id, contact_id, note)
		values($1, $2, $3)
		`
		if _, err := tx.db.Exec(ctx, forward, params.UserID, params.ContactID, params.Note); err != nil {
			return wrapPgError(err)
		}

		reverse := `
		insert into contacts(user_id, contact_id)
		values($1, $2)
		on conflict do nothing
		`
		if _, err := tx.db.Exec(ctx, reverse, params.ContactID, params.UserID); err != nil {
			return wrapPgError(err)
		}
		return nil
	})
}

func (r *Repository) UpdateContact(ctx context.Context, params model.UpdateContactParams) error {
	query := `
	update contacts
	set note = $3
	where user_id = $1 and contact_id = $2
	`

	tag, err := r.db.Exec(ctx, query, params.UserID, params.ContactID, params.Note)
	if err != nil {
		return wrapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteContact(ctx context.Context, userID, contactID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `delete from contacts where user_id = $1 and contact_id = $2`, userID, contactID)
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetContact(ctx context.Context, userID, contactID uuid.UUID) (*model.Contact, error) {
	query := `
	select
		c.note
		, u.id
		, u.photo_s3_key
		, u.first_name
		, u.last_name
		, u.bio
		, u.birth_date
		, u.speciality
		, st.content
		, u.role
	from contacts c
	join users u on u.id = c.contact_id
	left join statuses st on st.id = u.status_id
	where c.user_id = $1 and c.contact_id = $2
	`

	contact := &model.Contact{User: &model.User{}}
	err := r.db.QueryRow(ctx, query, userID, contactID).Scan(
		&contact.Note,
		&contact.User.ID,
		&contact.User.PhotoS3Key,
		&contact.User.FirstName,
		&contact.User.LastName,
		&contact.User.Bio,
		&contact.User.BirthDate,
		&contact.User.Speciality,
		&contact.User.Status,
		&contact.User.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}
	return contact, nil
}

func (r *Repository) GetContactsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Contact, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `select count(*) from contacts where user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count contacts: %w", err)
	}

	query := `
	select
		c.note
		, u.id
		, u.photo_s3_key
		, u.first_name
		, u.last_name
		, u.bio
		, u.birth_date
		, u.speciality
		, st.content
		, u.role
	from contacts c
	join users u on u.id = c.contact_id
	left join statuses st on st.id = u.status_id
	where c.user_id = $1
	order by u.last_name, u.first_name, u.id
	limit $2 offset $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to select contacts: %w", err)
	}
	defer rows.Close()

	contacts := make([]*model.Contact, 0)
	for rows.Next() {
		contact := &model.Contact{User: &model.User{}}
		if err := rows.Scan(
			&contact.Note,
			&contact.User.ID,
			&contact.User.PhotoS3Key,
			&contact.User.FirstName,
			&contact.User.LastName,
			&contact.User.Bio,
			&contact.User.BirthDate,
			&contact.User.Speciality,
			&contact.User.Status,
			&contact.User.Role,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan contact: %w", err)
		}
		contacts = append(contacts, contact)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return contacts, total, nil
}
