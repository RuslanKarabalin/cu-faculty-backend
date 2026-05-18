package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os/signal"
	"syscall"
	"time"

	"faculty/internal/config"
	"faculty/internal/cuclient"
	"faculty/internal/db"

	fiberzap "github.com/gofiber/contrib/v3/zap"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

const (
	dbPingTimeout   = 5 * time.Second
	shutdownTimeout = 30 * time.Second
	httpReqTimeout  = 10 * time.Second
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
		return nil, fmt.Errorf("init logger: %w", err)
	}

	dbPool, err := newDBPool(cfg)
	if err != nil {
		return nil, err
	}

	goose.SetLogger(zap.NewStdLog(logger))
	if err := db.RunMigrations(dbPool); err != nil {
		return nil, err
	}

	a := &App{
		Fiber:    newFiber(logger),
		Config:   cfg,
		DB:       dbPool,
		Logger:   logger,
		CuClient: cuclient.New(&http.Client{Timeout: httpReqTimeout}, cfg.CuBaseUrl),
	}
	a.registerRoutes()

	return a, nil
}

func newDBPool(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := buildPgDSN(cfg)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}
	poolConfig.MaxConns = cfg.PgMaxConns
	poolConfig.MinConns = cfg.PgMinConns

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), dbPingTimeout)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return pool, nil
}

func buildPgDSN(cfg *config.Config) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.PgUsername, cfg.PgPassword),
		Host:   net.JoinHostPort(cfg.PgHost, cfg.PgPort),
		Path:   cfg.PgDatabase,
	}
	return u.String()
}

func newFiber(logger *zap.Logger) *fiber.App {
	f := fiber.New()
	f.Use(fiberzap.New(fiberzap.Config{Logger: logger}))
	return f
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		a.Logger.Info("shutting down...")
		_ = a.Fiber.ShutdownWithTimeout(shutdownTimeout)
	}()

	return a.Fiber.Listen(a.Config.Addr, fiber.ListenConfig{DisableStartupMessage: true})
}

func (a *App) Close() {
	a.DB.Close()
	_ = a.Logger.Sync()
}
