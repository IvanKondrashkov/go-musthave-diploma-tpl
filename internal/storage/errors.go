package storage

import "errors"

var (
	ErrLoginAlreadyExists       = errors.New("login already exists")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrInvalidToken             = errors.New("invalid token")
	ErrInvalidSigningMethod     = errors.New("invalid signing method")
	ErrInvalidOrderNumber       = errors.New("invalid order number")
	ErrInvalidUserWithdraw      = errors.New("invalid user withdraw")
	ErrAlreadyExistsOrderNumber = errors.New("already exists order number")
	ErrUserNotFound             = errors.New("user not found")
	ErrUserBalanceNotFound      = errors.New("user balance not found")
	ErrOrdersIsEmpty            = errors.New("orders is empty found")
	ErrWithdrawalsIsEmpty       = errors.New("withdrawals is empty found")
)
