package middleware

import (
	"faculty/internal/model"

	"github.com/gofiber/fiber/v2"
)

const cuUserKey = "cuUser"

func SetCuUser(c *fiber.Ctx, u *model.CuUserResp) {
	c.Locals(cuUserKey, u)
}

func GetCuUser(c *fiber.Ctx) (*model.CuUserResp, bool) {
	u, ok := c.Locals(cuUserKey).(*model.CuUserResp)
	return u, ok && u != nil
}
