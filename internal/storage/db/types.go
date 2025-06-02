package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/logger"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/service"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	service.UserRepository
	service.OrderRepository
	service.BalanceRepository
	service.WithdrawRepository
	Logger *logger.ZapLogger
	pool   *pgxpool.Pool
}

func NewRepository(ctx context.Context, zl *logger.ZapLogger, dns string) (*Repository, error) {
	parseConfig, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, parseConfig)
	if err != nil {
		return nil, fmt.Errorf("open database connection error: %w", err)
	}

	m, err := migrate.New("file://migration", dns)
	if err != nil {
		return nil, fmt.Errorf("database migration error: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("database migration error: %w", err)
	}

	return &Repository{
		Logger: zl,
		pool:   pool,
	}, nil
}

func (pg *Repository) Close() {
	pg.pool.Close()
}
