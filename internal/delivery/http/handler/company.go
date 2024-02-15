package handler

import (
	"github.com/gofiber/fiber/v2"
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
	if err := h.usecase.Create(ctx.Context()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.Status(fiber.StatusNotImplemented).SendString("not implemented")
}
