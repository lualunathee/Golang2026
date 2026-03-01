package modules

import (
	"os"
	"time"
)

type PostgreConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string

	ExecTimeout time.Duration
}

func LoadPostgresConfig() *PostgreConfig {
	cfg := &PostgreConfig{
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		Username:    os.Getenv("DB_USER"),
		Password:    os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		SSLMode:     os.Getenv("DB_SSLMODE"),
		ExecTimeout: 5 * time.Second,
	}

	// дефолт, чтобы не падало
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	return cfg
}
