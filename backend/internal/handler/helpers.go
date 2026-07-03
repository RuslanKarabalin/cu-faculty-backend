package handler

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"faculty/internal/middleware"
	"faculty/internal/model"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
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

type photoUploader interface {
	Upload(ctx context.Context, key, contentType string, body io.Reader, size int64) error
}

func bindMultipartData(c fiber.Ctx, out any) error {
	data := c.FormValue("data")
	if data == "" {
		return respondError(c, fiber.StatusBadRequest, "data part is required")
	}
	if err := json.Unmarshal([]byte(data), out); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid data part")
	}
	return nil
}

func bindOptionalMultipartData(c fiber.Ctx, out any) (bool, error) {
	data := c.FormValue("data")
	if data == "" {
		return false, nil
	}
	if err := json.Unmarshal([]byte(data), out); err != nil {
		return false, respondError(c, fiber.StatusBadRequest, "invalid data part")
	}
	return true, nil
}

func uploadOptionalPhoto(c fiber.Ctx, storage photoUploader, logger *zap.Logger, keyPrefix string) (string, error) {
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		return "", nil
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", respondError(c, fiber.StatusBadRequest, "photo must be an image")
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.Error("failed to open uploaded photo", zap.Error(err))
		return "", respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	defer func() { _ = file.Close() }()

	key := keyPrefix + uuid.NewString()
	if err := storage.Upload(c.Context(), key, contentType, file, fileHeader.Size); err != nil {
		logger.Error("failed to upload photo", zap.Error(err))
		return "", respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return key, nil
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
