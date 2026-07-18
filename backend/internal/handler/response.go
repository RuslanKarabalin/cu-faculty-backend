package handler

import (
	"errors"

	"faculty/internal/apierr"
	"faculty/internal/repository"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func respondError(status int, msg string) error {
	return apierr.New(status, msg)
}

func respondValidation(msg string) error {
	return apierr.NewCode(fiber.StatusUnprocessableEntity, apierr.CodeValidation, msg)
}

func respondBindError() error {
	return apierr.NewCode(fiber.StatusBadRequest, apierr.CodeBadRequest, "invalid request")
}

func unexpectedError(logger *zap.Logger, msg string, err error) error {
	switch {
	case errors.Is(err, repository.ErrValidation):
		return apierr.NewCode(fiber.StatusUnprocessableEntity, apierr.CodeValidation, "one or more fields are invalid")
	case errors.Is(err, repository.ErrInvalidRefID):
		return apierr.NewCode(fiber.StatusBadRequest, apierr.CodeInvalidReference, "a referenced entity does not exist")
	case errors.Is(err, repository.ErrDuplicate):
		return apierr.NewCode(fiber.StatusConflict, apierr.CodeAlreadyExists, "resource already exists")
	}
	logger.Error(msg, zap.Error(err))
	return apierr.NewCode(fiber.StatusInternalServerError, apierr.CodeInternal, "internal server error")
}

func ErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		if ae, ok := errors.AsType[*apierr.Error](err); ok {
			return apierr.WriteCode(c, ae.Status, ae.Code, ae.Message)
		}

		switch {
		case errors.Is(err, repository.ErrValidation):
			return apierr.WriteCode(c, fiber.StatusUnprocessableEntity, apierr.CodeValidation, "one or more fields are invalid")
		case errors.Is(err, repository.ErrNotFound):
			return apierr.WriteCode(c, fiber.StatusNotFound, apierr.CodeNotFound, "resource not found")
		}

		if fe, ok := errors.AsType[*fiber.Error](err); ok {
			return apierr.WriteCode(c, fe.Code, apierr.CodeForStatus(fe.Code), fe.Message)
		}

		logger.Error("unhandled error", zap.Error(err))
		return apierr.WriteCode(c, fiber.StatusInternalServerError, apierr.CodeInternal, "internal server error")
	}
}
