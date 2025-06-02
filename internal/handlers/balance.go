package handlers

import (
	"encoding/json"
	"net/http"
)

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
