package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
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

func (h *UserHandler) initRoutes(router *echo.Group) {
	users := router.Group("/users")
	{
		users.POST("/sign-up", h.signUp)
		users.POST("/sign-in", h.signIn)
		users.POST("/confirm", h.confirm)
		users.POST("/refresh-tokens", h.refreshTokens)
	}
}

type userSignUpRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CompanyID string `json:"company_id"`
}

type userSignUpResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *UserHandler) signUp(ctx echo.Context) error {
	var req userSignUpRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecase.SignUp(ctx.Request().Context(), usecase.UserSignUpRequest{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		CompanyID: req.CompanyID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return ErrorResponse(ctx, http.StatusConflict, domain.ErrUserAlreadyExists)
		}
		if errors.Is(err, domain.ErrCompanyNotFound) {
			return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrCompanyNotFound)
		}

		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK, userSignUpResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

type userSignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userSignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *UserHandler) signIn(ctx echo.Context) error {
	var req userSignInRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecase.SignIn(ctx.Request().Context(), usecase.UserSignInRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ErrorResponse(ctx, http.StatusNotFound, domain.ErrUserNotFound)
		}

		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK, userSignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

type userRefreshTokensRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type userRefreshTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *UserHandler) refreshTokens(ctx echo.Context) error {
	var req userRefreshTokensRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecase.RefreshTokens(ctx.Request().Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			return ErrorResponse(ctx, http.StatusUnauthorized, domain.ErrUnauthorized)
		}

		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK, userRefreshTokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

type userConfirmRequest struct {
	ID string `json:"id"`
}

func (h *UserHandler) confirm(ctx echo.Context) error {
	var req userConfirmRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	if err := h.usecase.Confirm(ctx.Request().Context(), req.ID); err != nil {
		if errors.Is(err, domain.ErrUserConfirmation) {
			return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrUserConfirmation)
		}

		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK)
}
