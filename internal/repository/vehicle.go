package repository

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/syncrepair/backend/internal/domain"
	"github.com/syncrepair/backend/internal/util"
)

type VehicleRepository interface {
	Create(ctx context.Context, vehicle domain.Vehicle) error
	Delete(ctx context.Context, id string) error
}

type vehicleRepository struct {
	db        *pgxpool.Pool
	sb        squirrel.StatementBuilderType
	tableName string
}

func NewVehicleRepository(db *pgxpool.Pool, sb squirrel.StatementBuilderType, tableName string) VehicleRepository {
	return &vehicleRepository{
		db:        db,
		sb:        sb,
		tableName: tableName,
	}
}

func (r *vehicleRepository) Create(ctx context.Context, vehicle domain.Vehicle) error {
	sql, args, err := r.sb.Insert(r.tableName).
		Columns("id", "make", "model", "year", "vin", "plate_number", "client_id").
		Values(vehicle.ID, vehicle.Make, vehicle.Model, vehicle.Year, vehicle.VIN, vehicle.PlateNumber, vehicle.ClientID).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		if errors.Is(util.ParsePgErr(err), util.PgErrForeignKey) {
			return domain.ErrClientNotFound
		}

		return err
	}

	return nil
}

func (r *vehicleRepository) Delete(ctx context.Context, id string) error {
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
