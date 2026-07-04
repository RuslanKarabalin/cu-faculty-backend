package middleware

import (
	"errors"

	"faculty/internal/apierr"
	"faculty/internal/cuclient"

	"github.com/gofiber/fiber/v3"
)

func Auth(client *cuclient.Client, publicPaths map[string]struct{}) fiber.Handler {
	return func(c fiber.Ctx) error {
		if _, ok := publicPaths[c.Path()]; ok {
			return c.Next()
		}
		cookie := c.Cookies(cuclient.CookieName)
		if cookie == "" {
			return apierr.WriteCode(c, fiber.StatusUnauthorized, apierr.CodeUnauthorized, "authentication cookie not provided")
		}
		cuUser, err := client.Authorize(c.Context(), cookie)
		if err != nil {
			if errors.Is(err, cuclient.ErrUnauthorized) {
				return apierr.WriteCode(c, fiber.StatusUnauthorized, apierr.CodeUnauthorized, "authentication cookie rejected by upstream")
			}
			return apierr.WriteCode(c, fiber.StatusBadGateway, apierr.CodeUpstream, "authentication upstream is unavailable")
		}
		SetCuUser(c, cuUser)
		return c.Next()
	}
}
