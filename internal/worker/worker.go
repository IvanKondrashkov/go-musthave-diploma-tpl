package worker

import (
	"context"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/logger"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"

	"go.uber.org/zap"
)

// SendOrderAccrualRequest отправляет задачу на пакетное обновление заказов в очередь обработки
// Принимает:
// ctx - контекст для контроля времени выполнения
// event - models.AccrualRequest
func (w *Worker) SendOrderAccrualRequest(ctx context.Context, event models.AccrualRequest) {
	select {
	case w.resultCh <- event:
	case <-ctx.Done():
		return
	}
}

// RunJobUpdateOrderAccrual запускает воркер для обработки задач обновления заказов
// Принимает:
// ctx - контекст для контроля времени выполнения
func (w *Worker) RunJobUpdateOrderAccrual(ctx context.Context) {
	defer w.wg.Done()
	for event := range w.resultCh {
		if err := w.orderService.CreateOrderAccrual(event); err != nil && ctx.Err() == nil {
			w.errorCh <- err
			continue
		}

		if err := w.orderService.UpdateOrderAccrual(ctx, event); err != nil && ctx.Err() == nil {
			w.errorCh <- err
		}
	}
}

// ErrorListener обрабатывает ошибки от воркеров
// Принимает:
// ctx - контекст для контроля времени выполнения
// zl - логгер для записи ошибок
// Возвращает канал для сигнализации завершения ErrorListener
func (w *Worker) ErrorListener(ctx context.Context, zl *logger.ZapLogger) <-chan struct{} {
	defer close(w.doneCh)

	for err := range w.errorCh {
		select {
		case <-ctx.Done():
			zl.Log.Debug("order update error (shutdown)", zap.Error(err))
		default:
			zl.Log.Debug("order update error", zap.Error(err))
		}
	}
	return w.doneCh
}

// Close останавливает воркеры и освобождает ресурсы
func (w *Worker) Close() {
	close(w.resultCh)
	w.wg.Wait()
	close(w.errorCh)

	<-w.doneCh
}
