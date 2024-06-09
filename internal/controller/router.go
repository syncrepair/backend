package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/ziflex/lecho/v3"
	"net/http"
)

func NewRouter(logger zerolog.Logger) *echo.Echo {
	e := echo.New()

	e.Logger = lecho.From(logger)
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		code := http.StatusInternalServerError

		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			code = httpError.Code
		}

		if code >= http.StatusInternalServerError {
			ErrorResponse(ctx, code, domain.ErrInternalServer, err) // TODO: change router logic
		}
	}

	e.Use(requestLoggerMiddleware(logger))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogLevel: log.OFF,
	}))

	return e
}

func requestLoggerMiddleware(log zerolog.Logger) echo.MiddlewareFunc {
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
				Int64("latency", req.Latency.Milliseconds()).
				Msg(req.Method)

			return nil
		},
	})
}
