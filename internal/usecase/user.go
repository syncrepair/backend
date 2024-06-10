package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
	"github.com/syncrepair/backend/pkg/auth"
)

type UserUsecase interface {
	SignUp(ctx context.Context, req UserSignUpRequest) (UserTokens, error)
	SignIn(ctx context.Context, req UserSignInRequest) (UserTokens, error)
}

type userUsecase struct {
	repository     repository.UserRepository
	passwordHasher auth.PasswordHasher
	tokensManager  auth.TokensManager
}

func NewUserUsecase(repository repository.UserRepository, passwordHasher auth.PasswordHasher, tokensManager auth.TokensManager) UserUsecase {
	return &userUsecase{
		repository:     repository,
		passwordHasher: passwordHasher,
		tokensManager:  tokensManager,
	}
}

type UserSignUpRequest struct {
	Name      string
	Email     string
	Password  string
	CompanyID string
}

func (uc *userUsecase) SignUp(ctx context.Context, req UserSignUpRequest) (UserTokens, error) {
	id := util.GenerateID()

	if err := uc.repository.Create(ctx, domain.User{
		ID:          id,
		Name:        req.Name,
		Email:       req.Email,
		Password:    uc.passwordHasher.Hash(req.Password),
		CompanyID:   req.CompanyID,
		IsConfirmed: false,
	}); err != nil {
		return UserTokens{}, err
	}

	tokens, err := uc.generateTokens(id)
	if err != nil {
		return UserTokens{}, err
	}

	return tokens, nil
}

type UserSignInRequest struct {
	Email    string
	Password string
}

func (uc *userUsecase) SignIn(ctx context.Context, req UserSignInRequest) (UserTokens, error) {
	user, err := uc.repository.FindByCredentials(ctx, req.Email, uc.passwordHasher.Hash(req.Password))
	if err != nil {
		return UserTokens{}, err
	}

	tokens, err := uc.generateTokens(user.ID)
	if err != nil {
		return UserTokens{}, err
	}

	return tokens, nil
}

type UserTokens struct {
	AccessToken  string
	RefreshToken string
}

func (uc *userUsecase) generateTokens(id string) (UserTokens, error) {
	accessToken, err := uc.tokensManager.NewAccessToken(id)
	if err != nil {
		return UserTokens{}, err
	}

	refreshToken, err := uc.tokensManager.NewRefreshToken()
	if err != nil {
		return UserTokens{}, err
	}

	// TODO: storing sessions in db

	return UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
