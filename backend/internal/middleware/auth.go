package middleware

import (
	"errors"

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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bff.cookie not provided"})
		}
		cuUser, err := client.Authorize(c.Context(), cookie)
		if err != nil {
			if errors.Is(err, cuclient.ErrUnauthorized) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bff.cookie rejected by CU"})
			}
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "upstream error"})
		}
		SetCuUser(c, cuUser)
		return c.Next()
	}
}
