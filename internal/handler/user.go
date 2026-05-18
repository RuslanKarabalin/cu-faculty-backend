package handler

import (
	"context"
	"errors"

	"faculty/internal/cuclient"
	"faculty/internal/middleware"
	"faculty/internal/model"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userService interface {
	GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error)
}

type eduPlaceService interface {
	GetEduPlacesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EduPlace, error)
}

type registrationService interface {
	Register(ctx context.Context, cuUser model.CuUserResp, cookie string) (*model.User, bool, error)
}

type UserHandler struct {
	userService         userService
	eduPlaceService     eduPlaceService
	registrationService registrationService
	logger              *zap.Logger
}

func NewUserHandler(
	userService userService,
	eduPlaceService eduPlaceService,
	registrationService registrationService,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		userService:         userService,
		eduPlaceService:     eduPlaceService,
		registrationService: registrationService,
		logger:              logger,
	}
}

func (h *UserHandler) Register(c fiber.Ctx) error {
	cuUserResp, ok := middleware.GetCuUser(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if cuUserResp.ID == (uuid.UUID{}) || cuUserResp.FirstName == "" || cuUserResp.LastName == "" || cuUserResp.BirthDate == "" {
		h.logger.Error("incomplete user data from CU API")
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "incomplete user data from upstream"})
	}

	cookie := c.Cookies(cuclient.CookieName)
	user, isNewUser, err := h.registrationService.Register(c.Context(), *cuUserResp, cookie)
	if err != nil {
		if errors.Is(err, service.ErrInvalidBirthDate) || errors.Is(err, service.ErrInvalidUpstreamData) {
			h.logger.Error("invalid data from CU API", zap.Error(err))
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "invalid data from upstream"})
		}
		h.logger.Error("failed to register user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	statusCode := fiber.StatusCreated
	if !isNewUser {
		statusCode = fiber.StatusOK
	}
	return c.Status(statusCode).JSON(user)
}

func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	limit, offset := q.Normalize()

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

func (h *UserHandler) GetUserEduPlaces(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	places, err := h.eduPlaceService.GetEduPlacesByUserID(c.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get edu places", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	return c.JSON(places)
}
