package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
)

func Init(log zerolog.Logger) *echo.Echo {
	e := echo.New()

	e.Logger = lecho.From(log)
	e.Use(requestLoggerMiddleware(log))

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
