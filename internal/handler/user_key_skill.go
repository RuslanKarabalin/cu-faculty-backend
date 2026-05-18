package handler

import (
	"context"
	"errors"
	"strconv"

	"faculty/internal/model"
	"faculty/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userKeySkillService interface {
	GetUserKeySkills(ctx context.Context, userID uuid.UUID) ([]*model.Skill, error)
	AddUserKeySkill(ctx context.Context, userID uuid.UUID, skillID int) (*model.Skill, error)
	DeleteUserKeySkill(ctx context.Context, userID uuid.UUID, skillID int) error
}

type UserKeySkillHandler struct {
	service userKeySkillService
	logger  *zap.Logger
}

func NewUserKeySkillHandler(service userKeySkillService, logger *zap.Logger) *UserKeySkillHandler {
	return &UserKeySkillHandler{service: service, logger: logger}
}

func (h *UserKeySkillHandler) GetUserKeySkills(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid user id")
	}

	skills, err := h.service.GetUserKeySkills(c.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get user key skills", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(skills)
}

func (h *UserKeySkillHandler) GetMyKeySkills(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	skills, err := h.service.GetUserKeySkills(c.Context(), cuUser.ID)
	if err != nil {
		h.logger.Error("failed to get user key skills", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(skills)
}

func (h *UserKeySkillHandler) AddMyKeySkill(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	skillID, err := strconv.Atoi(c.Params("skillId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid skill id")
	}

	skill, err := h.service.AddUserKeySkill(c.Context(), cuUser.ID, skillID)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidRefID) || errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "key skill not found")
		}
		h.logger.Error("failed to add user key skill", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.Status(fiber.StatusCreated).JSON(skill)
}

func (h *UserKeySkillHandler) DeleteMyKeySkill(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	skillID, err := strconv.Atoi(c.Params("skillId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid skill id")
	}

	if err := h.service.DeleteUserKeySkill(c.Context(), cuUser.ID, skillID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "key skill not assigned")
		}
		h.logger.Error("failed to delete user key skill", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
