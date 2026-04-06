package handler

import (
	"context"
	"errors"

	"faculty/internal/model"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userService interface {
	CreateUser(ctx context.Context, cuUser model.CuUserResp) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error)
}

type UserHandler struct {
	userService userService
	logger      *zap.Logger
}

func NewUserHandler(userService userService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	cuUserResp, ok := c.Locals("cuUser").(*model.CuUserResp)
	if !ok || cuUserResp == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if cuUserResp.Id == (uuid.UUID{}) || cuUserResp.FirstName == "" || cuUserResp.LastName == "" || cuUserResp.BirthDate == "" {
		h.logger.Error("incomplete user data from CU API")
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "incomplete user data from upstream"})
	}

	statusCode := fiber.StatusCreated
	if err := h.userService.CreateUser(c.Context(), *cuUserResp); err != nil {
		if errors.Is(err, service.ErrInvalidBirthDate) {
			h.logger.Error("invalid birth_date from CU API", zap.String("birth_date", cuUserResp.BirthDate))
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "invalid data from upstream"})
		}
		if errors.Is(err, service.ErrAlreadyRegistered) {
			statusCode = fiber.StatusOK
		} else {
			h.logger.Error("failed to create user", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
	}

	user, err := h.userService.GetUserByID(c.Context(), cuUserResp.Id)
	if err != nil {
		h.logger.Error("failed to fetch user after register", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(statusCode).JSON(user)
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 1
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := h.userService.GetAllUsers(c.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to get users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	return c.JSON(model.Page[*model.User]{
		Data:   users,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}
