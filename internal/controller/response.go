package controller

import (
	"github.com/labstack/echo/v4"
)

type successResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func SuccessResponse(ctx echo.Context, code int, data ...interface{}) error {
	res := successResponse{
		Status: "success",
	}

	if len(data) == 1 {
		res.Data = data[0]
	} else if len(data) > 1 {
		res.Data = data
	}

	return ctx.JSON(code, res)
}

func ErrorResponse(ctx echo.Context, code int, err error, appErr ...error) error {
	res := errorResponse{
		Status:  "error",
		Message: err.Error(),
	}

	return ctx.JSON(code, res)
}
