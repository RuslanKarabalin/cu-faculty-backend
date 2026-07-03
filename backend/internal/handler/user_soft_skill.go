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

type userSoftSkillService interface {
	GetUserSoftSkills(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Skill, int, error)
	AddUserSoftSkill(ctx context.Context, userID uuid.UUID, skillID int) (*model.Skill, error)
	DeleteUserSoftSkill(ctx context.Context, userID uuid.UUID, skillID int) error
}

type UserSoftSkillHandler struct {
	service userSoftSkillService
	logger  *zap.Logger
}

func NewUserSoftSkillHandler(service userSoftSkillService, logger *zap.Logger) *UserSoftSkillHandler {
	return &UserSoftSkillHandler{service: service, logger: logger}
}

func (h *UserSoftSkillHandler) GetUserSoftSkills(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid user id")
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}
	limit, offset := q.Normalize()

	skills, total, err := h.service.GetUserSoftSkills(c.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get user soft skills", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(model.Page[*model.Skill]{Data: skills, Total: total, Limit: limit, Offset: offset})
}

func (h *UserSoftSkillHandler) GetMySoftSkills(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	var q model.PageQuery
	if err := c.Bind().Query(&q); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}
	limit, offset := q.Normalize()

	skills, total, err := h.service.GetUserSoftSkills(c.Context(), cuUser.ID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get user soft skills", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(model.Page[*model.Skill]{Data: skills, Total: total, Limit: limit, Offset: offset})
}

func (h *UserSoftSkillHandler) AddMySoftSkill(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	skillID, err := strconv.Atoi(c.Params("skillId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid skill id")
	}

	skill, err := h.service.AddUserSoftSkill(c.Context(), cuUser.ID, skillID)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidRefID) || errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "soft skill not found")
		}
		h.logger.Error("failed to add user soft skill", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.Status(fiber.StatusCreated).JSON(skill)
}

func (h *UserSoftSkillHandler) DeleteMySoftSkill(c fiber.Ctx) error {
	cuUser, err := currentUser(c, h.logger)
	if err != nil {
		return err
	}

	skillID, err := strconv.Atoi(c.Params("skillId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid skill id")
	}

	if err := h.service.DeleteUserSoftSkill(c.Context(), cuUser.ID, skillID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return respondError(c, fiber.StatusNotFound, "soft skill not assigned")
		}
		h.logger.Error("failed to delete user soft skill", zap.Error(err))
		return respondError(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
