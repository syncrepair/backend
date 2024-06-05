package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
	"github.com/syncrepair/backend/pkg/hasher"
)

type UserUsecase interface {
	SignUp(ctx context.Context, input UserSignUpInput) (domain.UserTokens, error)
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

type UserSignUpInput struct {
	Name     string
	Email    string
	Password string
}

func (uc *userUsecase) SignUp(ctx context.Context, input UserSignUpInput) (domain.UserTokens, error) {
	if err := uc.repository.Create(ctx, domain.User{
		ID:          util.GenerateID(),
		Name:        input.Name,
		Email:       input.Email,
		Password:    uc.passwordHasher.Hash(input.Password),
		IsConfirmed: false,
	}); err != nil {
		return domain.UserTokens{}, err
	}

	return domain.UserTokens{}, nil
}
