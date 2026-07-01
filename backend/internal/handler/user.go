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
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req model.UpdateUserRequest) (*model.User, error)
	SetPhoto(ctx context.Context, id uuid.UUID, key string) (*model.User, error)
}

type registrationService interface {
	Register(ctx context.Context, cuUser model.CuUserResp, cookie string) (*model.User, bool, error)
}

type photoStorage interface {
	Upload(ctx context.Context, key, contentType string, body io.Reader, size int64) error
	PresignDownload(ctx context.Context, key string) (string, error)
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
		h.logger.Error("failed to register user", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}

	h.attachPhotoURL(c.Context(), user)

	statusCode := fiber.StatusCreated
	if !isNewUser {
		statusCode = fiber.StatusOK
	}
	return c.Status(statusCode).JSON(user)
}

// UploadMyPhoto stores the uploaded image in object storage and wires it to the
// current user's profile in a single request. The client sends the file as
// multipart form-data under the "photo" field.
func (h *UserHandler) UploadMyPhoto(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "photo file is required")
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return respondError(c, fiber.StatusBadRequest, "photo must be an image")
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.logger.Error("failed to open uploaded photo", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	defer func() { _ = file.Close() }()

	key := photoKeyPrefix(cuUser.ID) + uuid.NewString()
	if err := h.storage.Upload(c.Context(), key, contentType, file, fileHeader.Size); err != nil {
		h.logger.Error("failed to upload photo", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}

	user, err := h.userService.SetPhoto(c.Context(), cuUser.ID, key)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "user not found")
		}
		h.logger.Error("failed to set user photo", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	h.attachPhotoURL(c.Context(), user)
	return c.JSON(user)
}

func (h *UserHandler) GetMe(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	user, err := h.userService.GetUserByID(c.Context(), cuUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "user not found")
		}
		h.logger.Error("failed to get current user", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
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
	if err := c.Bind().JSON(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	user, err := h.userService.UpdateUser(c.Context(), cuUser.ID, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "user not found")
		}
		if errors.Is(err, repository.ErrInvalidRefID) {
			return respondError(c, fiber.StatusBadRequest, "invalid status id")
		}
		h.logger.Error("failed to update current user", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	h.attachPhotoURL(c.Context(), user)
	return c.JSON(user)
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
		h.logger.Error("failed to get user by id", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	h.attachPhotoURL(c.Context(), user)
	return c.JSON(user)
}

func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}
	limit, offset := q.Normalize()

	users, total, err := h.userService.GetAllUsers(c.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to get users", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
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
