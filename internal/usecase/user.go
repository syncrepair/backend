package usecase

import (
	"context"
	"github.com/syncrepair/backend/internal/model"
	"github.com/syncrepair/backend/internal/utils/id"
	"github.com/syncrepair/backend/pkg/password"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
}

type UserUsecase struct {
	repository        UserRepository
	companyRepository CompanyRepository
	passwordHasher    password.Hasher
}

func NewUserUsecase(repository UserRepository, companyRepository CompanyRepository, passwordHasher password.Hasher) *UserUsecase {
	return &UserUsecase{
		repository:        repository,
		companyRepository: companyRepository,
		passwordHasher:    passwordHasher,
	}
}

func (u *UserUsecase) SignUp(ctx context.Context, input *model.UserSignUpInput) error {
	var companyID string

	if input.CompanyCode != "" {
		company, err := u.companyRepository.GetByCode(ctx, input.CompanyCode)
		if err != nil {
			return err
		}

		companyID = company.ID
	}

	user := &model.User{
		ID:        id.Generate(),
		Name:      input.Name,
		Email:     input.Email,
		Password:  u.passwordHasher.Hash(input.Password),
		CompanyID: companyID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.repository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}
