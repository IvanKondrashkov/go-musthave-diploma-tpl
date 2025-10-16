package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config содержит конфигурационные параметры приложения,
// которые могут быть установлены через переменные окружения.
type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`            // Адрес сервера в формате host:port
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"` // Адрес сервера системы лояльности в формате host:port
	DatabaseURI          string `env:"DATABASE_URI"`           // URI для подключения к БД
	LogLevel             string `env:"LOG_LEVEL"`              // Уровень логирования (DEBUG, INFO, WARN, ERROR)
	AuthKey              string `env:"AUTH_KEY"`               // Ключ для аутентификации

	TerminationTimeout int `env:"TERMINATION_TIMEOUT"` // Таймаут завершения работы (в секундах)
	WorkerCount        int `env:"WORKER_COUNT"`        // Количество воркеров
}

// Глобальные переменные конфигурации со значениями по умолчанию
var (
	RunAddress           = "localhost:8080"
	DatabaseURI          = ""
	AccrualSystemAddress = ""
	LogLevel             = "INFO"
	AuthKey              = []byte("6368616e676520746869732070617373776f726420746f206120736563726574")

	TerminationTimeout = time.Second * 30
	WorkerCount        = 10
)

// ParseConfig загружает конфигурацию приложения из:
// 1. Аргументов командной строки (имеют наивысший приоритет)
// 2. Переменных окружения
// 3. Значений по умолчанию
//
// Возвращает ошибку если не удалось распарсить конфигурацию.
func ParseConfig() error {
	flag.StringVar(&RunAddress, "a", RunAddress, "Base host host:port")
	flag.StringVar(&DatabaseURI, "d", DatabaseURI, "Base url db connection")
	flag.StringVar(&AccrualSystemAddress, "r", AccrualSystemAddress, "Base host accrual system")
	flag.Parse()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("config parse error: %w", err)
	}

	if envRunAddress := cfg.RunAddress; envRunAddress != "" {
		RunAddress = envRunAddress
	}

	if envDatabaseURI := cfg.DatabaseURI; envDatabaseURI != "" {
		DatabaseURI = envDatabaseURI
	}

	if envLogLevel := cfg.LogLevel; envLogLevel != "" {
		LogLevel = envLogLevel
	}

	if envAccrualSystemAddress := cfg.AccrualSystemAddress; envAccrualSystemAddress != "" {
		AccrualSystemAddress = envAccrualSystemAddress
	}

	if envAuthKey := cfg.AuthKey; envAuthKey != "" {
		AuthKey = []byte(envAuthKey)
	}

	if envTerminationTimeout := cfg.TerminationTimeout; envTerminationTimeout != 0 {
		TerminationTimeout = time.Duration(envTerminationTimeout)
	}

	if envWorkerCount := cfg.WorkerCount; envWorkerCount != 0 {
		WorkerCount = envWorkerCount
	}
	return err
}
