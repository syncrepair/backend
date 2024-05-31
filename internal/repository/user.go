package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/syncrepair/backend/internal/domain"
)

type UserRepository interface {
	Create(context.Context, *domain.User) (*domain.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	_, err := r.db.Exec(ctx, "insert into users (id, name, email, password, is_confirmed) values ($1, $2, $3, $4, $5)", user.ID, user.Name, user.Email, user.Password, user.IsConfirmed)
	if err != nil {
		return nil, err
	}

	return user, nil
}
