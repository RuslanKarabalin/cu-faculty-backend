package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"faculty/internal/middleware"
	"faculty/internal/model"
	"faculty/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// fakePhotoStorage records the keys it was asked to delete. Its method set
// satisfies photoStorage, newsPhotoStorage and eventPhotoStorage alike.
type fakePhotoStorage struct {
	deleted []string
	err     error
}

func (s *fakePhotoStorage) Upload(context.Context, string, string, io.Reader, int64) error {
	return nil
}

func (s *fakePhotoStorage) PresignDownload(context.Context, string) (string, error) {
	return "", nil
}

func (s *fakePhotoStorage) Delete(_ context.Context, key string) error {
	s.deleted = append(s.deleted, key)
	return s.err
}

// The embedded interfaces leave every method we do not exercise nil, so an
// accidental call panics instead of silently succeeding.
type fakeUserService struct {
	userService
	key *string
	err error
}

func (s *fakeUserService) DeletePhoto(context.Context, uuid.UUID) (*string, error) {
	return s.key, s.err
}

type fakeNewsService struct {
	newsService
	key *string
	err error
}

func (s *fakeNewsService) DeletePhoto(context.Context, uuid.UUID, uuid.UUID) (*string, error) {
	return s.key, s.err
}

type fakeEventService struct {
	eventService
	key *string
	err error
}

func (s *fakeEventService) DeletePhoto(context.Context, uuid.UUID, uuid.UUID) (*string, error) {
	return s.key, s.err
}

// authedApp mounts h on method+path with a CU user already in the context, the
// way the real auth middleware would leave it.
func authedApp(t *testing.T, method, path string, h fiber.Handler) *fiber.App {
	t.Helper()
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler(zap.NewNop())})
	app.Use(func(c fiber.Ctx) error {
		middleware.SetCuUser(c, &model.CuUserResp{ID: uuid.New()})
		return c.Next()
	})
	app.Add([]string{method}, path, h)
	return app
}

// photoDeleteCase is shared by all three entities: the wiring differs, the
// contract does not.
type photoDeleteCase struct {
	name        string
	key         *string
	err         error
	wantStatus  int
	wantDeleted []string
}

func photoDeleteCases() []photoDeleteCase {
	return []photoDeleteCase{
		{
			name:        "фото было — 204 и ключ удалён из хранилища",
			key:         new("news/1/photo.jpg"),
			wantStatus:  fiber.StatusNoContent,
			wantDeleted: []string{"news/1/photo.jpg"},
		},
		{
			name:        "фото не было — 204 и хранилище не трогаем (идемпотентность)",
			key:         nil,
			wantStatus:  fiber.StatusNoContent,
			wantDeleted: nil,
		},
		{
			name:        "сущность не найдена или чужая — 404",
			err:         repository.ErrNotFound,
			wantStatus:  fiber.StatusNotFound,
			wantDeleted: nil,
		},
		{
			name:        "ошибка БД — 500",
			err:         errors.New("boom"),
			wantStatus:  fiber.StatusInternalServerError,
			wantDeleted: nil,
		},
	}
}

func assertDeleted(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("удалено из хранилища %v, ожидалось %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("удалён ключ %q, ожидался %q", got[i], want[i])
		}
	}
}

func TestDeleteMyPhoto(t *testing.T) {
	for _, tc := range photoDeleteCases() {
		t.Run(tc.name, func(t *testing.T) {
			storage := &fakePhotoStorage{}
			h := NewUserHandler(&fakeUserService{key: tc.key, err: tc.err}, nil, storage, zap.NewNop())
			app := authedApp(t, http.MethodDelete, "/api/me/photo", h.DeleteMyPhoto)

			resp := do(t, app, httptest.NewRequest(http.MethodDelete, "/api/me/photo", nil))
			if resp.StatusCode != tc.wantStatus {
				t.Errorf("статус %d, ожидался %d", resp.StatusCode, tc.wantStatus)
			}
			assertDeleted(t, storage.deleted, tc.wantDeleted)
		})
	}
}

