package middleware

import (
	"context"

	"faculty/internal/apierr"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const CodeAccountDeleted = "account_deleted"

type deletedAccountChecker interface {
	IsUserDeleted(ctx context.Context, id uuid.UUID) (bool, error)
}

// DenyDeletedAccount rejects requests from soft-deleted faculty accounts.
// Unregistered CU users (no faculty row) are allowed through so they can register.
func DenyDeletedAccount(checker deletedAccountChecker, logger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		cuUser, ok := GetCuUser(c)
		if !ok {
			return c.Next()
		}

		deleted, err := checker.IsUserDeleted(c.Context(), cuUser.ID)
		if err != nil {
			logger.Error("failed to check deleted account", zap.Error(err))
			return apierr.NewCode(fiber.StatusInternalServerError, apierr.CodeInternal, "internal server error")
		}
		if deleted {
			return apierr.NewCode(fiber.StatusForbidden, CodeAccountDeleted, "account is deleted")
		}
		return c.Next()
	}
}
