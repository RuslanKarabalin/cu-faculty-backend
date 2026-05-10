package handler

import (
	"context"
	"errors"
	"strings"
	"time"

	"faculty/internal/cuclient"
	"faculty/internal/middleware"
	"faculty/internal/model"
	"faculty/internal/repository"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userService interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error)
}

type eduPlaceService interface {
	GetEduPlacesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EduPlace, error)
}

type registrationService interface {
	Register(ctx context.Context, cuUser model.CuUserResp, eduPlaces []model.CreateEduPlaceParams) (*model.User, bool, error)
}

type UserHandler struct {
	userService         userService
	eduPlaceService     eduPlaceService
	registrationService registrationService
	logger              *zap.Logger
	cuClient            *cuclient.Client
}

func NewUserHandler(
	userService userService,
	eduPlaceService eduPlaceService,
	registrationService registrationService,
	logger *zap.Logger,
	cuClient *cuclient.Client,
) *UserHandler {
	return &UserHandler{
		userService:         userService,
		eduPlaceService:     eduPlaceService,
		registrationService: registrationService,
		logger:              logger,
		cuClient:            cuClient,
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

	existing, err := h.userService.GetUserByID(c.Context(), cuUserResp.ID)
	if err == nil {
		return c.Status(fiber.StatusOK).JSON(existing)
	}
	if !errors.Is(err, repository.ErrNotFound) {
		h.logger.Error("failed to lookup user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	cookie := c.Cookies("bff.cookie")
	cuEduPlaces, err := h.cuClient.StudentEduInfo(c.Context(), cookie)
	if err != nil {
		h.logger.Error("failed to fetch user edu place", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	eduPlaceParams := make([]model.CreateEduPlaceParams, 0, len(cuEduPlaces))
	for _, e := range cuEduPlaces {
		t, err := time.Parse("2006-01-02", e.EducationProgram.StartDate)
		if err != nil {
			h.logger.Error("failed to parse edu place date", zap.Error(err))
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "invalid data from upstream"})
		}
		eduPlaceParams = append(eduPlaceParams, model.CreateEduPlaceParams{
			UniversityId:   207,
			Grade:          strings.ToLower(e.EducationProgram.Level),
			Specialization: e.EducationProgram.Name,
			StartYear:      t.Year(),
			IsStudyingNow:  true,
		})
	}

	user, isNewUser, err := h.registrationService.Register(c.Context(), *cuUserResp, eduPlaceParams)
	if err != nil {
		if errors.Is(err, service.ErrInvalidBirthDate) {
			h.logger.Error("invalid birth_date from CU API", zap.String("birth_date", cuUserResp.BirthDate))
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
	type Query struct {
		Limit  int `query:"limit"`
		Offset int `query:"offset"`
	}

	var q Query
	if err := c.Bind().Query(&q); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	limit := q.Limit
	offset := q.Offset

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
