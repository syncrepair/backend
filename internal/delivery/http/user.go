package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/usecase"
	"net/http"
)

func (h *Handler) initUserRoutes(router *echo.Group) {
	users := router.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
		users.POST("/confirm", h.userConfirm)
		users.POST("/refresh-tokens", h.userRefreshTokens)
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

func (h *Handler) userSignUp(ctx echo.Context) error {
	var req userSignUpRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecases.UserUsecase.SignUp(ctx.Request().Context(), usecase.UserSignUpRequest{
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

func (h *Handler) userSignIn(ctx echo.Context) error {
	var req userSignInRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecases.UserUsecase.SignIn(ctx.Request().Context(), usecase.UserSignInRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrUserNotFound)
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

func (h *Handler) userRefreshTokens(ctx echo.Context) error {
	var req userRefreshTokensRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecases.UserUsecase.RefreshTokens(ctx.Request().Context(), req.RefreshToken)
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

func (h *Handler) userConfirm(ctx echo.Context) error {
	var req userConfirmRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	if err := h.usecases.UserUsecase.Confirm(ctx.Request().Context(), req.ID); err != nil {
		if errors.Is(err, domain.ErrUserConfirmation) {
			return ErrorResponse(ctx, http.StatusBadRequest, domain.ErrUserConfirmation)
		}

		return ErrorResponse(ctx, http.StatusInternalServerError, err)
	}

	return SuccessResponse(ctx, http.StatusOK)
}
