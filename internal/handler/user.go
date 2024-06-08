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

type userSignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) SignUp(ctx echo.Context) error {
	var req userSignUpRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecase.SignUp(util.Ctx(ctx), usecase.UserSignUpRequest{
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

type userSignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) SignIn(ctx echo.Context) error {
	var req userSignInRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecase.SignIn(util.Ctx(ctx), usecase.UserSignInRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ErrorResponse(ctx, http.StatusNotFound, domain.ErrUserNotFound)
		}

		return ErrorResponse(ctx, http.StatusInternalServerError, domain.ErrInternalServer, err)
	}

	return SuccessResponse(ctx, http.StatusOK, tokens)
}
