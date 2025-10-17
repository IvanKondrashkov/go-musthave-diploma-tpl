package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/IvanKondrashkov/go-market/internal/models"
	customError "github.com/IvanKondrashkov/go-market/internal/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// BalanceWithdraw списание баллов с накопительного счёта в счёт оплаты нового заказа, сохранение в БД PostgreSQL
// Принимает:
// - ctx: контекст с информацией о пользователе
// - tx: транзакцию
// - userID: идентификатор пользователя
// - withdraw: models.UserWithdrawals
// Возвращает:
// - ошибку, если не достаточно денег для списания (ErrInvalidUserWithdraw) или возникли проблемы при сохранении
func (pg *Repository) BalanceWithdraw(ctx context.Context, tx pgx.Tx, userID uuid.UUID, withdraw models.UserWithdrawals) error {
	query := `
	SELECT current FROM balances WHERE user_id = $1 FOR UPDATE
	`

	var current float64
	err := tx.QueryRow(ctx, query, userID).Scan(&current)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("withdraw balance error: %w", customError.ErrInvalidUserWithdraw)
		}
		return fmt.Errorf("withdraw balance error: %w", err)
	}

	if current < withdraw.Sum {
		return fmt.Errorf("withdraw balance error: %w", customError.ErrInvalidUserWithdraw)
	}

	query = `
	UPDATE balances SET current = current - $1, withdrawn = withdrawn + $2 WHERE user_id = $3
	`

	_, err = tx.Exec(ctx, query, withdraw.Sum, withdraw.Sum, userID)
	if err != nil {
		return fmt.Errorf("update withdraw balance error: %w", err)
	}

	query = `
	INSERT INTO withdrawals(id, user_id, order_number, sum) VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(ctx, query, uuid.New(), userID, withdraw.Order, withdraw.Sum)
	if err != nil {
		return fmt.Errorf("save withdraw error: %w", err)
	}
	return err
}

// GetWithdrawals история выводов средств, получение данных из БД PostgreSQL
// Принимает:
// - ctx: контекст с информацией о пользователе
// - userID: идентификатор пользователя
// Возвращает:
// - []*models.UserWithdrawals или ошибку, если возникли проблемы при получении данных
func (pg *Repository) GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]*models.UserWithdrawals, error) {
	query := `
	SELECT order_number, sum, processed_at FROM withdrawals WHERE user_id = $1 ORDER BY processed_at DESC
	`

	rows, err := pg.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get all withdrawals error: %w", err)
	}
	defer rows.Close()

	var withdrawals []*models.UserWithdrawals
	for rows.Next() {
		var w models.UserWithdrawals
		if err = rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt); err != nil {
			return nil, fmt.Errorf("get all withdrawals error: %w", err)
		}
		withdrawals = append(withdrawals, &w)
	}
	return withdrawals, err
}
