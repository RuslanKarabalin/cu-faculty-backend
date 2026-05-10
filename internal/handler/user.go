package handler

import (
	"context"
	"errors"
	"strings"
	"time"

	"faculty/internal/cuclient"
	"faculty/internal/middleware"
	"faculty/internal/model"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userService interface {
	CreateUser(ctx context.Context, cuUser model.CuUserResp) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error)
}

type eduPlaceService interface {
	CreateEduPlace(ctx context.Context, params model.CreateEduPlaceParams) error
	GetEduPlacesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EduPlace, error)
}

type UserHandler struct {
	userService     userService
	eduPlaceService eduPlaceService
	logger          *zap.Logger
	cuClient        *cuclient.Client
}

func NewUserHandler(userService userService, eduPlaceService eduPlaceService, logger *zap.Logger, cuClient *cuclient.Client) *UserHandler {
	return &UserHandler{
		userService:     userService,
		eduPlaceService: eduPlaceService,
		logger:          logger,
		cuClient:        cuClient,
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

	user, err := h.userService.GetUserByID(c.Context(), cuUserResp.ID)
	if err != nil {
		h.logger.Error("failed to fetch user after register", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	cookie := c.Cookies("bff.cookie")
	eduPlaces, err := h.cuClient.StudentEduInfo(c.Context(), cookie)
	if err != nil {
		h.logger.Error("failed to fetch user edu place", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	for _, eduPlace := range eduPlaces {
		t, err := time.Parse("2006-01-02", eduPlace.EducationProgram.StartDate)
		if err != nil {
			h.logger.Error("failed to parse edu place date", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		params := model.CreateEduPlaceParams{
			UserId:         cuUserResp.ID,
			UniversityId:   207,
			Grade:          strings.ToLower(eduPlace.EducationProgram.Level),
			Specialization: eduPlace.EducationProgram.Name,
			StartYear:      t.Year(),
			IsStudyingNow:  true,
		}

		if err := h.eduPlaceService.CreateEduPlace(c.Context(), params); err != nil {
			h.logger.Error("failed to create edu place", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
	}

	return c.Status(statusCode).JSON(user)
}

func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	type Query struct {
		limit  int `query:"limit"`
		offset int `query:"offset"`
	}

	var q Query
	if err := c.Bind().Query(&q); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	limit := q.limit
	offset := q.offset

	if limit == 0 {
		limit = 20
	}

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
