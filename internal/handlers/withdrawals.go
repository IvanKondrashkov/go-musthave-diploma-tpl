package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/IvanKondrashkov/go-market/internal/models"
	customError "github.com/IvanKondrashkov/go-market/internal/storage"
)

// BalanceWithdraw списывает баллы с баланса
// @Summary Withdraw balance
// @Description Запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
// @Tags withdrawals
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param withdraw body models.UserWithdrawals true "Данные для списания"
// @Success 200 {string} string "Успешная обработка запроса"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 401 {string} string "Пользователь не аутентифицирован"
// @Failure 402 {string} string "На счету недостаточно средств"
// @Failure 422 {string} string "Неверный номер заказа"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/balance/withdraw [post]
func (app *App) BalanceWithdraw(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var withdraw models.UserWithdrawals
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&withdraw)
	if err != nil {
		http.Error(res, "Invalid request format!", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	err = app.withdrawService.BalanceWithdraw(req.Context(), withdraw)
	switch {
	case err == nil:
		res.WriteHeader(http.StatusOK)
		return
	case errors.Is(err, customError.ErrInvalidOrderNumber):
		http.Error(res, "Invalid order number!", http.StatusUnprocessableEntity)
		return
	case errors.Is(err, customError.ErrInvalidUserWithdraw):
		http.Error(res, "Are not enough funds in the account!", http.StatusPaymentRequired)
		return
	default:
		http.Error(res, "Internal server error!", http.StatusInternalServerError)
		return
	}
}

// GetWithdrawals возвращает историю выводов средств
// @Summary Get withdrawals history
// @Description Получение информации о выводе средств с накопительного счёта пользователем. Сортировка по времени вывода от новых к старым.
// @Tags withdrawals
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.UserWithdrawals "История выводов средств"
// @Success 204 {string} string "Нет ни одного списания"
// @Failure 401 {string} string "Пользователь не аутентифицирован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/withdrawals [get]
func (app *App) GetWithdrawals(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	withdrawals, err := app.withdrawService.GetWithdrawals(req.Context())
	switch {
	case err == nil:
		res.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(res)
		err = enc.Encode(&withdrawals)
		if err != nil {
			http.Error(res, "Invalid response format!", http.StatusBadRequest)
			return
		}
		return
	case errors.Is(err, customError.ErrWithdrawalsIsEmpty):
		http.Error(res, "Withdrawals is empty!", http.StatusNoContent)
		return
	default:
		http.Error(res, "Internal server error!", http.StatusInternalServerError)
		return
	}
}
