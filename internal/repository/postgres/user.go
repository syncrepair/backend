package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/syncrepair/backend/internal/model"
)

const usersTable = "users"

type UserRepository struct {
	db *pgxpool.Pool
	sb squirrel.StatementBuilderType
}

func NewUserRepository(db *pgxpool.Pool, statementBuilder squirrel.StatementBuilderType) *UserRepository {
	return &UserRepository{
		db: db,
		sb: statementBuilder,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	sql, args, err := r.sb.Insert(usersTable).
		Columns("id", "name", "email", "password", "company_id", "created_at", "updated_at").
		Values(user.ID, user.Name, user.Email, user.Password, user.CompanyID, user.CreatedAt, user.UpdatedAt).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
