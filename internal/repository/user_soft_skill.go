package repository

import (
	"context"
	"errors"
	"faculty/internal/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) AddUserSoftSkill(ctx context.Context, userID uuid.UUID, skillID int) error {
	query := `
	insert into user_soft_skills(user_id, soft_skill_id)
	values($1, $2)
	on conflict do nothing
	`

	if _, err := r.db.Exec(ctx, query, userID, skillID); err != nil {
		return wrapPgError(err)
	}
	return nil
}

func (r *Repository) DeleteUserSoftSkill(ctx context.Context, userID uuid.UUID, skillID int) error {
	tag, err := r.db.Exec(ctx, `delete from user_soft_skills where user_id = $1 and soft_skill_id = $2`, userID, skillID)
	if err != nil {
		return fmt.Errorf("failed to delete user soft skill: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetSoftSkillByID(ctx context.Context, id int) (*model.Skill, error) {
	s := &model.Skill{}
	err := r.db.QueryRow(ctx, `select id, name from soft_skills where id = $1`, id).Scan(&s.ID, &s.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get soft skill by id: %w", err)
	}
	return s, nil
}

func (r *Repository) GetUserSoftSkills(ctx context.Context, userID uuid.UUID) ([]*model.Skill, error) {
	query := `
	select ss.id, ss.name
	from user_soft_skills uss
	join soft_skills ss on ss.id = uss.soft_skill_id
	where uss.user_id = $1
	order by ss.name, ss.id
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select user soft skills: %w", err)
	}
	defer rows.Close()

	skills := make([]*model.Skill, 0)
	for rows.Next() {
		s := &model.Skill{}
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, fmt.Errorf("failed to scan soft skill: %w", err)
		}
		skills = append(skills, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return skills, nil
}
