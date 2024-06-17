package http

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/auth"
	"net/http"
)

type Handler struct {
	usecases      Usecases
	tokensManager auth.TokensManager
	log           zerolog.Logger
}

type Usecases struct {
	UserUsecase    usecase.UserUsecase
	CompanyUsecase usecase.CompanyUsecase
	ServiceUsecase usecase.ServiceUsecase
	ClientUsecase  usecase.ClientUsecase
}

func NewHandler(log zerolog.Logger, tokensManager auth.TokensManager, usecases Usecases) *Handler {
	return &Handler{
		usecases:      usecases,
		tokensManager: tokensManager,
		log:           log,
	}
}

func (h *Handler) Init() *echo.Echo {
	r := echo.New()

	r.Use(h.requestLoggingMiddleware())

	r.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	h.initAPI(r)

	return r
}

func (h *Handler) initAPI(router *echo.Echo) {
	api := router.Group("/api")
	{
		h.initUserRoutes(api)
		h.initCompanyRoutes(api)
		h.initServiceRoutes(api)
		h.initClientRoutes(api)
	}
}
