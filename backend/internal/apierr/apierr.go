package apierr

import "github.com/gofiber/fiber/v3"

const (
	CodeBadRequest       = "bad_request"
	CodeValidation       = "validation_error"
	CodeUnauthorized     = "unauthorized"
	CodeForbidden        = "forbidden"
	CodeNotFound         = "not_found"
	CodeConflict         = "conflict"
	CodeAlreadyExists    = "already_exists"
	CodeInvalidReference = "invalid_reference"
	CodePayloadTooLarge  = "payload_too_large"
	CodeUpstream         = "upstream_error"
	CodeUnavailable      = "unavailable"
	CodeInternal         = "internal_error"
)

type Response struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func CodeForStatus(status int) string {
	switch status {
	case fiber.StatusBadRequest:
		return CodeBadRequest
	case fiber.StatusUnauthorized:
		return CodeUnauthorized
	case fiber.StatusForbidden:
		return CodeForbidden
	case fiber.StatusNotFound:
		return CodeNotFound
	case fiber.StatusConflict:
		return CodeConflict
	case fiber.StatusRequestEntityTooLarge:
		return CodePayloadTooLarge
	case fiber.StatusUnprocessableEntity:
		return CodeValidation
	case fiber.StatusBadGateway:
		return CodeUpstream
	case fiber.StatusServiceUnavailable:
		return CodeUnavailable
	default:
		return CodeInternal
	}
}

func Write(c fiber.Ctx, status int, msg string) error {
	return WriteCode(c, status, CodeForStatus(status), msg)
}

func WriteCode(c fiber.Ctx, status int, code, msg string) error {
	return c.Status(status).JSON(Response{Error: msg, Code: code})
}
