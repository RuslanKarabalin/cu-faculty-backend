package main

import (
	"context"
	"encoding/json"
	"faculty/internal/config"
	"faculty/internal/db"
	"faculty/internal/model"
	"log"
	"net/url"
	"os"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3/client"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't create zap logger %s", err)
		os.Exit(1)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal("Can't sync logger after shutdown", zap.Error(err))
		}
	}()

	sugar := logger.Sugar()

	cfg := config.ReadConfig()

	conn, err := pgxpool.New(ctx, cfg.GetPostgresUrl())
	if err != nil {
		sugar.Error("Cannot connect to PostgreSQL", zap.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close()

	restClient := client.New()

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
			})
		}
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/me", func(c *fiber.Ctx) error {
		ck := &model.Cookie{}
		if err := c.BodyParser(ck); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if ck.Cookie == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": "cookie is empty"})
		}

		url, err := url.JoinPath(cfg.CuBaseUrl, "api", "student-hub", "students", "me")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err})
		}

		resp, err := restClient.
			SetCookie("bff.cookie", ck.Cookie).
			Get(url)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err})
		}

		var data map[string]any
		if err := json.Unmarshal(resp.Body(), &data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err})
		}

		userId, ok := data["id"]
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": "id not found"})
		}

		return c.JSON(fiber.Map{"id": userId})
	})

	sugar.Fatal(app.Listen(cfg.Addr))
}
