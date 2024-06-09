package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"net/http"
)

type CompanyController struct {
	usecase usecase.CompanyUsecase
}

func NewCompanyController(usecase usecase.CompanyUsecase) *CompanyController {
	return &CompanyController{
		usecase: usecase,
	}
}

func (h *CompanyController) InitRoutes(router *echo.Group) {
	companies := router.Group("/companies")
	{
		companies.POST("", h.Create)
	}
}

type companyCreateRequest struct {
	Name string `json:"name"`
}

type companyCreateResponse struct {
	ID string `json:"id"`
}

func (h *CompanyController) Create(ctx echo.Context) error {
	var req companyCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id, err := h.usecase.Create(Ctx(ctx), usecase.CompanyCreateRequest{
		Name: req.Name,
	})
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, domain.ErrInternalServer, err)
	}

	return SuccessResponse(ctx, http.StatusOK, companyCreateResponse{
		ID: id,
	})
}
