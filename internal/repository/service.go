package repository

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/util"
)

type ServiceRepository interface {
	Create(ctx context.Context, service domain.Service) error
	Delete(ctx context.Context, id string) error
}

type serviceRepository struct {
	db        *pgxpool.Pool
	sb        squirrel.StatementBuilderType
	tableName string
}

func NewServiceRepository(db *pgxpool.Pool, sb squirrel.StatementBuilderType, tableName string) ServiceRepository {
	return &serviceRepository{
		db:        db,
		sb:        sb,
		tableName: tableName,
	}
}

func (r *serviceRepository) Create(ctx context.Context, service domain.Service) error {
	sql, args, err := r.sb.Insert(r.tableName).
		Columns("id", "name", "description", "price", "company_id").
		Values(service.ID, service.Name, service.Description, service.Price, service.CompanyID).
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

func (r *serviceRepository) Delete(ctx context.Context, id string) error {
	sql, args, err := r.sb.Delete(r.tableName).
		Where(squirrel.Eq{"id": id}).
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
