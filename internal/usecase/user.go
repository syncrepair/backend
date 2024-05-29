package usecase

import (
	"context"
	"fmt"
	"github.com/syncrepair/backend/internal/entity"
	"github.com/syncrepair/backend/internal/repository"
)

type UserUsecase interface {
	SignUp(context.Context, *entity.User) (*entity.UserTokens, error)
}

type userUsecase struct {
	repository repository.UserRepository
}

func NewUserUsecase(repository repository.UserRepository) UserUsecase {
	return &userUsecase{
		repository: repository,
	}
}

func (uc *userUsecase) SignUp(ctx context.Context, user *entity.User) (*entity.UserTokens, error) {
	createdUser, err := uc.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &entity.UserTokens{
		AccessToken:  fmt.Sprintf("%s", createdUser.ID),
		RefreshToken: fmt.Sprintf("%s", createdUser.ID),
	}, nil
}
