package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type PGConfig struct {
	Name     string `yaml:"database_name" env:"DATABASE_NAME" required:"true"`
	Host     string `yaml:"database_host" env:"DATABASE_HOST" required:"true"`
	Port     string `yaml:"database_port" env:"DATABASE_PORT" required:"true"`
	User     string `yaml:"database_user" env:"DATABASE_USER" required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
}

func PGConfigLoad() (*PGConfig, error) {
	cfgPath := os.Getenv(configPath)
	if len(cfgPath) == 0 {
		return nil, fmt.Errorf("%s environment not found", configPath)
	}

	if _, err := os.Stat(cfgPath); err != nil {
		return nil, fmt.Errorf("%s file not found", cfgPath)
	}

	var pgCfg PGConfig

	// Читаем конфиг-файл и заполняем нашу структуру
	if err := cleanenv.ReadConfig(configPath, &pgCfg); err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return &pgCfg, nil
}
