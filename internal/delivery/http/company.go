package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"net/http"
)

func (h *Handler) initCompanyRoutes(router *echo.Group) {
	companies := router.Group("/companies")
	{
		companies.POST("", h.companyCreate)

		authenticated := companies.Group("", h.authMiddleware())
		{
			authenticated.GET("", h.companyGet)
			authenticated.PUT("", h.companyUpdate)
			authenticated.DELETE("", h.companyDelete)
		}
	}
}

type companyCreateRequest struct {
	Name     string                 `json:"name"`
	Logo     string                 `json:"logo"`
	Location domain.CompanyLocation `json:"location"`
	Settings domain.CompanySettings `json:"settings"`
}

// @Summary      Create
// @Description  Create new company
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        request body companyCreateRequest true "Request body"
// @Success      201 {object} idResponse
// @Failure      400,500 {object} response
// @Router       /companies [post]
func (h *Handler) companyCreate(ctx echo.Context) error {
	var req companyCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id, err := h.usecases.Company.Create(ctx.Request().Context(), usecase.CompanyCreateInput{
		Name:     req.Name,
		Logo:     req.Logo,
		Location: req.Location,
		Settings: req.Settings,
	})
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusCreated, idResponse{
		ID: id,
	})
}

// @Summary      Get
// @Description  Get company
// @Security     UserAuth
// @Tags         companies
// @Accept       json
// @Produce      json
// @Success      200 {object} domain.Company
// @Failure      400,404,500 {object} response
// @Router       /companies [get]
func (h *Handler) companyGet(ctx echo.Context) error {
	company, err := h.usecases.Company.GetByID(ctx.Request().Context(), getCompanyIDFromCtx(ctx))
	if err != nil {
		if errors.Is(err, domain.ErrCompanyNotFound) {
			return newResponse(ctx, http.StatusNotFound, domain.ErrCompanyNotFound)
		}

		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK, company)
}

type companyUpdateRequest struct {
	Name     string                 `json:"name"`
	Logo     string                 `json:"logo"`
	Location domain.CompanyLocation `json:"location"`
	Settings domain.CompanySettings `json:"settings"`
}

// @Summary      Update
// @Description  Update company
// @Security     UserAuth
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        request body companyUpdateRequest true "Request body"
// @Success      200 {object} response
// @Failure      400,500 {object} response
// @Router       /companies [put]
func (h *Handler) companyUpdate(ctx echo.Context) error {
	var req companyUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	if err := h.usecases.Company.Update(ctx.Request().Context(), getCompanyIDFromCtx(ctx), usecase.CompanyUpdateInput{
		Name:     req.Name,
		Logo:     req.Logo,
		Location: req.Location,
		Settings: req.Settings,
	}); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK)
}

// @Summary      Delete
// @Description  Delete company
// @Security     UserAuth
// @Tags         companies
// @Accept       json
// @Produce      json
// @Success      200 {object} response
// @Failure      500 {object} response
// @Router       /companies [delete]
func (h *Handler) companyDelete(ctx echo.Context) error {
	if err := h.usecases.Company.Delete(ctx.Request().Context(), getCompanyIDFromCtx(ctx)); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK)
}
