package handler

import (
	"faculty/internal/middleware"
	"faculty/internal/model"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func respondError(c fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{"error": msg})
}

func currentUser(c fiber.Ctx, logger *zap.Logger) (*model.CuUserResp, error) {
	cuUser, ok := middleware.GetCuUser(c)
	if !ok {
		logger.Error("cu user missing from context on authenticated route")
		return nil, respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return cuUser, nil
}
