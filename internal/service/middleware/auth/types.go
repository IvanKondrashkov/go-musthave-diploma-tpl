package auth

import (
	"fmt"
	"time"

	"github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/config"
	customError "github.com/IvanKondrashkov/go-musthave-diploma-tpl/internal/storage"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	AuthCookie      string = "auth"
	tokenExpiration        = 24 * time.Hour
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

func ValidateToken(tokenString string) (*uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("validate token error: %w", customError.ErrInvalidSigningMethod)
		}
		return config.AuthKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate token error: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &claims.UserID, nil
	}

	return nil, fmt.Errorf("validate token error: %w", customError.ErrInvalidToken)
}

func GenerateToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(tokenExpiration)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.AuthKey)
}
