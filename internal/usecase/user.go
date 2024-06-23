package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/repository"
	"github.com/syncrepair/backend/pkg/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserUsecase struct {
	repository      repository.UserRepository
	passwordHasher  auth.PasswordHasher
	tokensManager   auth.TokensManager
	refreshTokenTTL time.Duration
}

func NewUserUsecase(repository repository.UserRepository, passwordHasher auth.PasswordHasher, tokensManager auth.TokensManager, refreshTokenTTL time.Duration) *UserUsecase {
	return &UserUsecase{
		repository:      repository,
		passwordHasher:  passwordHasher,
		tokensManager:   tokensManager,
		refreshTokenTTL: refreshTokenTTL,
	}
}

type UserSignUpInput struct {
	Name      string
	Email     string
	Password  string
	CompanyID string
}

func (uc *UserUsecase) SignUp(ctx context.Context, input UserSignUpInput) (UserTokens, error) {
	id := primitive.NewObjectID()

	companyID, err := primitive.ObjectIDFromHex(input.CompanyID)
	if err != nil {
		return UserTokens{}, err
	}

	if err := uc.repository.Create(ctx, domain.User{
		ID:          id,
		Name:        input.Name,
		Email:       input.Email,
		Password:    uc.passwordHasher.Hash(input.Password),
		IsConfirmed: false,
		Session:     domain.UserSession{},
		CompanyID:   companyID,
	}); err != nil {
		return UserTokens{}, err
	}

	return uc.createSession(ctx, id, companyID)
}

type UserSignInInput struct {
	Email    string
	Password string
}

func (uc *UserUsecase) SignIn(ctx context.Context, input UserSignInInput) (UserTokens, error) {
	user, err := uc.repository.GetByCredentials(ctx, input.Email, uc.passwordHasher.Hash(input.Password))
	if err != nil {
		return UserTokens{}, err
	}

	return uc.createSession(ctx, user.ID, user.CompanyID)
}

func (uc *UserUsecase) GetByID(ctx context.Context, id string) (domain.User, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}

	return uc.repository.GetByID(ctx, userID)
}

type UserUpdateInput struct {
	Name      string
	Email     string
	Password  string
	CompanyID string
}

func (uc *UserUsecase) Update(ctx context.Context, id string, input UserUpdateInput) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	companyID, err := primitive.ObjectIDFromHex(input.CompanyID)
	if err != nil {
		return err
	}

	return uc.repository.Update(ctx, domain.User{
		ID:        userID,
		Name:      input.Name,
		Email:     input.Email,
		Password:  uc.passwordHasher.Hash(input.Password),
		CompanyID: companyID,
	})
}

func (uc *UserUsecase) Delete(ctx context.Context, id string) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return uc.repository.Delete(ctx, userID)
}

func (uc *UserUsecase) Confirm(ctx context.Context, id string) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return uc.repository.Confirm(ctx, userID)
}

type UserTokens struct {
	AccessToken  string
	RefreshToken string
}

func (uc *UserUsecase) Refresh(ctx context.Context, refreshToken string) (UserTokens, error) {
	user, err := uc.repository.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return UserTokens{}, err
	}

	return uc.createSession(ctx, user.ID, user.CompanyID)
}

func (uc *UserUsecase) createSession(ctx context.Context, userID, companyID primitive.ObjectID) (UserTokens, error) {
	accessToken, err := uc.tokensManager.NewAccessToken(auth.Claims{
		UserID:    userID.Hex(),
		CompanyID: companyID.Hex(),
	})
	if err != nil {
		return UserTokens{}, err
	}

	refreshToken, err := uc.tokensManager.NewRefreshToken()
	if err != nil {
		return UserTokens{}, err
	}

	if err := uc.repository.SetSession(ctx, userID, domain.UserSession{
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(uc.refreshTokenTTL),
	}); err != nil {
		return UserTokens{}, err
	}

	return UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
