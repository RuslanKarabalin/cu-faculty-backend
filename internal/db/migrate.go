package db

import (
	"embed"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFs embed.FS

func RunMigrations(conn *pgxpool.Pool) {
	goose.SetBaseFS(migrationsFs)
	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}
	db := stdlib.OpenDBFromPool(conn)
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	if err := db.Close(); err != nil {
		log.Fatalf("failed to close database connection while migrations running: %v", err)
	}
}
