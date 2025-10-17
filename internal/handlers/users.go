package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/IvanKondrashkov/go-market/internal/models"
	"github.com/IvanKondrashkov/go-market/internal/service/middleware/auth"
	customError "github.com/IvanKondrashkov/go-market/internal/storage"
)

// Register регистрирует нового пользователя
// @Summary User registration
// @Description Регистрация нового пользователя в маркете «Гофермарт»
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserCredentials true "Данные для регистрации"
// @Success 200 {string} string "Пользователь успешно зарегистрирован и аутентифицирован"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 409 {string} string "Логин уже занят"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/register [post]
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

// Login аутентифицирует пользователя
// @Summary User login
// @Description Аутентификация пользователя в маркете «Гофермарт»
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserCredentials true "Данные для аутентификации"
// @Success 200 {string} string "Пользователь успешно аутентифицирован"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 401 {string} string "Неверная пара логин/пароль"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/login [post]
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
