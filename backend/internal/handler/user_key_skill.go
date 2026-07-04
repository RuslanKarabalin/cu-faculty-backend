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
	GetUserKeySkills(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Skill, int, error)
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

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	skills, total, err := h.service.GetUserKeySkills(c.Context(), userID, limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get user key skills", err)
	}
	return c.JSON(model.Page[*model.Skill]{Data: skills, Total: total, Limit: limit, Offset: offset})
}

func (h *UserKeySkillHandler) GetMyKeySkills(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondBindError(c)
	}
	limit, offset := q.Normalize()

	skills, total, err := h.service.GetUserKeySkills(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		return unexpectedError(c, h.logger, "failed to get user key skills", err)
	}
	return c.JSON(model.Page[*model.Skill]{Data: skills, Total: total, Limit: limit, Offset: offset})
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
		return unexpectedError(c, h.logger, "failed to add user key skill", err)
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
		return unexpectedError(c, h.logger, "failed to delete user key skill", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
