package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel             string `env:"LOG_LEVEL"`
	AuthKey              string `env:"AUTH_KEY"`

	TerminationTimeout int `env:"TERMINATION_TIMEOUT"`
	WorkerCount        int `env:"WORKER_COUNT"`
}

var (
	RunAddress           = "localhost:8080"
	DatabaseURI          = ""
	AccrualSystemAddress = ""
	LogLevel             = "INFO"
	AuthKey              = []byte("6368616e676520746869732070617373776f726420746f206120736563726574")

	TerminationTimeout = time.Second * 30
	WorkerCount        = 10
)

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
