package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (pg *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return pg.pool.Begin(ctx)
}
