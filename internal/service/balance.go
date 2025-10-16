package service

import (
	"context"
	"fmt"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"
	customContext "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/service/middleware/auth"

	"github.com/google/uuid"
)

// SaveBalance сохранение баллов накопительного счёта
// Принимает:
// - ctx: контекст с информацией о пользователе
// - userID: идентификатор пользователя
// - event: models.AccrualResponse
// Возвращает:
// - ошибку, если возникли проблемы при сохранении
func (s *BalanceService) SaveBalance(ctx context.Context, userID uuid.UUID, event models.AccrualResponse) error {
	tx, err := s.Runner.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("open transactional error: %w", err)
	}

	err = s.BalanceRepository.SaveBalance(ctx, tx, userID, event)
	if err != nil {
		_ = s.Runner.Rollback(ctx, tx)
		return fmt.Errorf("save balance error: %w", err)
	}
	return s.Runner.Commit(ctx, tx)
}

// GetBalance получение текущего баланса пользователя, получение данных из БД PostgreSQL
// Принимает:
// - ctx: контекст с информацией о пользователе
// Возвращает:
// - *models.UserBalance или ошибку, если возникли проблемы при получении данных
func (s *BalanceService) GetBalance(ctx context.Context) (*models.UserBalance, error) {
	userID := customContext.GetContextUserID(ctx)
	balance, err := s.BalanceRepository.GetBalance(ctx, *userID)
	if err != nil {
		return nil, fmt.Errorf("get balance error: %w", err)
	}
	return balance, err
}
