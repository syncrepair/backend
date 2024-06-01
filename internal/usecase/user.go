package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
	"github.com/syncrepair/backend/pkg/hasher"
)

type UserUsecase interface {
	SignUp(context.Context, domain.User) (domain.UserTokens, error)
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

func (uc *userUsecase) SignUp(ctx context.Context, user domain.User) (domain.UserTokens, error) {
	user.ID = util.GenerateID()
	user.IsConfirmed = false
	user.Password = uc.passwordHasher.Hash(user.Password)

	if err := uc.repository.Create(ctx, user); err != nil {
		return domain.UserTokens{}, err
	}

	return domain.UserTokens{}, nil
}
