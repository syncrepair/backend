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
	Confirm(ctx context.Context, id string) error
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
	sql, args, err := r.sb.Insert(r.tableName).
		Columns("id", "name", "email", "password", "company_id", "is_confirmed").
		Values(user.ID, user.Name, user.Email, user.Password, user.CompanyID, user.IsConfirmed).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		if errors.Is(util.ParsePgErr(err), util.PgErrAlreadyExists) {
			return domain.ErrUserAlreadyExists
		}
		if errors.Is(util.ParsePgErr(err), util.PgErrForeignKey) {
			return domain.ErrCompanyNotFound
		}

		return err
	}

	return nil
}

func (r *userRepository) FindByCredentials(ctx context.Context, email string, password string) (domain.User, error) {
	sql, args, err := r.sb.Select("id", "name", "email", "password", "company_id", "is_confirmed").
		From(r.tableName).
		Where(squirrel.And{
			squirrel.Eq{"email": email},
			squirrel.Eq{"password": password},
		}).
		ToSql()
	if err != nil {
		return domain.User{}, err
	}

	var user domain.User

	if err = r.db.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CompanyID, &user.IsConfirmed); err != nil {
		if errors.Is(util.ParsePgErr(err), util.PgErrNotFound) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) Confirm(ctx context.Context, id string) error {
	sql, args, err := r.sb.Update(r.tableName).
		Set("is_confirmed", true).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Query(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
