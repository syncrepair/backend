package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"net/http"
)

func (h *Handler) initClientRoutes(router *echo.Group) {
	services := router.Group("/clients", h.authMiddleware())
	{
		services.POST("", h.clientCreate)
		services.DELETE("", h.clientDelete)
	}
}

type clientCreateRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type clientCreateResponse struct {
	ID string `json:"id"`
}

func (h *Handler) clientCreate(ctx echo.Context) error {
	var req clientCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id, err := h.usecases.ClientUsecase.Create(ctx.Request().Context(), usecase.ClientCreateRequest{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		CompanyID:   getCompanyIDFromCtx(ctx),
	})
	if err != nil {
		if errors.Is(err, domain.ErrCompanyNotFound) {
			return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrCompanyNotFound)
		}

		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK, clientCreateResponse{
		ID: id,
	})
}

type clientDeleteRequest struct {
	ID string `json:"id"`
}

func (h *Handler) clientDelete(ctx echo.Context) error {
	var req clientDeleteRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	if err := h.usecases.ClientUsecase.Delete(ctx.Request().Context(), req.ID); err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK)
}
