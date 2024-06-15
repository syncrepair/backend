package http

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"net/http"
)

type Handlers struct {
	User    *UserHandler
	Company *CompanyHandler
}

func NewHandler(log zerolog.Logger, handlers Handlers) *echo.Echo {
	h := echo.New()

	h.Use(requestLoggingMiddleware(log))

	h.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	initAPI(h, handlers)

	return h
}

func initAPI(h *echo.Echo, handlers Handlers) {
	api := h.Group("/api")
	{
		handlers.User.initRoutes(api)
		handlers.Company.initRoutes(api)
	}
}
