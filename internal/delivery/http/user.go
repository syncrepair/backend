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
		users.POST("/confirm/:id", h.userConfirm)
		users.POST("/refresh", h.userRefresh)

		authenticated := users.Group("", h.authMiddleware())
		{
			authenticated.GET("", h.userGet)
			authenticated.PUT("", h.userUpdate)
			authenticated.DELETE("", h.userDelete)
		}
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

// @Summary      Sign up
// @Description  Create new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body userSignUpRequest true "Request body"
// @Success      201 {object} userSignUpResponse
// @Failure      400,500 {object} response
// @Router       /users/sign-up [post]
func (h *Handler) userSignUp(ctx echo.Context) error {
	var req userSignUpRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecases.User.SignUp(ctx.Request().Context(), usecase.UserSignUpInput{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		CompanyID: req.CompanyID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return newResponse(ctx, http.StatusBadRequest, domain.ErrUserAlreadyExists)
		}

		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusCreated, userSignUpResponse{
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

// @Summary      Sign in
// @Description  Sign in to user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body userSignInRequest true "Request body"
// @Success      200 {object} userSignInResponse
// @Failure      400,500 {object} response
// @Router       /users/sign-in [post]
func (h *Handler) userSignIn(ctx echo.Context) error {
	var req userSignInRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecases.User.SignIn(ctx.Request().Context(), usecase.UserSignInInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return newResponse(ctx, http.StatusBadRequest, domain.ErrUserNotFound)
		}

		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK, userSignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

type userRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type userRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary      Refresh tokens
// @Description  Refresh user tokens
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body userRefreshRequest true "Request body"
// @Success      200 {object} userRefreshResponse
// @Failure      400,401,500 {object} response
// @Router       /users/refresh [post]
func (h *Handler) userRefresh(ctx echo.Context) error {
	var req userRefreshRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.usecases.User.Refresh(ctx.Request().Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return newResponse(ctx, http.StatusUnauthorized, domain.ErrUnauthorized)
		}

		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK, userRefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Summary      Confirm
// @Description  Confirm user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} response
// @Failure      500 {object} response
// @Router       /users/confirm/{id} [post]
func (h *Handler) userConfirm(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := h.usecases.User.Confirm(ctx.Request().Context(), id); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK)
}

// @Summary      Get
// @Description  Get user
// @Security     UserAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200 {object} domain.User
// @Failure      400,500 {object} response
// @Router       /users [get]
func (h *Handler) userGet(ctx echo.Context) error {
	user, err := h.usecases.User.GetByID(ctx.Request().Context(), getUserIDFromCtx(ctx))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return newResponse(ctx, http.StatusBadRequest, domain.ErrUserNotFound)
		}

		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK, user)
}

type userUpdateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary      Update
// @Description  Update user
// @Security     UserAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body userUpdateRequest true "Request body"
// @Success      200 {object} response
// @Failure      400,500 {object} response
// @Router       /users [put]
func (h *Handler) userUpdate(ctx echo.Context) error {
	var req userUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return newResponse(ctx, http.StatusBadRequest, domain.ErrBadRequest)
	}

	if err := h.usecases.User.Update(ctx.Request().Context(), getUserIDFromCtx(ctx), usecase.UserUpdateInput{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		CompanyID: getCompanyIDFromCtx(ctx),
	}); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK)
}

// @Summary      Delete
// @Description  Delete user
// @Security     UserAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200 {object} response
// @Failure      500 {object} response
// @Router       /users [delete]
func (h *Handler) userDelete(ctx echo.Context) error {
	if err := h.usecases.User.Delete(ctx.Request().Context(), getUserIDFromCtx(ctx)); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err)
	}

	return newResponse(ctx, http.StatusOK)
}
