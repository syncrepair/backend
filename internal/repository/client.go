package repository

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/util"
)

type ClientRepository interface {
	Create(ctx context.Context, client domain.Client) error
	Delete(ctx context.Context, id string) error
}

type clientRepository struct {
	db        *pgxpool.Pool
	sb        squirrel.StatementBuilderType
	tableName string
}

func NewClientRepository(db *pgxpool.Pool, sb squirrel.StatementBuilderType, tableName string) ClientRepository {
	return &clientRepository{
		db:        db,
		sb:        sb,
		tableName: tableName,
	}
}

func (r *clientRepository) Create(ctx context.Context, client domain.Client) error {
	sql, args, err := r.sb.Insert(r.tableName).
		Columns("id", "name", "phone_number", "company_id").
		Values(client.ID, client.Name, client.PhoneNumber, client.CompanyID).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		if errors.Is(util.ParsePgErr(err), util.PgErrForeignKey) {
			return domain.ErrCompanyNotFound
		}

		return err
	}

	return nil
}

func (r *clientRepository) Delete(ctx context.Context, id string) error {
	sql, args, err := r.sb.Delete(r.tableName).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		if errors.Is(util.ParsePgErr(err), util.PgErrForeignKey) {
			return domain.ErrCompanyNotFound
		}

		return err
	}

	return nil
}
