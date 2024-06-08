package repository

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/util"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByCredentials(ctx context.Context, email string, password string) (domain.User, error)
}

type userRepository struct {
	db        *pgxpool.Pool
	sb        squirrel.StatementBuilderType
	tableName string
}

func NewUserRepository(db *pgxpool.Pool, sb squirrel.StatementBuilderType, tableName string) UserRepository {
	return &userRepository{
		db:        db,
		sb:        sb,
		tableName: tableName,
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	sql, args := r.sb.Insert(r.tableName).
		Columns("id", "name", "email", "password", "company_id", "is_confirmed").
		Values(user.ID, user.Name, user.Email, user.Password, user.CompanyID, user.IsConfirmed).
		MustSql()

	_, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		if errors.Is(util.ParsePgErr(err), util.PgErrAlreadyExists) {
			return domain.ErrUserAlreadyExists
		}

		return err
	}

	return nil
}

func (r *userRepository) FindByCredentials(ctx context.Context, email string, password string) (domain.User, error) {
	sql, args := r.sb.Select("id", "name", "email", "password", "company_id", "is_confirmed").
		From(r.tableName).
		Where(squirrel.And{
			squirrel.Eq{"email": email},
			squirrel.Eq{"password": password},
		}).
		MustSql()

	var user domain.User

	err := r.db.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CompanyID, &user.IsConfirmed)
	if err != nil {
		if errors.Is(util.ParsePgErr(err), util.PgErrNotFound) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}
