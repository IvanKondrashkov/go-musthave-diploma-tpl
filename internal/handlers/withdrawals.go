package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"
	customError "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/storage"
)

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
