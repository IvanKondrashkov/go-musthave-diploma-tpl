package handlers

import (
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/service/middleware/auth"
	"net/http"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/config"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/logger"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/service"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/worker"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type UserService interface {
	Register(res http.ResponseWriter, req *http.Request)
	Login(res http.ResponseWriter, req *http.Request)
}

type OrderService interface {
	SaveOrder(res http.ResponseWriter, req *http.Request)
	GetOrders(res http.ResponseWriter, req *http.Request)
}

type BalanceService interface {
	GetBalance(res http.ResponseWriter, req *http.Request)
}

type WithdrawService interface {
	BalanceWithdraw(res http.ResponseWriter, req *http.Request)
	GetWithdrawals(res http.ResponseWriter, req *http.Request)
}

type App struct {
	URL             string
	worker          *worker.Worker
	userService     *service.UserService
	orderService    *service.OrderService
	balanceService  *service.BalanceService
	withdrawService *service.WithdrawService
}

type Handler struct {
	Logger *logger.ZapLogger
	app    *App
}

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

func NewHandler(zl *logger.ZapLogger, app *App) *Handler {
	return &Handler{
		Logger: zl,
		app:    app,
	}
}

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

func NewServer(r *chi.Mux) *http.Server {
	return &http.Server{
		Addr:         config.RunAddress,
		Handler:      r,
		ReadTimeout:  config.TerminationTimeout,
		WriteTimeout: config.TerminationTimeout,
	}
}
