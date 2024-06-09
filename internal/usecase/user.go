package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
	"github.com/syncrepair/backend/pkg/auth"
)

type UserUsecase interface {
	SignUp(ctx context.Context, req UserSignUpRequest) (string, error)
	SignIn(ctx context.Context, req UserSignInRequest) (string, error)
}

type userUsecase struct {
	repository     repository.UserRepository
	passwordHasher auth.PasswordHasher
	jwtManager     auth.JWTManager
}

func NewUserUsecase(repository repository.UserRepository, passwordHasher auth.PasswordHasher, jwtManager auth.JWTManager) UserUsecase {
	return &userUsecase{
		repository:     repository,
		passwordHasher: passwordHasher,
		jwtManager:     jwtManager,
	}
}

type UserSignUpRequest struct {
	Name      string
	Email     string
	Password  string
	CompanyID string
}

func (uc *userUsecase) SignUp(ctx context.Context, req UserSignUpRequest) (string, error) {
	id := util.GenerateID()

	if err := uc.repository.Create(ctx, domain.User{
		ID:          id,
		Name:        req.Name,
		Email:       req.Email,
		Password:    uc.passwordHasher.Hash(req.Password),
		CompanyID:   req.CompanyID,
		IsConfirmed: false,
	}); err != nil {
		return "", err
	}

	return uc.jwtManager.GenerateToken(id), nil
}

type UserSignInRequest struct {
	Email    string
	Password string
}

func (uc *userUsecase) SignIn(ctx context.Context, req UserSignInRequest) (string, error) {
	user, err := uc.repository.FindByCredentials(ctx, req.Email, uc.passwordHasher.Hash(req.Password))
	if err != nil {
		return "", err
	}

	return uc.jwtManager.GenerateToken(user.ID), nil
}
