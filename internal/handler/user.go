package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"github.com/syncrepair/backend/internal/util"
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

type UserSignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) SignUp(ctx echo.Context) error {
	var req UserSignUpRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecase.SignUp(util.Ctx(ctx), domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return ErrorResponse(ctx, http.StatusConflict, domain.ErrUserAlreadyExists)
		}

		return ErrorResponse(ctx, http.StatusInternalServerError, domain.ErrInternalServer, err)
	}

	return SuccessResponse(ctx, http.StatusOK, tokens)
}

func (h *UserHandler) SignIn(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Sign In")
}
