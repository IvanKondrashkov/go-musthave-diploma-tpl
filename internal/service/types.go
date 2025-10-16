package service

import (
	"context"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/logger"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Runner интерфейс для работы с транзакциями
type Runner interface {
	// BeginTx начинает новую транзакцию
	BeginTx(ctx context.Context) (pgx.Tx, error)
	// Commit коммитит изменения
	Commit(ctx context.Context, tx pgx.Tx) error
	// Rollback откатывает транзакцию
	Rollback(ctx context.Context, tx pgx.Tx) error
}

// UserRepository интерфейс для пользовательских операций с данными
type UserRepository interface {
	Runner
	// Register регистрирует пользователя в системе
	Register(ctx context.Context, tx pgx.Tx, login, password string) (*uuid.UUID, error)
	// Authenticate аутентификация и авторизация пользователя в системе
	Authenticate(ctx context.Context, login, password string) (*uuid.UUID, error)
	// GenerateToken генерация JWT токена доступа
	GenerateToken(userID uuid.UUID) (string, error)
}

// OrderRepository интерфейс для операций с заказами
type OrderRepository interface {
	Runner
	// SaveOrder сохранение нового заказа
	SaveOrder(ctx context.Context, tx pgx.Tx, userID uuid.UUID, orderNumber string) error
	// GetOrders получение списка заказов пользователя
	GetOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
	// GetUserByOrderNumber получение пользователя по номеру заказа
	GetUserByOrderNumber(ctx context.Context, orderNumber string) (*uuid.UUID, error)
	// UpdateOrderAccrual обновление заказа
	UpdateOrderAccrual(ctx context.Context, tx pgx.Tx, event models.AccrualResponse) error
}

// BalanceRepository интерфейс для операций с балансом
type BalanceRepository interface {
	Runner
	// SaveBalance сохранение баллов накопительного счёта
	SaveBalance(ctx context.Context, tx pgx.Tx, userID uuid.UUID, event models.AccrualResponse) error
	// GetBalance получение текущего баланса пользователя
	GetBalance(ctx context.Context, userID uuid.UUID) (*models.UserBalance, error)
}

// WithdrawRepository интерфейс для операций со списанием балов накопительной системы
type WithdrawRepository interface {
	Runner
	// BalanceWithdraw списание баллов с накопительного счёта в счёт оплаты нового заказа
	BalanceWithdraw(ctx context.Context, tx pgx.Tx, userID uuid.UUID, withdraw models.UserWithdrawals) error
	// GetWithdrawals данные о выводе средств
	GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]*models.UserWithdrawals, error)
}

// UserService реализует бизнес-логику сервиса пользователя
type UserService struct {
	Logger         *logger.ZapLogger // Логгер для записи событий
	UserRepository UserRepository    // Репозиторий для работы с пользователями
	Runner         Runner            // Для работы с транзакциями
}

// BalanceService реализует бизнес-логику сервиса балансов
type BalanceService struct {
	Logger            *logger.ZapLogger // Логгер для записи событий
	BalanceRepository BalanceRepository // Репозиторий для работы с балансами
	Runner            Runner            // Для работы с транзакциями
}

// WithdrawService реализует бизнес-логику сервиса списаний
type WithdrawService struct {
	Logger             *logger.ZapLogger  // Логгер для записи событий
	WithdrawRepository WithdrawRepository // Репозиторий для работы со списаниями
	Runner             Runner             // Для работы с транзакциями
}

// OrderService реализует бизнес-логику сервиса заказов
type OrderService struct {
	Logger          *logger.ZapLogger // Логгер для записи событий
	OrderRepository OrderRepository   // Сервис для работы с пользователями
	BalanceService  *BalanceService   // Сервис для работы с балансами
	WithdrawService *WithdrawService  // Сервис для работы со списаниями
	Runner          Runner            // Для работы с транзакциями
}

// NewUserService создает новый экземпляр сервиса пользователей
// Принимает:
// - zl: логгер
// - ru: реализация интерфейса Runner
// - userRepository: реализация интерфейса UserRepository
// Возвращает инициализированный UserService
func NewUserService(zl *logger.ZapLogger, ru Runner, userRepository UserRepository) *UserService {
	return &UserService{
		Logger:         zl,
		Runner:         ru,
		UserRepository: userRepository,
	}
}

// NewWithdrawService создает новый экземпляр сервиса списаний
// Принимает:
// - zl: логгер
// - ru: реализация интерфейса Runner
// - withdrawRepository: реализация интерфейса WithdrawRepository
// Возвращает инициализированный WithdrawService
func NewWithdrawService(zl *logger.ZapLogger, ru Runner, withdrawRepository WithdrawRepository) *WithdrawService {
	return &WithdrawService{
		Logger:             zl,
		Runner:             ru,
		WithdrawRepository: withdrawRepository,
	}
}

// NewBalanceService создает новый экземпляр сервиса балансов
// Принимает:
// - zl: логгер
// - ru: реализация интерфейса Runner
// - balanceRepository: реализация интерфейса BalanceRepository
// Возвращает инициализированный BalanceService
func NewBalanceService(zl *logger.ZapLogger, ru Runner, balanceRepository BalanceRepository) *BalanceService {
	return &BalanceService{
		Logger:            zl,
		Runner:            ru,
		BalanceRepository: balanceRepository,
	}
}

// NewOrderService создает новый экземпляр сервиса заказов
// Принимает:
// - zl: логгер
// - ru: реализация интерфейса Runner
// - orderRepository: реализация интерфейса OrderRepository
// - balanceService: сервис балансов BalanceService
// - withdrawService: сервис списаний WithdrawService
// Возвращает инициализированный OrderService
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
