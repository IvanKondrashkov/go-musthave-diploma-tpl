package auth

import (
	"net/http"
)

// Authentication middleware проверяет/устанавливает аутентификацию пользователя
// Принимает:
// h - следующий обработчик в цепочке
// Возвращает обработчик с проверкой аутентификации
func Authentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AuthCookie)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := ValidateToken(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := SetContextUserID(r.Context(), *userID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
