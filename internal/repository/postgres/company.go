package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/syncrepair/backend/internal/model"
)

const companiesTable = "companies"

type CompanyRepository struct {
	db *pgxpool.Pool
	sb squirrel.StatementBuilderType
}

func NewCompanyRepository(db *pgxpool.Pool, statementBuilder squirrel.StatementBuilderType) *CompanyRepository {
	return &CompanyRepository{
		db: db,
		sb: statementBuilder,
	}
}

func (r *CompanyRepository) Create(ctx context.Context, company *model.Company) error {
	sql, args, err := r.sb.Insert(companiesTable).
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
