package env

import (
	"fmt"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
	"lost-items-service/internal/config"
)

type pgConfig struct {
	Name     string `yaml:"database_name" env-required:"true"`
	Host     string `env:"DATABASE_HOST" env-required:"true"`
	Port     string `yaml:"database_port" env-required:"true"`
	User     string `yaml:"database_user" env-required:"true"`
	SSLMode  string `yaml:"database_ssl_mode" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
}

func PGConfigLoad(configPath string) (*pgConfig, error) {
	path, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	var pgCfg pgConfig

	// Читаем конфиг-файл и заполняем нашу структуру
	if err = cleanenv.ReadConfig(path, &pgCfg); err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if _, err = strconv.Atoi(pgCfg.Port); err != nil {
		return nil, fmt.Errorf("invalid database port: %s", err)
	}

	return &pgCfg, nil
}

func (cfg *pgConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
	)
}
