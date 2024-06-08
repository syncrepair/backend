package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
	"github.com/syncrepair/backend/pkg/hasher"
)

type UserUsecase interface {
	SignUp(ctx context.Context, req UserSignUpRequest) (domain.UserTokens, error)
	SignIn(ctx context.Context, req UserSignInRequest) (domain.UserTokens, error)
}

type userUsecase struct {
	repository     repository.UserRepository
	passwordHasher hasher.Hasher
}

func NewUserUsecase(repository repository.UserRepository, passwordHasher hasher.Hasher) UserUsecase {
	return &userUsecase{
		repository:     repository,
		passwordHasher: passwordHasher,
	}
}

type UserSignUpRequest struct {
	Name     string
	Email    string
	Password string
}

func (uc *userUsecase) SignUp(ctx context.Context, req UserSignUpRequest) (domain.UserTokens, error) {
	if err := uc.repository.Create(ctx, domain.User{
		ID:          util.GenerateID(),
		Name:        req.Name,
		Email:       req.Email,
		Password:    uc.passwordHasher.Hash(req.Password),
		IsConfirmed: false,
	}); err != nil {
		return domain.UserTokens{}, err
	}

	// TODO: generating JWT

	return domain.UserTokens{}, nil
}

type UserSignInRequest struct {
	Email    string
	Password string
}

func (uc *userUsecase) SignIn(ctx context.Context, req UserSignInRequest) (domain.UserTokens, error) {
	_, err := uc.repository.FindByCredentials(ctx, req.Email, uc.passwordHasher.Hash(req.Password))
	if err != nil {
		return domain.UserTokens{}, err
	}

	// TODO: generating JWT

	return domain.UserTokens{}, nil
}
