package handler

import (
	"errors"

	"faculty/internal/apierr"
	"faculty/internal/repository"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func respondError(c fiber.Ctx, status int, msg string) error {
	return apierr.Write(c, status, msg)
}

func respondValidation(c fiber.Ctx, msg string) error {
	return apierr.WriteCode(c, fiber.StatusUnprocessableEntity, apierr.CodeValidation, msg)
}

func respondBindError(c fiber.Ctx) error {
	return apierr.WriteCode(c, fiber.StatusBadRequest, apierr.CodeBadRequest, "invalid request")
}

func unexpectedError(c fiber.Ctx, logger *zap.Logger, msg string, err error) error {
	switch {
	case errors.Is(err, repository.ErrValidation):
		return apierr.WriteCode(c, fiber.StatusUnprocessableEntity, apierr.CodeValidation, "one or more fields are invalid")
	case errors.Is(err, repository.ErrInvalidRefID):
		return apierr.WriteCode(c, fiber.StatusBadRequest, apierr.CodeInvalidReference, "a referenced entity does not exist")
	case errors.Is(err, repository.ErrDuplicate):
		return apierr.WriteCode(c, fiber.StatusConflict, apierr.CodeAlreadyExists, "resource already exists")
	}
	logger.Error(msg, zap.Error(err))
	return apierr.WriteCode(c, fiber.StatusInternalServerError, apierr.CodeInternal, "internal server error")
}

func ErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
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
