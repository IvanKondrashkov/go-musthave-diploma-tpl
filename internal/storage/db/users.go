package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	customError "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (pg *Repository) Register(ctx context.Context, tx pgx.Tx, login, password string) (*uuid.UUID, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("register user error: %w", err)
	}

	query := `
	INSERT INTO users(id, login, password_hash)
	VALUES ($1, $2, $3) RETURNING id
	`

	var userID *uuid.UUID
	err = tx.QueryRow(ctx, query, uuid.New(), login, string(hashedPassword)).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("register user error: %w", customError.ErrLoginAlreadyExists)
	}
	return userID, err
}

func (pg *Repository) Authenticate(ctx context.Context, login, password string) (*uuid.UUID, error) {
	query := `
	SELECT id, password_hash FROM users WHERE login = $1
	`

	var userID uuid.UUID
	var passwordHash string
	err := pg.pool.QueryRow(ctx, query, login).Scan(&userID, &passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("authenticate user error: %w", customError.ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("authenticate user error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("authenticate user error: %w", customError.ErrInvalidCredentials)
	}
	return &userID, err
}
