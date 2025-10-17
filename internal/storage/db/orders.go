package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/IvanKondrashkov/go-market/internal/models"
	customError "github.com/IvanKondrashkov/go-market/internal/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetUserByOrderNumber получение пользователя по номеру заказа, получение данных из БД PostgreSQL
// Принимает:
// - ctx: контекст с информацией о пользователе
// - orderNumber: номер заказа
// Возвращает:
// - userID или ошибку, если не найден пользователь (ErrUserNotFound) или возникли проблемы при получении данных
func (pg *Repository) GetUserByOrderNumber(ctx context.Context, orderNumber string) (*uuid.UUID, error) {
	query := `
	SELECT user_id FROM orders WHERE number = $1
	`

	var userID *uuid.UUID
	err := pg.pool.QueryRow(ctx, query, orderNumber).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get user by order number error: %w", customError.ErrUserNotFound)
		}
		return nil, fmt.Errorf("get user by order number error: %w", err)
	}
	return userID, err
}

// GetOrders получение списка заказов пользователя, получение данных из БД PostgreSQL
// Принимает:
// - ctx: контекст с информацией о пользователе
// - userID: идентификатор пользователя
// Возвращает:
// - []*models.Order или ошибку, если возникли проблемы при получении данных
func (pg *Repository) GetOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	query := `
	SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC
	`

	rows, err := pg.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get all orders error: %w", err)
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err = rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			return orders, fmt.Errorf("get all orders error: %w", err)
		}
		orders = append(orders, &o)
	}
	return orders, err
}

// SaveOrder сохранение нового заказа, сохранение в БД PostgreSQL
// Принимает:
// - ctx: контекст с информацией о пользователе
// - tx: транзакцию
// - userID: идентификатор пользователя
// - orderNumber: номер заказа
// Возвращает:
// - ошибку, если возникли проблемы при сохранении
func (pg *Repository) SaveOrder(ctx context.Context, tx pgx.Tx, userID uuid.UUID, orderNumber string) error {
	query := `
	INSERT INTO orders(id, user_id, number, status, uploaded_at)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := tx.Exec(ctx, query, uuid.New(), userID, orderNumber, models.NEW, time.Now())
	if err != nil {
		return fmt.Errorf("save order error: %w", err)
	}
	return err
}

// UpdateOrderAccrual обновление заказа, обновление в БД PostgreSQL
// Принимает:
// - ctx: контекст с информацией о пользователе
// - tx: транзакцию
// - event: models.AccrualResponse
// Возвращает:
// - ошибку, если возникли проблемы при обновлении
func (pg *Repository) UpdateOrderAccrual(ctx context.Context, tx pgx.Tx, event models.AccrualResponse) error {
	query := `
	UPDATE orders SET status = $1, accrual = $2, processed_at = $3 WHERE number = $4
	`

	_, err := tx.Exec(ctx, query, event.Status, event.Accrual, time.Now(), event.Order)
	if err != nil {
		return fmt.Errorf("update order error: %w", err)
	}
	return err
}
