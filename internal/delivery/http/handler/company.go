package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/syncrepair/backend/internal/model"
	"github.com/syncrepair/backend/internal/usecase"
)

type CompanyHandler struct {
	usecase *usecase.CompanyUsecase
}

func NewCompanyHandler(usecase *usecase.CompanyUsecase) *CompanyHandler {
	return &CompanyHandler{
		usecase: usecase,
	}
}

func (h *CompanyHandler) Create(ctx *fiber.Ctx) error {
	inp := new(model.CompanyCreateInput)

	if err := ctx.BodyParser(inp); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := inp.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := h.usecase.Create(ctx.Context(), inp); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.SendStatus(fiber.StatusCreated)
}
