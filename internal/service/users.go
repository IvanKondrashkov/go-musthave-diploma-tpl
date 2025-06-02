package service

import (
	"context"
	"fmt"
	customError "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/storage"

	"github.com/google/uuid"
)

func (s *UserService) Register(ctx context.Context, login, password string) (*uuid.UUID, error) {
	tx, err := s.Runner.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("open transactional error: %w", err)
	}

	userID, err := s.UserRepository.Register(ctx, tx, login, password)
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, fmt.Errorf("register user error: %w", customError.ErrLoginAlreadyExists)
	}
	return userID, tx.Commit(ctx)
}

func (s *UserService) Authenticate(ctx context.Context, login, password string) (*uuid.UUID, error) {
	userID, err := s.UserRepository.Authenticate(ctx, login, password)
	if err != nil {
		return nil, fmt.Errorf("authenticate user error: %w", customError.ErrInvalidCredentials)
	}
	return userID, err
}
