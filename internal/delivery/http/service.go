package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"net/http"
)

func (h *Handler) initServiceRoutes(router *echo.Group) {
	services := router.Group("/services", h.authMiddleware())
	{
		services.POST("", h.serviceCreate)
	}
}

type serviceCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type serviceCreateResponse struct {
	ID string `json:"id"`
}

func (h *Handler) serviceCreate(ctx echo.Context) error {
	var req serviceCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id, err := h.usecases.ServiceUsecase.Create(ctx.Request().Context(), usecase.ServiceCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CompanyID:   getCompanyIDFromCtx(ctx),
	})
	if err != nil {
		if errors.Is(err, domain.ErrCompanyNotFound) {
			return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrCompanyNotFound)
		}

		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK, serviceCreateResponse{
		ID: id,
	})
}
