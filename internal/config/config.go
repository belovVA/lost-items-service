package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type PGConfig interface {
	DSN() string
}

type HTTPConfig interface {
	GetPort() string
	GetHost() string
	GetTimeout() time.Duration
	GetIdleTimeout() time.Duration
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type JWTConfig interface {
	GetSecret() string
}

func LoadConfig(cfgPath string) (string, error) {
	if _, err := os.Stat(cfgPath); err != nil {
		return "", fmt.Errorf("%s file not found", cfgPath)
	}
	return cfgPath, nil
}

func LoadEnv(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}
