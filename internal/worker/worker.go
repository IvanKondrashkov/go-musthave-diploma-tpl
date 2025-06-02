package worker

import (
	"context"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/logger"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"

	"go.uber.org/zap"
)

func (w *Worker) SendOrderAccrualRequest(ctx context.Context, event models.AccrualRequest) {
	select {
	case w.resultCh <- event:
	case <-ctx.Done():
		return
	}
}

func (w *Worker) RunJobUpdateOrderAccrual(ctx context.Context) {
	defer w.wg.Done()
	for event := range w.resultCh {
		if err := w.orderService.CreateOrderAccrual(event); err != nil {
			w.errorCh <- err
			continue
		}

		if err := w.orderService.UpdateOrderAccrual(ctx, event); err != nil {
			w.errorCh <- err
		}
	}
}

func (w *Worker) ErrorListener(zl *logger.ZapLogger) {
	for err := range w.errorCh {
		zl.Log.Debug("order update error", zap.Error(err))
	}
}

func (w *Worker) Close() {
	close(w.resultCh)
	w.wg.Wait()
	close(w.errorCh)
}
