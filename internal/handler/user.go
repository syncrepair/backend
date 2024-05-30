package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/usecase"
	"net/http"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) Routes(router *echo.Group) {
	users := router.Group("/users")
	{
		users.POST("/signup", h.SignUp)
		users.POST("/signin", h.SignIn)
	}
}

func (h *UserHandler) SignUp(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Sign Up")
}

func (h *UserHandler) SignIn(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Sign In")
}
