package http

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/syncrepair/backend/docs"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/pkg/auth"
)

type Handler struct {
	usecases      Usecases
	tokensManager auth.TokensManager
	log           zerolog.Logger
}

type Usecases struct {
	User    *usecase.UserUsecase
	Company *usecase.CompanyUsecase
	Service *usecase.ServiceUsecase
	Client  *usecase.ClientUsecase
	Vehicle *usecase.VehicleUsecase
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

	r.GET("/swagger/*", echoSwagger.WrapHandler)

	r.Use(h.requestLoggingMiddleware())

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
