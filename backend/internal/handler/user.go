package handler

import (
	"context"
	"errors"
	"io"
	"strings"

	"faculty/internal/cuclient"
	"faculty/internal/model"
	"faculty/internal/repository"
	"faculty/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userService interface {
	GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error)
	SearchUsers(ctx context.Context, userID uuid.UUID, search string, limit, offset int) ([]*model.UserSearchResult, int, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req model.UpdateUserRequest) (*model.User, error)
	SetPhoto(ctx context.Context, id uuid.UUID, key string) (*model.User, *string, error)
	GetProfileCompleteness(ctx context.Context, id uuid.UUID) (*model.ProfileCompleteness, error)
}

type registrationService interface {
	Register(ctx context.Context, cuUser model.CuUserResp, cookie string) (*model.User, bool, error)
}

type photoStorage interface {
	Upload(ctx context.Context, key, contentType string, body io.Reader, size int64) error
	PresignDownload(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type UserHandler struct {
	userService         userService
	registrationService registrationService
	storage             photoStorage
	logger              *zap.Logger
}

func NewUserHandler(
	userService userService,
	registrationService registrationService,
	storage photoStorage,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		userService:         userService,
		registrationService: registrationService,
		storage:             storage,
		logger:              logger,
	}
}

func photoKeyPrefix(userID uuid.UUID) string {
	return "photos/" + userID.String() + "/"
}

func (h *UserHandler) attachPhotoURL(ctx context.Context, u *model.User) {
	if u == nil || u.PhotoS3Key == nil {
		return
	}
	url, err := h.storage.PresignDownload(ctx, *u.PhotoS3Key)
	if err != nil {
		h.logger.Warn("failed to presign photo url", zap.Error(err))
		return
	}
	u.PhotoURL = &url
}

func (h *UserHandler) Register(c fiber.Ctx) error {
	cuUserResp, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	if cuUserResp.ID == (uuid.UUID{}) || cuUserResp.FirstName == "" || cuUserResp.LastName == "" || cuUserResp.BirthDate == "" {
		h.logger.Error("incomplete user data from CU API")
		return respondError(c, fiber.StatusBadGateway, "incomplete user data from upstream")
	}

	cookie := c.Cookies(cuclient.CookieName)
	user, isNewUser, err := h.registrationService.Register(c.Context(), *cuUserResp, cookie)
	if err != nil {
		if errors.Is(err, service.ErrInvalidBirthDate) || errors.Is(err, service.ErrInvalidUpstreamData) {
			h.logger.Error("invalid data from CU API", zap.Error(err))
			return respondError(c, fiber.StatusBadGateway, "invalid data from upstream")
		}
		return unexpectedError(c, h.logger, "failed to register user", err)
	}

	user, err = h.applyPhoto(c, user)
	if err != nil {
		return err
	}

	h.attachPhotoURL(c.Context(), user)

	statusCode := fiber.StatusCreated
	if !isNewUser {
		statusCode = fiber.StatusOK
	}
	return c.Status(statusCode).JSON(user)
}

func (h *UserHandler) updateUser(c fiber.Ctx, id uuid.UUID, req model.UpdateUserRequest) (*model.User, error) {
	user, err := h.userService.UpdateUser(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, respondError(c, fiber.StatusNotFound, "user not found")
		}
		if errors.Is(err, repository.ErrInvalidRefID) {
			return nil, respondError(c, fiber.StatusBadRequest, "invalid status id")
		}
		return nil, unexpectedError(c, h.logger, "failed to update current user", err)
	}
	return user, nil
}

func (h *UserHandler) getUser(c fiber.Ctx, id uuid.UUID) (*model.User, error) {
	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, respondError(c, fiber.StatusNotFound, "user not found")
		}
		return nil, unexpectedError(c, h.logger, "failed to get current user", err)
	}
	return user, nil
}

func (h *UserHandler) applyPhoto(c fiber.Ctx, user *model.User) (*model.User, error) {
	key, err := uploadOptionalPhoto(c, h.storage, h.logger, photoKeyPrefix(user.ID))
	if err != nil {
		return nil, err
	}
	if key == "" {
		return user, nil
	}

	updated, oldKey, err := h.userService.SetPhoto(c.Context(), user.ID, key)
	if err != nil {
		deletePhoto(c.Context(), h.storage, h.logger, key)
		if errors.Is(err, repository.ErrNotFound) {
			return nil, respondError(c, fiber.StatusNotFound, "user not found")
		}
		return nil, unexpectedError(c, h.logger, "failed to set user photo", err)
	}
	deleteReplacedPhoto(c.Context(), h.storage, h.logger, oldKey, key)
	return updated, nil
}

func (h *UserHandler) GetMe(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	user, err := h.getUser(c, cuUser.ID)
	if err != nil {
		return err
	}
	h.attachPhotoURL(c.Context(), user)
	return c.JSON(user)
}

func (h *UserHandler) UpdateMe(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var req model.UpdateUserRequest
	present, err := bindOptionalMultipartData(c, &req)
	if err != nil {
		return err
	}

	var user *model.User
	if present {
		user, err = h.updateUser(c, cuUser.ID, req)
	} else {
		user, err = h.getUser(c, cuUser.ID)
	}
	if err != nil {
		return err
	}

	user, err = h.applyPhoto(c, user)
	if err != nil {
		return err
	}

	h.attachPhotoURL(c.Context(), user)
	return c.JSON(user)
}

func (h *UserHandler) GetMyCompleteness(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	completeness, err := h.userService.GetProfileCompleteness(c.Context(), cuUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "user not found")
		}
		return unexpectedError(c, h.logger, "failed to get profile completeness", err)
	}
	return c.JSON(completeness)
}

func (h *UserHandler) GetStudentByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "user not found")
		}
		return unexpectedError(c, h.logger, "failed to get user by id", err)
	}
	h.attachPhotoURL(c.Context(), user)
	return c.JSON(user)
}

func (h *UserHandler) SearchUsers(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	search := strings.TrimSpace(c.Query("q"))
	if search == "" {
		return respondError(c, fiber.StatusBadRequest, "query parameter q is required")
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	results, total, err := h.userService.SearchUsers(c.Context(), cuUser.ID, search, limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to search users", err)
	}
	for _, res := range results {
		h.attachPhotoURL(c.Context(), res.User)
	}
	return c.JSON(model.Page[*model.UserSearchResult]{
		Data:   results,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	users, total, err := h.userService.GetAllUsers(c.Context(), limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get users", err)
	}
	for _, u := range users {
		h.attachPhotoURL(c.Context(), u)
	}
	return c.JSON(model.Page[*model.User]{
		Data:   users,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}
