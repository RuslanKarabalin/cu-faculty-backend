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

// presignPhoto returns a presigned download URL for the given S3 key, or nil if
// there is no key or presigning fails.
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
