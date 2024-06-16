package http

import (
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"net/http"
)

func (h *Handler) initCompanyRoutes(router *echo.Group) {
	companies := router.Group("/companies")
	{
		companies.POST("", h.companyCreate)
	}
}

type companyCreateRequest struct {
	Name string `json:"name"`
}

type companyCreateResponse struct {
	ID string `json:"id"`
}

func (h *Handler) companyCreate(ctx echo.Context) error {
	var req companyCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id, err := h.usecases.CompanyUsecase.Create(ctx.Request().Context(), usecase.CompanyCreateRequest{
		Name: req.Name,
	})
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK, companyCreateResponse{
		ID: id,
	})
}
