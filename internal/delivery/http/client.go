package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"net/http"
)

func (h *Handler) initClientRoutes(router *echo.Group) {
	clients := router.Group("/clients", h.authMiddleware())
	{
		clients.POST("", h.clientCreate)
		clients.GET("", h.clientGetAll)
		clients.GET("/:id", h.clientGetByID)
		clients.PUT("/:id", h.clientUpdate)
		clients.DELETE("/:id", h.clientDelete)

		vehicles := clients.Group("/:clientID/vehicles", h.authMiddleware())
		{
			vehicles.POST("", h.clientVehicleCreate)
			vehicles.GET("", h.clientVehicleGetAll)
			vehicles.GET("/:id", h.clientVehicleGetByID)
			vehicles.PUT("/:id", h.clientVehicleUpdate)
			vehicles.DELETE("/:id", h.clientVehicleDelete)
		}
	}
}

type clientCreateRequest struct {
	Name        string                 `json:"name"`
	PhoneNumber string                 `json:"phone_number"`
	Vehicles    []domain.ClientVehicle `json:"vehicles"`
	Settings    domain.ClientSettings  `json:"settings"`
}

// @Summary      Create
// @Description  Create client
// @Security     UserAuth
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        request body clientCreateRequest true "Request body"
// @Success      201 {object} idResponse
// @Failure      400,500 {object} response
// @Router       /clients [post]
func (h *Handler) clientCreate(ctx echo.Context) error {
	var req clientCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id, err := h.usecases.Client.Create(ctx.Request().Context(), usecase.ClientCreateInput{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Vehicles:    req.Vehicles,
		Settings:    req.Settings,
		CompanyID:   getCompanyIDFromCtx(ctx),
	})
	if err != nil {
		if errors.Is(err, domain.ErrClientAlreadyExists) {
			return newResponse(ctx, http.StatusBadRequest, domain.ErrClientAlreadyExists)
		}

		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusCreated, idResponse{
		ID: id,
	})
}

// @Summary      Get all
// @Description  Get all clients
// @Security     UserAuth
// @Tags         clients
// @Accept       json
// @Produce      json
// @Success      200 {object} []domain.Client
// @Failure      404,500 {object} response
// @Router       /clients [get]
func (h *Handler) clientGetAll(ctx echo.Context) error {
	clients, err := h.usecases.Client.GetAll(ctx.Request().Context(), getCompanyIDFromCtx(ctx))
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK, clients)
}

// @Summary      Get by ID
// @Description  Get client by ID
// @Security     UserAuth
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        id path string true "Client ID"
// @Success      200 {object} domain.Client
// @Failure      404,500 {object} response
// @Router       /clients/{id} [get]
func (h *Handler) clientGetByID(ctx echo.Context) error {
	id := ctx.Param("id")
	client, err := h.usecases.Client.GetByID(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			return newResponse(ctx, http.StatusNotFound, domain.ErrClientNotFound)
		}

		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK, client)
}

type clientUpdateRequest struct {
	Name        string                 `json:"name"`
	PhoneNumber string                 `json:"phone_number"`
	Vehicles    []domain.ClientVehicle `json:"vehicles"`
	Settings    domain.ClientSettings  `json:"settings"`
}

// @Summary      Update
// @Description  Update client by ID
// @Security     UserAuth
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        id path string true "Client ID"
// @Param        request body clientUpdateRequest true "Request body"
// @Success      200 {object} response
// @Failure      404,500 {object} response
// @Router       /clients/{id} [put]
func (h *Handler) clientUpdate(ctx echo.Context) error {
	var req clientUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	id := ctx.Param("id")

	if err := h.usecases.Client.Update(ctx.Request().Context(), id, usecase.ClientUpdateInput{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Vehicles:    req.Vehicles,
		Settings:    req.Settings,
		CompanyID:   getCompanyIDFromCtx(ctx),
	}); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK)
}

// @Summary      Delete
// @Description  Delete client by ID
// @Security     UserAuth
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        id path string true "Client ID"
// @Success      200 {object} response
// @Failure      500 {object} response
// @Router       /clients/{id} [delete]
func (h *Handler) clientDelete(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := h.usecases.Client.Delete(ctx.Request().Context(), id); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK)
}

