package storage

import "errors"

// Пакет storage содержит определения ошибок хранилища
var (
	// ErrLoginAlreadyExists - возникает при попытке создать дублирующую сущность
	ErrLoginAlreadyExists = errors.New("login already exists")
	// ErrInvalidCredentials - возникает при передаче невалидных кредов пользователя
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInvalidToken - возникает при передаче невалидного токена авторизации
	ErrInvalidToken = errors.New("invalid token")
	// ErrInvalidSigningMethod - возникает при передаче невалидного метода подписи токена
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	// ErrInvalidOrderNumber - возникает при передаче невалидного номера заказа
	ErrInvalidOrderNumber = errors.New("invalid order number")
	// ErrInvalidUserWithdraw - возникает при недостаточном кол-ве денег для списания
	ErrInvalidUserWithdraw = errors.New("invalid user withdraw")
	// ErrAlreadyExistsOrderNumber - возникает при попытке создать дублирующую сущность
	ErrAlreadyExistsOrderNumber = errors.New("already exists order number")
	// ErrUserNotFound - возникает при попытке получить сущность которая не существует
	ErrUserNotFound = errors.New("user not found")
	// ErrUserBalanceNotFound - возникает при попытке получить сущность которая не существует
	ErrUserBalanceNotFound = errors.New("user balance not found")
	// ErrOrdersIsEmpty - возникает при попытке получить сущность которая пустая
	ErrOrdersIsEmpty = errors.New("orders is empty found")
	// ErrWithdrawalsIsEmpty - возникает при попытке получить сущность которая пустая
	ErrWithdrawalsIsEmpty = errors.New("withdrawals is empty found")
)
