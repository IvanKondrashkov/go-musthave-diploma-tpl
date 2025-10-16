package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/config"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/handlers"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/logger"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/service"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/storage/db"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/worker"

	"go.uber.org/zap"
)

// @title Gophermart API
// @version 1.0
// @description маркет «Гофермарт»

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name Authorization

// @tag.name auth "Аутентификация и регистрация"
// @tag.name orders "Заказы"
// @tag.name balance "Баланс"
// @tag.name withdrawals "Списания"
func main() {
	err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	zl, err := logger.NewZapLogger(config.LogLevel)
	if err != nil {
		return err
	}
	defer zl.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), config.TerminationTimeout)
	defer cancel()

	newRepository, err := db.NewRepository(ctx, zl, config.DatabaseURI)
	if err != nil {
		return err
	}
	defer newRepository.Close()
	newRunner := newRepository

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		newUserService := service.NewUserService(zl, newRunner, newRepository)
		newBalanceService := service.NewBalanceService(zl, newRunner, newRepository)
		newWithdrawService := service.NewWithdrawService(zl, newRunner, newRepository)
		newOrderService := service.NewOrderService(zl, newRunner, newRepository, newBalanceService, newWithdrawService)
		newWorker := worker.NewWorker(ctx, config.WorkerCount, zl, newOrderService)
		newApp := handlers.NewApp(newWorker, newUserService, newOrderService, newBalanceService, newWithdrawService)
		newHandler := handlers.NewHandler(zl, newApp)
		newRouter := handlers.NewRouter(newHandler)
		newHTTPServer := handlers.NewServer(newRouter)

		defer newWorker.Close()

		return runServer(zl, newHTTPServer)
	}
}

func runServer(zl *logger.ZapLogger, httpServer *http.Server) error {
	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		zl.Log.Info("HTTP server starting", zap.String("address", config.RunAddress))
		errChan <- httpServer.ListenAndServe()
	}()

	select {
	case sig := <-sigChan:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), config.TerminationTimeout)
		defer cancel()

		zl.Log.Info("Received signal, shutting down gracefully", zap.String("signal", sig.String()))
		zl.Log.Info("Stopping HTTP server...")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			zl.Log.Error("Server shutdown failed", zap.Error(err))
			return err
		}

		zl.Log.Info("Server stopped gracefully")
		return nil

	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			zl.Log.Error("Server error", zap.Error(err))
			return err
		}
		return nil
	}
}
