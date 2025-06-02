package worker

import (
	"context"
	"sync"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/logger"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/service"
)

const (
	bufCh = 100
)

type Worker struct {
	wg           sync.WaitGroup
	orderService *service.OrderService
	resultCh     chan models.AccrualRequest
	errorCh      chan error
}

func NewWorker(workerCount int, zl *logger.ZapLogger, orderService *service.OrderService) *Worker {
	w := &Worker{
		orderService: orderService,
		resultCh:     make(chan models.AccrualRequest, bufCh),
		errorCh:      make(chan error, bufCh),
	}

	go w.ErrorListener(zl)

	for i := 0; i < workerCount; i++ {
		w.wg.Add(1)
		go w.RunJobUpdateOrderAccrual(context.Background())
	}
	return w
}
