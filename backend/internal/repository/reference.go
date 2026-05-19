package repository

import (
	"context"
	"faculty/internal/model"
	"fmt"
)

func (r *Repository) GetStatuses(ctx context.Context) ([]*model.Status, error) {
	rows, err := r.db.Query(ctx, `select id, content from statuses order by id`)
	if err != nil {
		return nil, fmt.Errorf("failed to select statuses: %w", err)
	}
	defer rows.Close()

	items := make([]*model.Status, 0)
	for rows.Next() {
		s := &model.Status{}
		if err := rows.Scan(&s.ID, &s.Content); err != nil {
			return nil, fmt.Errorf("failed to scan status: %w", err)
		}
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return items, nil
}

func (r *Repository) GetKeySkills(ctx context.Context) ([]*model.Skill, error) {
	rows, err := r.db.Query(ctx, `select id, name from key_skills order by name, id`)
	if err != nil {
		return nil, fmt.Errorf("failed to select key skills: %w", err)
	}
	defer rows.Close()

	items := make([]*model.Skill, 0)
	for rows.Next() {
		s := &model.Skill{}
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, fmt.Errorf("failed to scan key skill: %w", err)
		}
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return items, nil
}

func (r *Repository) GetSoftSkills(ctx context.Context) ([]*model.Skill, error) {
	rows, err := r.db.Query(ctx, `select id, name from soft_skills order by name, id`)
	if err != nil {
		return nil, fmt.Errorf("failed to select soft skills: %w", err)
	}
	defer rows.Close()

	items := make([]*model.Skill, 0)
	for rows.Next() {
		s := &model.Skill{}
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, fmt.Errorf("failed to scan soft skill: %w", err)
		}
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return items, nil
}

func (r *Repository) GetCompanies(ctx context.Context) ([]*model.Company, error) {
	rows, err := r.db.Query(ctx, `select id, name from companies order by name, id`)
	if err != nil {
		return nil, fmt.Errorf("failed to select companies: %w", err)
	}
	defer rows.Close()

	items := make([]*model.Company, 0)
	for rows.Next() {
		c := &model.Company{}
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return items, nil
}

func (r *Repository) GetWorkPositions(ctx context.Context) ([]*model.WorkPosition, error) {
	rows, err := r.db.Query(ctx, `select id, name from work_positions order by name, id`)
	if err != nil {
		return nil, fmt.Errorf("failed to select work positions: %w", err)
	}
	defer rows.Close()

	items := make([]*model.WorkPosition, 0)
	for rows.Next() {
		p := &model.WorkPosition{}
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, fmt.Errorf("failed to scan work position: %w", err)
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return items, nil
}

func (r *Repository) GetUniversities(ctx context.Context) ([]*model.University, error) {
	rows, err := r.db.Query(ctx, `select id, name, short_name from universities order by short_name, id`)
	if err != nil {
		return nil, fmt.Errorf("failed to select universities: %w", err)
	}
	defer rows.Close()

	items := make([]*model.University, 0)
	for rows.Next() {
		u := &model.University{}
		if err := rows.Scan(&u.ID, &u.Name, &u.ShortName); err != nil {
			return nil, fmt.Errorf("failed to scan university: %w", err)
		}
		items = append(items, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return items, nil
}

func (r *Repository) GetFaqs(ctx context.Context) ([]*model.Faq, error) {
	rows, err := r.db.Query(ctx, `select id, question, answer from faqs order by id`)
	if err != nil {
		return nil, fmt.Errorf("failed to select faqs: %w", err)
	}
	defer rows.Close()

	items := make([]*model.Faq, 0)
	for rows.Next() {
		f := &model.Faq{}
		if err := rows.Scan(&f.ID, &f.Question, &f.Answer); err != nil {
			return nil, fmt.Errorf("failed to scan faq: %w", err)
		}
		items = append(items, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return items, nil
}

func (r *Repository) GetEnumValues(ctx context.Context, typeName string) ([]string, error) {
	rows, err := r.db.Query(ctx, `select enumlabel from pg_enum where enumtypid = $1::regtype order by enumsortorder`, typeName)
	if err != nil {
		return nil, fmt.Errorf("failed to select enum values: %w", err)
	}
	defer rows.Close()

	values := make([]string, 0)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("failed to scan enum value: %w", err)
		}
		values = append(values, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return values, nil
}
