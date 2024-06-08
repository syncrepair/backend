package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/internal/util"
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

func (h *CompanyHandler) Routes(router *echo.Group) {
	companies := router.Group("/companies")
	{
		companies.POST("", h.Create)
	}
}

type companyCreateRequest struct {
	Name string `json:"name"`
}

func (h *CompanyHandler) Create(ctx echo.Context) error {
	var req companyCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id, err := h.usecase.Create(util.Ctx(ctx), usecase.CompanyCreateRequest{
		Name: req.Name,
	})
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, domain.ErrInternalServer, err)
	}

	return SuccessResponse(ctx, http.StatusOK, id)
}
