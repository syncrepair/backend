package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
)

type UserUsecase interface {
	SignUp(context.Context, domain.User) (domain.UserTokens, error)
}

type userUsecase struct {
	repository repository.UserRepository
}

func NewUserUsecase(repository repository.UserRepository) UserUsecase {
	return &userUsecase{
		repository: repository,
	}
}

func (uc *userUsecase) SignUp(ctx context.Context, user domain.User) (domain.UserTokens, error) {
	user.ID = util.GenerateID()
	user.IsConfirmed = false

	// TODO: password hashing

	if err := uc.repository.Create(ctx, user); err != nil {
		return domain.UserTokens{}, err
	}

	return domain.UserTokens{}, nil
}
