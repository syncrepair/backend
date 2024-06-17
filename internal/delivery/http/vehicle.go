package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"net/http"
)

func (h *Handler) initVehicleRoutes(router *echo.Group) {
	vehicles := router.Group("/vehicles", h.authMiddleware())
	{
		vehicles.POST("", h.vehicleCreate)
		vehicles.DELETE("", h.vehicleDelete)
	}
}

type vehicleCreateRequest struct {
	Make        string `json:"make"`
	Model       string `json:"model"`
	Year        uint   `json:"year"`
	VIN         string `json:"vin"`
	PlateNumber string `json:"plate_number"`
	ClientID    string `json:"client_id"`
}

type vehicleCreateResponse struct {
	ID string `json:"id"`
}

func (h *Handler) vehicleCreate(ctx echo.Context) error {
	var req vehicleCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id, err := h.usecases.VehicleUsecase.Create(ctx.Request().Context(), usecase.VehicleCreateRequest{
		Make:        req.Make,
		Model:       req.Model,
		Year:        req.Year,
		VIN:         req.VIN,
		PlateNumber: req.PlateNumber,
		ClientID:    req.ClientID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrClientNotFound)
		}

		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK, vehicleCreateResponse{
		ID: id,
	})
}

type vehicleDeleteRequest struct {
	ID string `json:"id"`
}

func (h *Handler) vehicleDelete(ctx echo.Context) error {
	var req vehicleDeleteRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	if err := h.usecases.VehicleUsecase.Delete(ctx.Request().Context(), req.ID); err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK)
}
