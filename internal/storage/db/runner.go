package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// BeginTx начинает новую транзакцию в базе данных
// Возвращает pgx.Tx транзакцию или ошибку если транзакция не может быть начата
func (pg *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return pg.pool.Begin(ctx)
}

// Commit сохраняет изменение в базе данных
// Возвращает ошибку если транзакция не может быть сохранена
func (pg *Repository) Commit(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

// Rollback откатывает транзакцию в базе данных
// Возвращает ошибку если транзакция не может быть сохранена
func (pg *Repository) Rollback(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}
