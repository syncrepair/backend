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
	sql, args, err := r.sb.
		Insert(companiesTable).
		Columns("id", "name", "code", "created_at", "updated_at").
		Values(company.ID, company.Name, company.Code, company.CreatedAt, company.UpdatedAt).
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

func (r *CompanyRepository) GetByCode(ctx context.Context, code string) (*model.Company, error) {
	sql, args, err := r.sb.
		Select("id", "name", "code", "created_at", "updated_at").
		From(companiesTable).
		Where(squirrel.Eq{"code": code}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var company model.Company

	if err := r.db.QueryRow(ctx, sql, args...).Scan(&company.ID, &company.Name, &company.Code, &company.CreatedAt, &company.UpdatedAt); err != nil {
		return nil, err
	}

	return &company, nil
}
