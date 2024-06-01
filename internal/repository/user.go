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
	Create(context.Context, domain.User) error
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
		Columns("id", "name", "email", "password", "is_confirmed").
		Values(user.ID, user.Name, user.Email, user.Password, user.IsConfirmed).
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
