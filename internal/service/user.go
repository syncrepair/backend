package service

import (
	"context"
	"fmt"
	"github.com/syncrepair/backend/internal/entity"
	"github.com/syncrepair/backend/internal/repository"
)

type User interface {
	SignUp(context.Context, *entity.User) (*entity.UserTokens, error)
}

type user struct {
	repository repository.User
}

func NewUser(repository repository.User) User {
	return &user{
		repository: repository,
	}
}

func (s *user) SignUp(ctx context.Context, user *entity.User) (*entity.UserTokens, error) {
	createdUser, err := s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &entity.UserTokens{
		AccessToken:  fmt.Sprintf("%s", createdUser.ID),
		RefreshToken: fmt.Sprintf("%s", createdUser.ID),
	}, nil
}
