package handler

import (
	"encoding/json"
	"net/url"

	"faculty/internal/model"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3/client"
)

type UserHandler struct {
	userService *service.UserService
	restClient  *client.Client
	cuBaseURL   string
}

func NewUserHandler(userService *service.UserService, restClient *client.Client, cuBaseURL string) *UserHandler {
	return &UserHandler{
		userService: userService,
		restClient:  restClient,
		cuBaseURL:   cuBaseURL,
	}
}

func (h *UserHandler) Me(c *fiber.Ctx) error {
	cookie := c.Cookies("bff.cookie")
	if cookie == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": "cookie is empty"})
	}

	u, err := url.JoinPath(h.cuBaseURL, "api", "student-hub", "students", "me")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}

	resp, err := h.restClient.SetCookie("bff.cookie", cookie).Get(u)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}

	var cuUserResp model.CuUserResp
	if err := json.Unmarshal(resp.Body(), &cuUserResp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}

	if err := h.userService.CreateUser(c.Context(), cuUserResp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}

	return c.JSON(cuUserResp)
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}
	return c.JSON(users)
}
