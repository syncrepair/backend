package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(url string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		panic("error connecting to postgres database: " + err.Error())
	}

	return pool
}
