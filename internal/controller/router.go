package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func NewRouter(logger zerolog.Logger) *echo.Echo {
	r := echo.New()
	r.Use(requestLoggingMiddleware(logger))

	return r
}

func requestLoggingMiddleware(log zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:  true,
		LogURI:     true,
		LogError:   true,
		LogLatency: true,
		LogMethod:  true,
		LogValuesFunc: func(c echo.Context, req middleware.RequestLoggerValues) error {
			log.Info().
				Str("uri", req.URI).
				Int("status", req.Status).
				Str("latency", fmt.Sprintf("%dms", req.Latency.Milliseconds())).
				Msg(req.Method)

			return nil
		},
	})
}
