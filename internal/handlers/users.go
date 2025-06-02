package handlers

import (
	"encoding/json"
	"errors"
	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/service/middleware/auth"
	"net/http"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/models"
	customError "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/storage"
)

func (app *App) Register(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var creds models.UserCredentials
	if err := json.NewDecoder(req.Body).Decode(&creds); err != nil {
		http.Error(res, "Invalid request format!", http.StatusBadRequest)
		return
	}

	userID, err := app.userService.Register(req.Context(), creds.Login, creds.Password)
	switch {
	case errors.Is(err, customError.ErrLoginAlreadyExists):
		http.Error(res, "Login already exists!", http.StatusConflict)
		return
	case err != nil:
		http.Error(res, "Internal server error!", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(*userID)
	if err != nil {
		http.Error(res, "Internal server error!", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:  auth.AuthCookie,
		Value: token,
		Path:  "/",
	})
	res.WriteHeader(http.StatusOK)
}

func (app *App) Login(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var credentials models.UserCredentials
	if err := json.NewDecoder(req.Body).Decode(&credentials); err != nil {
		http.Error(res, "Invalid request format!", http.StatusBadRequest)
		return
	}

	userID, err := app.userService.Authenticate(req.Context(), credentials.Login, credentials.Password)
	switch {
	case errors.Is(err, customError.ErrInvalidCredentials):
		http.Error(res, "Invalid credentials!", http.StatusUnauthorized)
		return
	case err != nil:
		http.Error(res, "Internal server error!", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(*userID)
	if err != nil {
		http.Error(res, "Internal server error!", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:  auth.AuthCookie,
		Value: token,
		Path:  "/",
	})
	res.WriteHeader(http.StatusOK)
}
