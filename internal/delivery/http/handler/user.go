package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/syncrepair/backend/internal/model"
	"github.com/syncrepair/backend/internal/usecase"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(usecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) SignUp(ctx *fiber.Ctx) error {
	inp := new(model.UserSignUpInput)

	if err := ctx.BodyParser(inp); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := inp.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := h.usecase.SignUp(ctx.Context(), inp); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.SendStatus(fiber.StatusCreated)
}
