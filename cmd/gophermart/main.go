package main

import (
	"context"
	"log"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/config"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/handlers"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/logger"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/service"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/storage/db"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/worker"

	"go.uber.org/zap"
)

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
		newWorker := worker.NewWorker(config.WorkerCount, zl, newOrderService)
		newApp := handlers.NewApp(newWorker, newUserService, newOrderService, newBalanceService, newWithdrawService)
		newHandler := handlers.NewHandler(zl, newApp)
		newRouter := handlers.NewRouter(newHandler)
		newServer := handlers.NewServer(newRouter)

		defer newWorker.Close()

		zl.Log.Info("Running server", zap.String("address", config.RunAddress))
		return newServer.ListenAndServe()
	}
}
