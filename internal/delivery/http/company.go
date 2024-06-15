package http

import (
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"net/http"
)

type CompanyHandler struct {
	usecase usecase.CompanyUsecase
}

func NewCompanyHandler(usecase usecase.CompanyUsecase) *CompanyHandler {
	return &CompanyHandler{
		usecase: usecase,
	}
}

func (h *CompanyHandler) initRoutes(router *echo.Group) {
	companies := router.Group("/companies")
	{
		companies.POST("", h.create)
	}
}

type companyCreateRequest struct {
	Name string `json:"name"`
}

type companyCreateResponse struct {
	ID string `json:"id"`
}

func (h *CompanyHandler) create(ctx echo.Context) error {
	var req companyCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id, err := h.usecase.Create(ctx.Request().Context(), usecase.CompanyCreateRequest{
		Name: req.Name,
	})
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK, companyCreateResponse{
		ID: id,
	})
}
