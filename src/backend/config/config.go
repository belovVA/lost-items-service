package config

import (
	"fmt"
	"os"
	"time"
)

const (
	configPath = "cfgPath"
)

type PGConfig interface {
	GetDSN() string
}

type HTTPConfig interface {
	GetPort() string
	GetHost() string
	GetTimeout() time.Duration
	GetIdleTimeout() time.Duration
}

type JWTConfig interface {
	GetSecret() string
}

func LoadConfig() (string, error) {
	cfgPath := os.Getenv(configPath)
	if len(cfgPath) == 0 {
		return "", fmt.Errorf("%s environment not found", configPath)
	}

	if _, err := os.Stat(cfgPath); err != nil {
		return "", fmt.Errorf("%s file not found", cfgPath)
	}

	return cfgPath, nil
}
