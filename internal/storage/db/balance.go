package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"
	customError "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// SaveBalance сохранение баллов накопительного счёта, сохранение в БД PostgreSQL
// Принимает:
// - ctx: контекст с информацией о пользователе
// - tx: транзакцию
// - userID: идентификатор пользователя
// - event: models.AccrualResponse
// Возвращает:
// - ошибку, если возникли проблемы при сохранении
func (pg *Repository) SaveBalance(ctx context.Context, tx pgx.Tx, userID uuid.UUID, event models.AccrualResponse) error {
	query := `
	INSERT INTO balances(user_id, current) VALUES ($1, $2)
	ON CONFLICT (user_id)
	DO UPDATE SET current = balances.current + $2
	`

	_, err := tx.Exec(ctx, query, userID, event.Accrual)
	if err != nil {
		return fmt.Errorf("save balance error: %w", err)
	}
	return err
}

// GetBalance получение текущего баланса пользователя, получение данных из БД PostgreSQL
// Принимает:
// - ctx: контекст с информацией о пользователе
// - userID: идентификатор пользователя
// Возвращает:
// - *models.UserBalance или ошибку, если не найден баланс пользователя (ErrUserBalanceNotFound) или возникли проблемы при получении данных
func (pg *Repository) GetBalance(ctx context.Context, userID uuid.UUID) (*models.UserBalance, error) {
	query := `
	SELECT current, withdrawn FROM balances WHERE user_id = $1
	`

	var balance models.UserBalance
	err := pg.pool.QueryRow(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get balance error: %w", customError.ErrUserBalanceNotFound)
		}
		return nil, fmt.Errorf("get balance error: %w", err)
	}
	return &balance, err
}
