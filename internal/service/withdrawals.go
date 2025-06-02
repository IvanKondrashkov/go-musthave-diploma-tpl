package service

import (
	"context"
	"fmt"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"
	customContext "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/service/middleware/auth"
	customError "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/storage"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/utils"
)

func (s *WithdrawService) BalanceWithdraw(ctx context.Context, withdraw models.UserWithdrawals) error {
	tx, err := s.Runner.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("open transactional error: %w", err)
	}

	if !utils.IsValidLuna(withdraw.Order) {
		return fmt.Errorf("invalid order number luna error: %w", customError.ErrInvalidOrderNumber)
	}

	userID := customContext.GetContextUserID(ctx)
	err = s.WithdrawRepository.BalanceWithdraw(ctx, tx, *userID, withdraw)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("withdraw balance error: %w", err)
	}
	return tx.Commit(ctx)
}

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
