package app

import (
	"context"
	"faculty/internal/config"
	"faculty/internal/cuclient"
	"faculty/internal/db"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type App struct {
	Fiber    *fiber.App
	Config   *config.Config
	DB       *pgxpool.Pool
	Logger   *zap.Logger
	CuClient *cuclient.Client
}

func New() (*App, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.GetPostgresUrl())
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dbPool.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
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

	httpClient := &http.Client{Timeout: 10 * time.Second}

	a := &App{
		Fiber:    f,
		Config:   cfg,
		DB:       dbPool,
		Logger:   logger,
		CuClient: cuclient.New(httpClient, cfg.CuBaseUrl),
	}

	if err := a.registerRoutes(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		a.Logger.Info("shutting down...")
		_ = a.Fiber.ShutdownWithTimeout(30 * time.Second)
	}()

	return a.Fiber.Listen(a.Config.Addr)
}

func (a *App) Close() {
	a.DB.Close()
	_ = a.Logger.Sync()
}
