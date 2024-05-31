package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(cfg Config) *pgxpool.Pool {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	))
	if err != nil {
		panic("error connecting to postgres database: " + err.Error())
	}

	return pool
}
