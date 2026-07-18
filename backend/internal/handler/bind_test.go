package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"faculty/internal/model"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func testApp(t *testing.T, h fiber.Handler) *fiber.App {
	t.Helper()
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler(zap.NewNop())})
	app.Post("/", h)
	return app
}

func do(t *testing.T, app *fiber.App, req *http.Request) *http.Response {
	t.Helper()
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	return resp
}

func TestBindJSONRunsValidation(t *testing.T) {
	app := testApp(t, func(c fiber.Ctx) error {
		var req model.CreatePostRequest
		if err := bindJSON(c, &req); err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusCreated)
	})

	cases := []struct {
		name string
		body string
		want int
	}{
		{"пустые title и content", `{"title":"","content":""}`, fiber.StatusUnprocessableEntity},
		{"нет полей вообще", `{}`, fiber.StatusUnprocessableEntity},
		{"title длиннее лимита", `{"title":"` + strings.Repeat("x", 100) + `","content":"c"}`, fiber.StatusUnprocessableEntity},
		{"невалидный JSON", `not-json`, fiber.StatusBadRequest},
		{"пустое тело", ``, fiber.StatusBadRequest},
		{"валидное тело", `{"title":"t","content":"c"}`, fiber.StatusCreated},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if got := do(t, app, req).StatusCode; got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func multipartReq(t *testing.T, data string) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if data != "" {
		if err := w.WriteField("data", data); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

func TestBindMultipartRunsValidation(t *testing.T) {
	app := testApp(t, func(c fiber.Ctx) error {
		var req model.CreateNewsRequest
		if err := bindMultipartData(c, &req); err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusCreated)
	})

	cases := []struct {
		name string
		data string
		want int
	}{
		{"пустые поля", `{"title":"","content":""}`, fiber.StatusUnprocessableEntity},
		{"битый JSON в data", `not-json`, fiber.StatusBadRequest},
		{"валидное тело", `{"title":"t","content":"c"}`, fiber.StatusCreated},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := do(t, app, multipartReq(t, tc.data)).StatusCode; got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}

	t.Run("нет поля data вовсе", func(t *testing.T) {
		if got := do(t, app, multipartReq(t, "")).StatusCode; got != fiber.StatusBadRequest {
			t.Errorf("got %d, want 400", got)
		}
	})

	t.Run("JSON вместо multipart", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"title":"t","content":"c"}`))
		req.Header.Set("Content-Type", "application/json")
		if got := do(t, app, req).StatusCode; got != fiber.StatusBadRequest {
			t.Errorf("got %d, want 400", got)
		}
	})
}

func TestCurrentUserReturnsErrorWhenMissing(t *testing.T) {
	app := testApp(t, func(c fiber.Ctx) error {
		user, err := currentUser(c, zap.NewNop())
		if err != nil {
			return err
		}
		if user == nil {
			t.Error("err == nil, но и user == nil: вызывающий код упадёт на разыменовании")
			return c.SendStatus(fiber.StatusOK)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	if got := do(t, app, httptest.NewRequest("POST", "/", nil)).StatusCode; got != fiber.StatusInternalServerError {
		t.Errorf("got %d, want 500", got)
	}
}

func TestCreatePostRequestValidate(t *testing.T) {
	if err := (model.CreatePostRequest{Title: "", Content: ""}).Validate(); err == nil {
		t.Error("пустые title/content должны быть невалидны")
	}
	if err := (model.CreatePostRequest{Title: "t", Content: "c"}).Validate(); err != nil {
		t.Errorf("валидное тело отвергнуто: %v", err)
	}
}
