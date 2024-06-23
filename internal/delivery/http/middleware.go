package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/syncrepair/backend/internal/domain"
	"net/http"
	"strings"
)

func (h *Handler) requestLoggingMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:  true,
		LogURI:     true,
		LogError:   true,
		LogLatency: true,
		LogMethod:  true,
		LogValuesFunc: func(c echo.Context, req middleware.RequestLoggerValues) error {
			h.log.Info().
				Str("uri", req.URI).
				Int("status", req.Status).
				Str("latency", fmt.Sprintf("%dms", req.Latency.Milliseconds())).
				Msg(req.Method)

			return nil
		},
	})
}

func (h *Handler) authMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			header := ctx.Request().Header.Get("Authorization")
			if header == "" {
				return newResponse(ctx, http.StatusUnauthorized, domain.ErrUnauthorized)
			}

			headerParts := strings.Split(header, " ")
			token := headerParts[1]

			claims, err := h.tokensManager.GetAccessTokenClaims(token)
			if err != nil {
				return newResponse(ctx, http.StatusUnauthorized, domain.ErrUnauthorized)
			}

			ctx.Set("userID", claims.UserID)
			ctx.Set("companyID", claims.CompanyID)

			return next(ctx)
		}
	}
}

func getUserIDFromCtx(ctx echo.Context) string {
	return ctx.Get("userID").(string)
}

func getCompanyIDFromCtx(ctx echo.Context) string {
	return ctx.Get("companyID").(string)
}
