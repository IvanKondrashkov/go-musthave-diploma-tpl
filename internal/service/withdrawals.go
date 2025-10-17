package service

import (
	"context"
	"fmt"

	"github.com/IvanKondrashkov/go-market/internal/models"
	customContext "github.com/IvanKondrashkov/go-market/internal/service/middleware/auth"
	customError "github.com/IvanKondrashkov/go-market/internal/storage"
	"github.com/IvanKondrashkov/go-market/internal/utils"
)

// BalanceWithdraw списание баллов с накопительного счёта в счёт оплаты нового заказа
// Принимает:
// - ctx: контекст с информацией о пользователе
// - withdraw: models.UserWithdrawals
// Возвращает:
// - ошибку, если номер заказа не валиден (ErrInvalidOrderNumber) или возникли проблемы при сохранении
func (s *WithdrawService) BalanceWithdraw(ctx context.Context, withdraw models.UserWithdrawals) error {
	tx, err := s.Runner.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("open transactional error: %w", err)
	}

	if !utils.IsValidLuna(withdraw.Order) {
		_ = s.Runner.Rollback(ctx, tx)
		return fmt.Errorf("invalid order number luna error: %w", customError.ErrInvalidOrderNumber)
	}

	userID := customContext.GetContextUserID(ctx)
	err = s.WithdrawRepository.BalanceWithdraw(ctx, tx, *userID, withdraw)
	if err != nil {
		_ = s.Runner.Rollback(ctx, tx)
		return fmt.Errorf("withdraw balance error: %w", err)
	}
	return s.Runner.Commit(ctx, tx)
}

// GetWithdrawals история выводов средств
// Принимает:
// - ctx: контекст с информацией о пользователе
// Возвращает:
// - []*models.UserWithdrawals или ошибку, если массив пустой (ErrWithdrawalsIsEmpty) или возникли проблемы при получении данных
func (s *WithdrawService) GetWithdrawals(ctx context.Context) ([]*models.UserWithdrawals, error) {
	userID := customContext.GetContextUserID(ctx)
	withdraws, err := s.WithdrawRepository.GetWithdrawals(ctx, *userID)
	if err != nil {
		return nil, fmt.Errorf("get all withdrawals error: %w", err)
	}

	if len(withdraws) == 0 {
		return withdraws, fmt.Errorf("get all withdrawals error: %w", customError.ErrWithdrawalsIsEmpty)
	}
	return withdraws, err
}
