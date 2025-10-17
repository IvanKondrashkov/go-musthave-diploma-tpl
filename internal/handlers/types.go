package handlers

import (
	"context"
	"net/http"

	"github.com/IvanKondrashkov/go-market/internal/models"
	"github.com/IvanKondrashkov/go-market/internal/service/middleware/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/IvanKondrashkov/go-market/internal/config"
	"github.com/IvanKondrashkov/go-market/internal/logger"
	"github.com/IvanKondrashkov/go-market/internal/service"
	"github.com/IvanKondrashkov/go-market/internal/worker"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// UserService определяет интерфейс для работы с пользователями
type UserService interface {
	// Register регистрирует пользователя в системе
	Register(ctx context.Context, tx pgx.Tx, login, password string) (*uuid.UUID, error)
	// Authenticate аутентификация и авторизация пользователя в системе
	Authenticate(ctx context.Context, login, password string) (*uuid.UUID, error)
}

// OrderService интерфейс для работы с заказами
type OrderService interface {
	// SaveOrder сохранение нового заказа
	SaveOrder(ctx context.Context, tx pgx.Tx, userID uuid.UUID, orderNumber string) error
	// GetOrders получение списка заказов пользователя
	GetOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
}

// BalanceService интерфейс для работы с балансами
type BalanceService interface {
	// GetBalance получение текущего баланса пользователя
	GetBalance(ctx context.Context, userID uuid.UUID) (*models.UserBalance, error)
}

// WithdrawService интерфейс для работы со списаниями
type WithdrawService interface {
	// BalanceWithdraw списание баллов с накопительного счёта в счёт оплаты нового заказа
	BalanceWithdraw(ctx context.Context, tx pgx.Tx, userID uuid.UUID, withdraw models.UserWithdrawals) error
	// GetWithdrawals данные о выводе средств
	GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]*models.UserWithdrawals, error)
}

// App представляет основное приложение с сервисами и воркером
type App struct {
	URL             string
	worker          *worker.Worker
	userService     *service.UserService
	orderService    *service.OrderService
	balanceService  *service.BalanceService
	withdrawService *service.WithdrawService
}

// Handler обрабатывает HTTP-запросы
type Handler struct {
	Logger *logger.ZapLogger
	app    *App
}

// NewApp создает новый экземпляр App
func NewApp(newWorker *worker.Worker, userService *service.UserService, orderService *service.OrderService,
	balanceService *service.BalanceService, withdrawService *service.WithdrawService) *App {
	return &App{
		URL:             config.RunAddress,
		worker:          newWorker,
		userService:     userService,
		orderService:    orderService,
		balanceService:  balanceService,
		withdrawService: withdrawService,
	}
}

// NewHandler создает новый обработчик HTTP-запросов
func NewHandler(zl *logger.ZapLogger, app *App) *Handler {
	return &Handler{
		Logger: zl,
		app:    app,
	}
}

// NewRouter создает маршрутизатор с middleware и обработчиками
func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.Compress(3))
	r.Post("/api/user/register", h.app.Register)
	r.Post("/api/user/login", h.app.Login)

	r.Route(`/api/user`, func(r chi.Router) {
		r.Use(auth.Authentication)
		r.Post("/orders", h.app.SaveOrder)
		r.Get("/orders", h.app.GetOrders)
		r.Get("/balance", h.app.GetBalance)
		r.Post("/balance/withdraw", h.app.BalanceWithdraw)
		r.Get("/withdrawals", h.app.GetWithdrawals)
	})
	return r
}

// NewServer создает HTTP-сервер с настройками
func NewServer(r *chi.Mux) *http.Server {
	return &http.Server{
		Addr:         config.RunAddress,
		Handler:      r,
		ReadTimeout:  config.TerminationTimeout,
		WriteTimeout: config.TerminationTimeout,
	}
}
