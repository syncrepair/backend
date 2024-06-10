package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/internal/util"
	"github.com/syncrepair/backend/pkg/auth"
	"github.com/syncrepair/backend/pkg/redis"
	"time"
)

type UserUsecase interface {
	SignUp(ctx context.Context, req UserSignUpRequest) (UserTokens, error)
	SignIn(ctx context.Context, req UserSignInRequest) (UserTokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (UserTokens, error)
}

type userUsecase struct {
	repository      repository.UserRepository
	passwordHasher  auth.PasswordHasher
	tokensManager   auth.TokensManager
	redisDB         *redis.Redis
	refreshTokenTTL time.Duration
}

func NewUserUsecase(repository repository.UserRepository, passwordHasher auth.PasswordHasher, tokensManager auth.TokensManager, redisDB *redis.Redis, refreshTokenTTL time.Duration) UserUsecase {
	return &userUsecase{
		repository:      repository,
		passwordHasher:  passwordHasher,
		tokensManager:   tokensManager,
		redisDB:         redisDB,
		refreshTokenTTL: refreshTokenTTL,
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

	return uc.createSession(ctx, id)
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

	return uc.createSession(ctx, user.ID)
}

func (uc *userUsecase) RefreshTokens(ctx context.Context, refreshToken string) (UserTokens, error) {
	var userID string

	if err := uc.redisDB.Get(ctx, refreshToken).Scan(&userID); err != nil {
		return UserTokens{}, err
	}

	return uc.createSession(ctx, userID)
}

type UserTokens struct {
	AccessToken  string
	RefreshToken string
}

func (uc *userUsecase) createSession(ctx context.Context, userID string) (UserTokens, error) {
	accessToken, err := uc.tokensManager.NewAccessToken(userID)
	if err != nil {
		return UserTokens{}, err
	}

	refreshToken, err := uc.tokensManager.NewRefreshToken()
	if err != nil {
		return UserTokens{}, err
	}

	uc.redisDB.Set(ctx, refreshToken, userID, uc.refreshTokenTTL)

	return UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