func TestDeleteNewsPhoto(t *testing.T) {
	for _, tc := range photoDeleteCases() {
		t.Run(tc.name, func(t *testing.T) {
			storage := &fakePhotoStorage{}
			h := NewNewsHandler(&fakeNewsService{key: tc.key, err: tc.err}, storage, zap.NewNop())
			app := authedApp(t, http.MethodDelete, "/api/news/:id/photo", h.DeleteNewsPhoto)

			target := "/api/news/" + uuid.NewString() + "/photo"
			resp := do(t, app, httptest.NewRequest(http.MethodDelete, target, nil))
			if resp.StatusCode != tc.wantStatus {
				t.Errorf("статус %d, ожидался %d", resp.StatusCode, tc.wantStatus)
			}
			assertDeleted(t, storage.deleted, tc.wantDeleted)
		})
	}
}

func TestDeleteEventPhoto(t *testing.T) {
	for _, tc := range photoDeleteCases() {
		t.Run(tc.name, func(t *testing.T) {
			storage := &fakePhotoStorage{}
			h := NewEventHandler(&fakeEventService{key: tc.key, err: tc.err}, storage, zap.NewNop())
			app := authedApp(t, http.MethodDelete, "/api/events/:id/photo", h.DeleteEventPhoto)

			target := "/api/events/" + uuid.NewString() + "/photo"
			resp := do(t, app, httptest.NewRequest(http.MethodDelete, target, nil))
			if resp.StatusCode != tc.wantStatus {
				t.Errorf("статус %d, ожидался %d", resp.StatusCode, tc.wantStatus)
			}
			assertDeleted(t, storage.deleted, tc.wantDeleted)
		})
	}
}

func TestDeletePhotoRejectsInvalidID(t *testing.T) {
	storage := &fakePhotoStorage{}

	newsHandler := NewNewsHandler(&fakeNewsService{}, storage, zap.NewNop())
	newsApp := authedApp(t, http.MethodDelete, "/api/news/:id/photo", newsHandler.DeleteNewsPhoto)
	if got := do(t, newsApp, httptest.NewRequest(http.MethodDelete, "/api/news/not-a-uuid/photo", nil)).StatusCode; got != fiber.StatusBadRequest {
		t.Errorf("news: статус %d, ожидался 400", got)
	}

	eventHandler := NewEventHandler(&fakeEventService{}, storage, zap.NewNop())
	eventApp := authedApp(t, http.MethodDelete, "/api/events/:id/photo", eventHandler.DeleteEventPhoto)
	if got := do(t, eventApp, httptest.NewRequest(http.MethodDelete, "/api/events/not-a-uuid/photo", nil)).StatusCode; got != fiber.StatusBadRequest {
		t.Errorf("events: статус %d, ожидался 400", got)
	}
}

// Удаление объекта из S3 — best-effort: запись в БД уже зафиксирована, и
// откатывать её из-за недоступного хранилища нельзя. Клиент получает 204.
func TestDeletePhotoIgnoresStorageFailure(t *testing.T) {
	storage := &fakePhotoStorage{err: errors.New("s3 unavailable")}
	h := NewUserHandler(&fakeUserService{key: new("photos/1/x")}, nil, storage, zap.NewNop())
	app := authedApp(t, http.MethodDelete, "/api/me/photo", h.DeleteMyPhoto)

	if got := do(t, app, httptest.NewRequest(http.MethodDelete, "/api/me/photo", nil)).StatusCode; got != fiber.StatusNoContent {
		t.Errorf("статус %d, ожидался 204", got)
	}
	assertDeleted(t, storage.deleted, []string{"photos/1/x"})
}

// Без пользователя в контексте ручка не должна дойти до сервиса.
func TestDeleteMyPhotoRequiresUser(t *testing.T) {
	storage := &fakePhotoStorage{}
	h := NewUserHandler(&fakeUserService{key: new("photos/1/x")}, nil, storage, zap.NewNop())

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler(zap.NewNop())})
	app.Add([]string{http.MethodDelete}, "/api/me/photo", h.DeleteMyPhoto)

	if got := do(t, app, httptest.NewRequest(http.MethodDelete, "/api/me/photo", nil)).StatusCode; got != fiber.StatusInternalServerError {
		t.Errorf("статус %d, ожидался 500", got)
	}
	assertDeleted(t, storage.deleted, nil)
}
