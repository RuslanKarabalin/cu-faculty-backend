package app

import (
	"context"
	"faculty/internal/config"
	"faculty/internal/db"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type App struct {
	Fiber  *fiber.App
	Config *config.Config
	DB     *pgxpool.Pool
	Logger *zap.Logger
}

func New() (*App, error) {
	cfg := config.ReadConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.New(context.Background(), cfg.GetPostgresUrl())
	if err != nil {
		return nil, err
	}

	goose.SetLogger(zap.NewStdLog(logger))
	err = db.RunMigrations(dbPool)
	if err != nil {
		return nil, err
	}

	f := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	f.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	a := &App{
		Fiber:  f,
		Config: cfg,
		DB:     dbPool,
		Logger: logger,
	}

	if err := a.registerRoutes(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	return a.Fiber.Listen(a.Config.Addr)
}
