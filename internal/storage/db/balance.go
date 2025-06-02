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
