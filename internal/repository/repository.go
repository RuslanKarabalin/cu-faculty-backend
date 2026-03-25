package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	pgPool *pgxpool.Pool
}

func Init(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pgPool: pool,
	}
}
