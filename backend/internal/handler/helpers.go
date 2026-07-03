package handler

import (
	"context"

	"faculty/internal/middleware"
	"faculty/internal/model"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func respondError(c fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{"error": msg})
}

type photoPresigner interface {
	PresignDownload(ctx context.Context, key string) (string, error)
}

type photoDeleter interface {
	Delete(ctx context.Context, key string) error
}

func deleteReplacedPhoto(ctx context.Context, storage photoDeleter, logger *zap.Logger, oldKey *string, newKey string) {
	if oldKey == nil || *oldKey == "" || *oldKey == newKey {
		return
	}
	if err := storage.Delete(ctx, *oldKey); err != nil {
		logger.Warn("failed to delete replaced photo", zap.String("key", *oldKey), zap.Error(err))
	}
}

func presignPhoto(ctx context.Context, storage photoPresigner, logger *zap.Logger, key *string) *string {
	if key == nil {
		return nil
	}
	url, err := storage.PresignDownload(ctx, *key)
	if err != nil {
		logger.Warn("failed to presign photo url", zap.Error(err))
		return nil
	}
	return &url
}

func currentUser(c fiber.Ctx, logger *zap.Logger) (*model.CuUserResp, error) {
	cuUser, ok := middleware.GetCuUser(c)
	if !ok {
		logger.Error("cu user missing from context on authenticated route")
		return nil, respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return cuUser, nil
}
