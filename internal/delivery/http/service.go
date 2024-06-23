package http

import (
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"net/http"
)

func (h *Handler) initServiceRoutes(router *echo.Group) {
	services := router.Group("/services", h.authMiddleware())
	{
		services.POST("", h.serviceCreate)
		services.GET("", h.serviceGetAll)
		services.GET("/:id", h.serviceGetByID)
		services.PUT("/:id", h.serviceUpdate)
		services.DELETE("/:id", h.serviceDelete)
	}
}

type serviceCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Duration    uint    `json:"duration"`
	Price       float64 `json:"price"`
}

// @Summary      Create
// @Description  Create service
// @Security     UserAuth
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        request body serviceCreateRequest true "Request body"
// @Success      201 {object} idResponse
// @Failure      400,500 {object} response
// @Router       /services [post]
func (h *Handler) serviceCreate(ctx echo.Context) error {
	var req serviceCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id, err := h.usecases.Service.Create(ctx.Request().Context(), usecase.ServiceCreateInput{
		Name:        req.Name,
		Description: req.Description,
		Duration:    req.Duration,
		Price:       req.Price,
		CompanyID:   getCompanyIDFromCtx(ctx),
	})
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusCreated, idResponse{
		ID: id,
	})
}

// @Summary      Get all
// @Description  Get all services
// @Security     UserAuth
// @Tags         services
// @Accept       json
// @Produce      json
// @Success      200 {object} []domain.Service
// @Failure      500 {object} response
// @Router       /services [get]
func (h *Handler) serviceGetAll(ctx echo.Context) error {
	services, err := h.usecases.Service.GetAll(ctx.Request().Context(), getCompanyIDFromCtx(ctx))
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK, services)
}

// @Summary      Get by ID
// @Description  Get service by ID
// @Security     UserAuth
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        id path string true "Service ID"
// @Success      200 {object} domain.Service
// @Failure      500 {object} response
// @Router       /services/{id} [get]
func (h *Handler) serviceGetByID(ctx echo.Context) error {
	id := ctx.Param("id")

	service, err := h.usecases.Service.GetByID(ctx.Request().Context(), getCompanyIDFromCtx(ctx), id)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK, service)
}

type serviceUpdateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Duration    uint    `json:"duration"`
	Price       float64 `json:"price"`
}

// @Summary      Update
// @Description  Update service by ID
// @Security     UserAuth
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        id path string true "Service ID"
// @Param        request body serviceUpdateRequest true "Request body"
// @Success      200 {object} response
// @Failure      400,500 {object} response
// @Router       /services [put]
func (h *Handler) serviceUpdate(ctx echo.Context) error {
	var req serviceUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id := ctx.Param("id")

	if err := h.usecases.Service.Update(ctx.Request().Context(), id, usecase.ServiceUpdateInput{
		Name:        req.Name,
		Description: req.Description,
		Duration:    req.Duration,
		Price:       req.Price,
		CompanyID:   getCompanyIDFromCtx(ctx),
	}); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK)
}

// @Summary      Delete
// @Description  Delete service by ID
// @Security     UserAuth
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        id path string true "Service ID"
// @Success      200 {object} response
// @Failure      500 {object} response
// @Router       /services/{id} [delete]
func (h *Handler) serviceDelete(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := h.usecases.Service.Delete(ctx.Request().Context(), id, getCompanyIDFromCtx(ctx)); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK)
}
