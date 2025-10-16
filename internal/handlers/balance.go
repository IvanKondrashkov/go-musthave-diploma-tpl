package handlers

import (
	"encoding/json"
	"net/http"
)

// GetBalance возвращает текущий баланс пользователя
// @Summary Get user balance
// @Description Получение текущего баланса счёта баллов лояльности пользователя
// @Tags balance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.UserBalance "Баланс пользователя"
// @Failure 401 {string} string "Пользователь не аутентифицирован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/balance [get]
func (app *App) GetBalance(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	balance, err := app.balanceService.GetBalance(req.Context())
	switch {
	case err == nil:
		res.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(res)
		err = enc.Encode(&balance)
		if err != nil {
			http.Error(res, "Invalid response format!", http.StatusBadRequest)
			return
		}
		return
	default:
		http.Error(res, "Internal server error!", http.StatusInternalServerError)
		return
	}
}
