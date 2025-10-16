package service

import (
	"context"
	"fmt"

	customError "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/storage"

	"github.com/google/uuid"
)

// Register регистрация пользователя
// Принимает:
// - ctx: контекст с информацией о пользователе
// - login: логин пользователя
// - password: пароль пользователя в зашифрованном виде
// Возвращает:
// - userID или ошибку, если пользователь уже существует (ErrLoginAlreadyExists) или возникли проблемы при сохранении
func (s *UserService) Register(ctx context.Context, login, password string) (*uuid.UUID, error) {
	tx, err := s.Runner.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("open transactional error: %w", err)
	}

	userID, err := s.UserRepository.Register(ctx, tx, login, password)
	if err != nil {
		_ = s.Runner.Rollback(ctx, tx)
		return nil, fmt.Errorf("register user error: %w", customError.ErrLoginAlreadyExists)
	}
	return userID, s.Runner.Commit(ctx, tx)
}

// Authenticate аутентификация и авторизация пользователя
// Принимает:
// - ctx: контекст с информацией о пользователе
// - login: логин пользователя
// - password: пароль пользователя
// Возвращает:
// - userID или ошибку, если данные не валидны (ErrInvalidCredentials) или возникли проблемы при авторизации
func (s *UserService) Authenticate(ctx context.Context, login, password string) (*uuid.UUID, error) {
	userID, err := s.UserRepository.Authenticate(ctx, login, password)
	if err != nil {
		return nil, fmt.Errorf("authenticate user error: %w", customError.ErrInvalidCredentials)
	}
	return userID, err
}
