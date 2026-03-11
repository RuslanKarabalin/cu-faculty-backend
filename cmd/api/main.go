package main

import (
	"fmt"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal("Can't sync logger after shutdown", zap.Error(err))
		}
	}()

	sugar := logger.Sugar()

	viper.SetDefault("APP_PORT", "8080")

	app := fiber.New()

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	// server := api.NewServer()
	// api.RegisterHandlers(app, server)

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	addr := fmt.Sprintf(":%s", viper.GetString("APP_PORT"))

	sugar.Fatal(app.Listen(addr))
}
