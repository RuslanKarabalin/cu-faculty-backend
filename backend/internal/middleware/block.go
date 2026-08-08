package middleware

import (
	"context"

	"faculty/internal/apierr"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type blockChecker interface {
	IsBlocked(ctx context.Context, blockerID, blockedUserID uuid.UUID) (bool, error)
	IsUserDeleted(ctx context.Context, id uuid.UUID) (bool, error)
}

// DenyIfBlockedByParam returns 404 when the user identified by param has blocked the caller
// or when that user account is soft-deleted.
func DenyIfBlockedByParam(param string, checker blockChecker, logger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		cuUser, ok := GetCuUser(c)
		if !ok {
			logger.Error("cu user missing from context on blocked-by guard")
			return apierr.NewCode(fiber.StatusInternalServerError, apierr.CodeInternal, "internal server error")
		}

		targetID, err := uuid.Parse(c.Params(param))
		if err != nil {
			return apierr.NewCode(fiber.StatusBadRequest, apierr.CodeBadRequest, "invalid user id")
		}
		if targetID == cuUser.ID {
			return c.Next()
		}

		deleted, err := checker.IsUserDeleted(c.Context(), targetID)
		if err != nil {
			logger.Error("failed to check deleted account", zap.Error(err))
			return apierr.NewCode(fiber.StatusInternalServerError, apierr.CodeInternal, "internal server error")
		}
		if deleted {
			return apierr.NewCode(fiber.StatusNotFound, apierr.CodeNotFound, "resource not found")
		}

		blocked, err := checker.IsBlocked(c.Context(), targetID, cuUser.ID)
		if err != nil {
			logger.Error("failed to check block", zap.Error(err))
			return apierr.NewCode(fiber.StatusInternalServerError, apierr.CodeInternal, "internal server error")
		}
		if blocked {
			return apierr.NewCode(fiber.StatusNotFound, apierr.CodeNotFound, "resource not found")
		}
		return c.Next()
	}
}