type clientVehicleCreateRequest struct {
	Make        string `json:"make"`
	Model       string `json:"model"`
	Year        uint   `json:"year"`
	VIN         string `json:"vin"`
	Distance    uint   `json:"distance"`
	PlateNumber string `json:"plate_number"`
}

// @Summary      Create vehicle
// @Description  Create vehicle for client
// @Security     UserAuth
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        clientID path string true "Client ID"
// @Param        request body clientVehicleCreateRequest true "Request body"
// @Success      201 {object} idResponse
// @Failure      404,500 {object} response
// @Router       /clients/{clientID}/vehicles [post]
func (h *Handler) clientVehicleCreate(ctx echo.Context) error {
	var req clientVehicleCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	clientID := ctx.Param("clientID")

	vehicleID, err := h.usecases.Vehicle.Create(ctx.Request().Context(), clientID, usecase.VehicleCreateInput{
		Make:        req.Make,
		Model:       req.Model,
		Year:        req.Year,
		VIN:         req.VIN,
		Distance:    req.Distance,
		PlateNumber: req.PlateNumber,
	})
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusCreated, idResponse{
		ID: vehicleID,
	})
}

// @Summary      Get all
// @Description  Get all client's vehicles
// @Security     UserAuth
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        clientID path string true "Client ID"
// @Success      200 {object} []domain.ClientVehicle
// @Failure      400,500 {object} response
// @Router       /clients/{clientID}/vehicles [get]
func (h *Handler) clientVehicleGetAll(ctx echo.Context) error {
	clientID := ctx.Param("clientID")

	vehicles, err := h.usecases.Vehicle.GetAll(ctx.Request().Context(), clientID)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			return newResponse(ctx, http.StatusBadRequest, domain.ErrClientNotFound)
		}

		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK, vehicles)
}

// @Summary      Get by ID
// @Description  Get client's vehicle by ID
// @Security     UserAuth
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        clientID path string true "Client ID"
// @Param        vehicleID path string true "Vehicle ID"
// @Success      200 {object} domain.ClientVehicle
// @Failure      500 {object} response
// @Router       /clients/{clientID}/vehicles/{vehicleID} [get]
func (h *Handler) clientVehicleGetByID(ctx echo.Context) error {
	clientID := ctx.Param("clientID")
	vehicleID := ctx.Param("id")

	vehicle, err := h.usecases.Vehicle.GetByID(ctx.Request().Context(), clientID, vehicleID)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK, vehicle)
}

type clientVehicleUpdateRequest struct {
	Make        string `json:"make"`
	Model       string `json:"model"`
	Year        uint   `json:"year"`
	VIN         string `json:"vin"`
	Distance    uint   `json:"distance"`
	PlateNumber string `json:"plate_number"`
}

// @Summary      Update vehicle by ID
// @Description  Update client's vehicle by ID
// @Security     UserAuth
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        clientID path string true "Client ID"
// @Param        vehicleID path string true "Vehicle ID"
// @Param        request body clientVehicleUpdateRequest true "Request body"
// @Success      200 {object} response
// @Failure      400,500 {object} response
// @Router       /clients/{clientID}/vehicles/{vehicleID} [put]
func (h *Handler) clientVehicleUpdate(ctx echo.Context) error {
	var req clientVehicleUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	clientID := ctx.Param("clientID")
	vehicleID := ctx.Param("id")

	if err := h.usecases.Vehicle.Update(ctx.Request().Context(), clientID, vehicleID, usecase.VehicleUpdateInput{
		Make:        req.Make,
		Model:       req.Model,
		Year:        req.Year,
		VIN:         req.VIN,
		Distance:    req.Distance,
		PlateNumber: req.PlateNumber,
	}); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK)
}

// @Summary      Delete vehicle by ID
// @Description  Delete client's vehicle by ID
// @Security     UserAuth
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        clientID path string true "Client ID"
// @Param        vehicleID path string true "Vehicle ID"
// @Success      200 {object} response
// @Failure      500 {object} response
// @Router       /clients/{clientID}/vehicles/{vehicleID} [delete]
func (h *Handler) clientVehicleDelete(ctx echo.Context) error {
	clientID := ctx.Param("clientID")
	vehicleID := ctx.Param("id")

	if err := h.usecases.Vehicle.Delete(ctx.Request().Context(), clientID, vehicleID); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK)
}
