package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/IvanKondrashkov/go-market/internal/models"
	customContext "github.com/IvanKondrashkov/go-market/internal/service/middleware/auth"
	customError "github.com/IvanKondrashkov/go-market/internal/storage"
)

// SaveOrder загружает номер заказа для расчета
// @Summary Upload order number
// @Description Загрузка номера заказа для расчета баллов лояльности. Номер заказа проверяется алгоритмом Луна.
// @Tags orders
// @Accept plain
// @Produce plain
// @Security ApiKeyAuth
// @Param order body string true "Номер заказа" Example("12345678903")
// @Success 200 {string} string "Заказ уже был загружен этим пользователем"
// @Success 202 {string} string "Новый номер заказа принят в обработку"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 401 {string} string "Пользователь не аутентифицирован"
// @Failure 409 {string} string "Номер заказа уже был загружен другим пользователем"
// @Failure 422 {string} string "Неверный формат номера заказа"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/orders [post]
func (app *App) SaveOrder(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Invalid request format!", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	order := string(body)
	userID, err := app.orderService.SaveOrder(req.Context(), order)
	switch {
	case err == nil:
		if userID != nil {
			res.WriteHeader(http.StatusOK)
		} else {
			userID = customContext.GetContextUserID(req.Context())
			event := models.AccrualRequest{
				UserID: *userID,
				Order:  order,
				Goods:  make([]models.Goods, 0),
			}

			go app.worker.SendOrderAccrualRequest(req.Context(), event)
			res.WriteHeader(http.StatusAccepted)
		}
		return
	case errors.Is(err, customError.ErrInvalidOrderNumber):
		http.Error(res, "Invalid order number!", http.StatusUnprocessableEntity)
		return
	case errors.Is(err, customError.ErrAlreadyExistsOrderNumber):
		http.Error(res, "Already exists order number for another user!", http.StatusConflict)
		return
	default:
		http.Error(res, "Internal server error!", http.StatusInternalServerError)
		return
	}
}

// GetOrders возвращает список загруженных номеров заказов
// @Summary Get user orders
// @Description Получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях.
// Сортировка по времени загрузки от новых к старым.
// @Tags orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Order "Список заказов"
// @Success 204 {string} string "Нет данных для ответа"
// @Failure 401 {string} string "Пользователь не аутентифицирован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/orders [get]
func (app *App) GetOrders(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	orders, err := app.orderService.GetOrders(req.Context())
	switch {
	case err == nil:
		res.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(res)
		err = enc.Encode(&orders)
		if err != nil {
			http.Error(res, "Invalid response format!", http.StatusBadRequest)
			return
		}
		return
	case errors.Is(err, customError.ErrOrdersIsEmpty):
		http.Error(res, "Orders is empty!", http.StatusNoContent)
		return
	default:
		http.Error(res, "Internal server error!", http.StatusInternalServerError)
		return
	}
}
