package main

import (
	"context"
	"faculty/internal/config"
	"faculty/internal/db"
	"os"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal("Can't sync logger after shutdown", zap.Error(err))
		}
	}()

	sugar := logger.Sugar()

	cfg := config.ReadConfig(sugar)

	conn, err := pgxpool.New(ctx, cfg.GetPostgresUrl())
	if err != nil {
		sugar.Error("Cannot connect to PostgreSQL", zap.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close()

	goose.SetLogger(zap.NewStdLog(logger))

	db.RunMigrations(conn)

	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		},
	)

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		if err := conn.Ping(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
				"db":     "unreachable",
			})
		}
		return c.JSON(fiber.Map{
			"status": "ok",
			"db":     "reachable",
		})
	})

	sugar.Fatal(app.Listen(cfg.Addr))
}
