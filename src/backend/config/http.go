package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	Port        string        `yaml:"port"  env-default:"8080"`
	Host        string        `yaml:"host"  env-default:"localhost"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func HTTPConfigLoad() (*HTTPConfig, error) {
	cfgPath := os.Getenv(configPath)
	if len(cfgPath) == 0 {
		return nil, fmt.Errorf("%s environment not found", configPath)
	}

	if _, err := os.Stat(cfgPath); err != nil {
		return nil, fmt.Errorf("%s file not found", cfgPath)
	}

	var httpCfg HTTPConfig

	// Читаем конфиг-файл и заполняем нашу структуру
	if err := cleanenv.ReadConfig(configPath, &httpCfg); err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return &httpCfg, nil
}
