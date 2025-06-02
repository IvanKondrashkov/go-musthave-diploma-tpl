package service

import (
	"context"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/logger"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Runner interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type UserRepository interface {
	Runner
	Register(ctx context.Context, tx pgx.Tx, login, password string) (*uuid.UUID, error)
	Authenticate(ctx context.Context, login, password string) (*uuid.UUID, error)
	GenerateToken(userID uuid.UUID) (string, error)
}

type OrderRepository interface {
	Runner
	SaveOrder(ctx context.Context, tx pgx.Tx, userID uuid.UUID, orderNumber string) error
	GetOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
	GetUserByOrderNumber(ctx context.Context, orderNumber string) (*uuid.UUID, error)
	UpdateOrderAccrual(ctx context.Context, tx pgx.Tx, event models.AccrualResponse) error
}

type BalanceRepository interface {
	Runner
	SaveBalance(ctx context.Context, tx pgx.Tx, userID uuid.UUID, event models.AccrualResponse) error
	GetBalance(ctx context.Context, userID uuid.UUID) (*models.UserBalance, error)
}

type WithdrawRepository interface {
	Runner
	BalanceWithdraw(ctx context.Context, tx pgx.Tx, userID uuid.UUID, withdraw models.UserWithdrawals) error
	GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]*models.UserWithdrawals, error)
}

type UserService struct {
	Logger         *logger.ZapLogger
	UserRepository UserRepository
	Runner         Runner
}

type BalanceService struct {
	Logger            *logger.ZapLogger
	BalanceRepository BalanceRepository
	Runner            Runner
}

type WithdrawService struct {
	Logger             *logger.ZapLogger
	WithdrawRepository WithdrawRepository
	Runner             Runner
}

type OrderService struct {
	Logger          *logger.ZapLogger
	OrderRepository OrderRepository
	BalanceService  *BalanceService
	WithdrawService *WithdrawService
	Runner          Runner
}

func NewUserService(zl *logger.ZapLogger, ru Runner, userRepository UserRepository) *UserService {
	return &UserService{
		Logger:         zl,
		Runner:         ru,
		UserRepository: userRepository,
	}
}

func NewWithdrawService(zl *logger.ZapLogger, ru Runner, withdrawRepository WithdrawRepository) *WithdrawService {
	return &WithdrawService{
		Logger:             zl,
		Runner:             ru,
		WithdrawRepository: withdrawRepository,
	}
}

func NewBalanceService(zl *logger.ZapLogger, ru Runner, balanceRepository BalanceRepository) *BalanceService {
	return &BalanceService{
		Logger:            zl,
		Runner:            ru,
		BalanceRepository: balanceRepository,
	}
}

func NewOrderService(zl *logger.ZapLogger, ru Runner,
	orderRepository OrderRepository, balanceService *BalanceService, withdrawService *WithdrawService) *OrderService {
	return &OrderService{
		Logger:          zl,
		Runner:          ru,
		OrderRepository: orderRepository,
		BalanceService:  balanceService,
		WithdrawService: withdrawService,
	}
}
