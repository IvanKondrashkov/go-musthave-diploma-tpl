package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"
	customContext "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/service/middleware/auth"
)

func (s *BalanceService) SaveBalance(ctx context.Context, userID uuid.UUID, event models.AccrualResponse) error {
	tx, err := s.Runner.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("open transactional error: %w", err)
	}

	err = s.BalanceRepository.SaveBalance(ctx, tx, userID, event)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("save balance error: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *BalanceService) GetBalance(ctx context.Context) (*models.UserBalance, error) {
	userID := customContext.GetContextUserID(ctx)
	balance, err := s.BalanceRepository.GetBalance(ctx, *userID)
	if err != nil {
		return nil, fmt.Errorf("get balance error: %w", err)
	}
	return balance, err
}
