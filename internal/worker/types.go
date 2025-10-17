package worker

import (
	"context"
	"sync"

	"github.com/IvanKondrashkov/go-market/internal/logger"
	"github.com/IvanKondrashkov/go-market/internal/models"
	"github.com/IvanKondrashkov/go-market/internal/service"
)

const (
	bufCh = 100
)

// Worker - структура для фоновой обработки задач обновления заказов
type Worker struct {
	wg           sync.WaitGroup             // Группа ожидания завершения воркеров
	orderService *service.OrderService      // Сервис для операций с заказами
	resultCh     chan models.AccrualRequest // Канал для задач обновления заказов
	errorCh      chan error                 // Канал для ошибок
	doneCh       chan struct{}              // Канал для сигнализации завершения ErrorListener
}

// NewWorker создает новый пул воркеров для обработки обновления заказов
// Принимает:
// - ctx: контекст для контроля времени выполнения
// - workerCount: количество воркеров
// - zl: логгер
// - orderService: сервис для операций с заказами
// Возвращает инициализированный Worker
func NewWorker(ctx context.Context, workerCount int, zl *logger.ZapLogger, orderService *service.OrderService) *Worker {
	w := &Worker{
		orderService: orderService,
		resultCh:     make(chan models.AccrualRequest, bufCh),
		errorCh:      make(chan error, bufCh),
	}

	go w.ErrorListener(ctx, zl)

	for i := 0; i < workerCount; i++ {
		w.wg.Add(1)
		go w.RunJobUpdateOrderAccrual(ctx)
	}
	return w
}
