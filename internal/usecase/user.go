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
	Confirm(ctx context.Context, id string) error
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

	return uc.createSession(ctx, id, req.CompanyID)
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

	return uc.createSession(ctx, user.ID, user.CompanyID)
}

func (uc *userUsecase) Confirm(ctx context.Context, id string) error {
	return uc.repository.Confirm(ctx, id)
}

func (uc *userUsecase) RefreshTokens(ctx context.Context, refreshToken string) (UserTokens, error) {
	var expiresAt time.Time
	var userID string
	var companyID string

	if err := uc.redisDB.HGet(ctx, refreshToken, "expires_at").Scan(&expiresAt); err != nil {
		return UserTokens{}, err
	}
	if time.Now().After(expiresAt) {
		uc.redisDB.HDel(ctx, refreshToken, "user_id", "company_id", "expires_at")
		return UserTokens{}, domain.ErrUnauthorized
	}

	if err := uc.redisDB.HGet(ctx, refreshToken, "user_id").Scan(&userID); err != nil {
		return UserTokens{}, err
	}

	if err := uc.redisDB.HGet(ctx, refreshToken, "company_id").Scan(&companyID); err != nil {
		return UserTokens{}, err
	}

	accessToken, err := uc.tokensManager.NewAccessToken(auth.Claims{
		UserID:    userID,
		CompanyID: companyID,
	})
	if err != nil {
		return UserTokens{}, err
	}

	return UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type UserTokens struct {
	AccessToken  string
	RefreshToken string
}

func (uc *userUsecase) createSession(ctx context.Context, userID, companyID string) (UserTokens, error) {
	accessToken, err := uc.tokensManager.NewAccessToken(auth.Claims{
		UserID:    userID,
		CompanyID: companyID,
	})
	if err != nil {
		return UserTokens{}, err
	}

	refreshToken, err := uc.tokensManager.NewRefreshToken()
	if err != nil {
		return UserTokens{}, err
	}

	uc.redisDB.HSet(ctx, refreshToken, "user_id", userID, "company_id", companyID, "expires_at", time.Now().Add(uc.refreshTokenTTL))

	return UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
