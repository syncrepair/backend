package http

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/syncrepair/backend/internal/domain"
	"net/http"
)

type dataResponse struct {
	Data  interface{} `json:"data"`
	Count uint        `json:"count"`
}

type idResponse struct {
	ID string `json:"id"`
}

type response struct {
	Message string `json:"message"`
}

func newResponse(ctx echo.Context, statusCode int, data ...interface{}) error {
	if len(data) == 0 {
		return ctx.JSON(statusCode, response{
			Message: http.StatusText(statusCode),
		})
	}

	switch data[0].(type) {
	case error:
		err := data[0].(error)
		msg := err.Error()

		if statusCode >= http.StatusInternalServerError {
			msg = domain.ErrInternalServer.Error()
			log.Error().
				Err(err).
				Msg("server error")
		}

		return ctx.JSON(statusCode, response{
			Message: msg,
		})
	default:
		return ctx.JSON(statusCode, data[0])
	}
}
