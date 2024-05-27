package service

import (
	"context"
	"fmt"
	"github.com/syncrepair/backend/internal/entity"
	"github.com/syncrepair/backend/internal/repository"
)

type UserService interface {
	SignUp(context.Context, *entity.User) (*entity.UserTokens, error)
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}

func (s *userService) SignUp(ctx context.Context, user *entity.User) (*entity.UserTokens, error) {
	createdUser, err := s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &entity.UserTokens{
		AccessToken:  fmt.Sprintf("%s", createdUser.ID),
		RefreshToken: fmt.Sprintf("%s", createdUser.ID),
	}, nil
}
