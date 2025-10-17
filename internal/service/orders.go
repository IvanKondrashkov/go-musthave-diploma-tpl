package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/IvanKondrashkov/go-market/internal/config"
	"github.com/IvanKondrashkov/go-market/internal/models"
	customContext "github.com/IvanKondrashkov/go-market/internal/service/middleware/auth"
	customError "github.com/IvanKondrashkov/go-market/internal/storage"
	"github.com/IvanKondrashkov/go-market/internal/utils"

	"github.com/google/uuid"
)

// SaveOrder сохранение нового заказа
// Принимает:
// - ctx: контекст с информацией о пользователе
// - orderNumber: номер заказа
// Возвращает:
// - userID или ошибку, если возникли проблемы при сохранении
func (s *OrderService) SaveOrder(ctx context.Context, orderNumber string) (*uuid.UUID, error) {
	tx, err := s.Runner.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("open transactional error: %w", err)
	}

	if !utils.IsValidLuna(orderNumber) {
		_ = s.Runner.Rollback(ctx, tx)
		return nil, fmt.Errorf("invalid order number luna error: %w", customError.ErrInvalidOrderNumber)
	}

	existsUserID, err := s.OrderRepository.GetUserByOrderNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, customError.ErrUserNotFound) {
		_ = s.Runner.Rollback(ctx, tx)
		return nil, fmt.Errorf("get user by order number error: %w", err)
	}

	userID := customContext.GetContextUserID(ctx)
	if existsUserID != nil && *existsUserID != *userID {
		_ = s.Runner.Rollback(ctx, tx)
		return nil, fmt.Errorf("order number already exists error: %w", customError.ErrAlreadyExistsOrderNumber)
	}

	if existsUserID != nil && *existsUserID == *userID {
		return userID, s.Runner.Commit(ctx, tx)
	}

	err = s.OrderRepository.SaveOrder(ctx, tx, *userID, orderNumber)
	if err != nil {
		_ = s.Runner.Rollback(ctx, tx)
		return nil, fmt.Errorf("save order error: %w", err)
	}
	return nil, s.Runner.Commit(ctx, tx)
}

// GetOrders получение списка заказов пользователя
// Принимает:
// - ctx: контекст с информацией о пользователе
// Возвращает:
// - []*models.Order или ошибку, если массив пуст (ErrOrdersIsEmpty) или возникли проблемы при получении данных
func (s *OrderService) GetOrders(ctx context.Context) ([]*models.Order, error) {
	userID := customContext.GetContextUserID(ctx)
	orders, err := s.OrderRepository.GetOrders(ctx, *userID)
	if err != nil {
		return nil, fmt.Errorf("get all orders error: %w", err)
	}

	if len(orders) == 0 {
		return orders, fmt.Errorf("get all orders error: %w", customError.ErrOrdersIsEmpty)
	}
	return orders, err
}

// CreateOrderAccrual отправка заказа в систему лояльности
// Принимает:
// - event: models.AccrualRequest
// Возвращает:
// - ошибку, если возникли проблемы при сохранении данных
func (s *OrderService) CreateOrderAccrual(event models.AccrualRequest) error {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	if err := encoder.Encode(event); err != nil {
		return fmt.Errorf("accrual encoding error: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("%s/api/orders", config.AccrualSystemAddress), &buf)
	if err != nil {
		return fmt.Errorf("request order accrual error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send order accrual error: %w", err)
	}
	defer res.Body.Close()
	return err
}

// UpdateOrderAccrual обновление заказа, данными полученными из системы лояльности
// Принимает:
// - ctx: контекст с информацией о пользователе
// - event: models.AccrualRequest
// Возвращает:
// - ошибку, если возникли проблемы при обновлении данных
func (s *OrderService) UpdateOrderAccrual(ctx context.Context, event models.AccrualRequest) error {
	tx, err := s.Runner.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("open transactional error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/orders/%s", config.AccrualSystemAddress, event.Order), http.NoBody)
	if err != nil {
		_ = s.Runner.Rollback(ctx, tx)
		return fmt.Errorf("request order accrual error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		_ = s.Runner.Rollback(ctx, tx)
		return fmt.Errorf("recive order accrual error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		_ = s.Runner.Rollback(ctx, tx)
		return fmt.Errorf("accrual returned status %d", res.StatusCode)
	}

	var respDto models.AccrualResponse
	if err = json.NewDecoder(res.Body).Decode(&respDto); err != nil {
		_ = s.Runner.Rollback(ctx, tx)
		return fmt.Errorf("accrual decoding error: %w", err)
	}

	err = s.BalanceService.SaveBalance(ctx, event.UserID, respDto)
	if err != nil {
		_ = s.Runner.Rollback(ctx, tx)
		return fmt.Errorf("save balance accrual error: %w", err)
	}

	err = s.OrderRepository.UpdateOrderAccrual(ctx, tx, respDto)
	if err != nil {
		_ = s.Runner.Rollback(ctx, tx)
		return fmt.Errorf("update order accrual error: %w", err)
	}
	return s.Runner.Commit(ctx, tx)
}
