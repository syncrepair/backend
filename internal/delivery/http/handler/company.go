package handler

import (
	"github.com/gofiber/fiber/v2"
)

type CompanyHandler struct{}

func NewCompanyHandler() *CompanyHandler {
	return &CompanyHandler{}
}

func (h *CompanyHandler) Create(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotImplemented).SendString("not implemented")
}
