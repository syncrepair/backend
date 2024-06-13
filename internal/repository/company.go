package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/syncrepair/backend/internal/domain"
)

type CompanyRepository interface {
	Create(ctx context.Context, company domain.Company) error
}

type companyRepository struct {
	db        *pgxpool.Pool
	sb        squirrel.StatementBuilderType
	tableName string
}

func NewCompanyRepository(db *pgxpool.Pool, sb squirrel.StatementBuilderType, tableName string) CompanyRepository {
	return &companyRepository{
		db:        db,
		sb:        sb,
		tableName: tableName,
	}
}

func (r *companyRepository) Create(ctx context.Context, company domain.Company) error {
	sql, args, err := r.sb.Insert(r.tableName).
		Columns("id", "name").
		Values(company.ID, company.Name).
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
