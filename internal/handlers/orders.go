package handlers

import (
	"encoding/json"
	"errors"
	customContext "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/service/middleware/auth"
	"io"
	"net/http"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"
	customError "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/storage"
)

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
