package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/syncrepair/backend/internal/domain"
	"net/http"
)

type successResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func SuccessResponse(ctx echo.Context, statusCode int, data ...interface{}) error {
	res := successResponse{
		Status: "success",
	}

	if len(data) == 1 {
		res.Data = data[0]
	} else if len(data) > 1 {
		res.Data = data
	}

	return ctx.JSON(statusCode, res)
}

func ErrorResponse(ctx echo.Context, statusCode int, err error) error {
	res := errorResponse{
		Status:  "error",
		Message: err.Error(),
	}

	if statusCode >= http.StatusInternalServerError {
		res.Message = domain.ErrInternalServer.Error()
		log.Error().
			Err(err).
			Msg("server error")
	}

	return ctx.JSON(statusCode, res)
}
