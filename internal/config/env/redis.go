package env

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"lost-items-service/internal/config"
)

type redisConfig struct {
	Host string `env:"REDIS_HOST" env-required:"true"`
	Port string `env:"REDIS_PORT" env-required:"true"`

	ConnectionTimeoutRedis time.Duration `yaml:"redis_connection_timeout_sec" env-default:"5s"`
	MaxIdleRedis           int           `yaml:"redis_max_idle" env-default:"5"`

	IdleTimeoutRedis time.Duration `yaml:"redis_idle_timeout_sec" env-default:"60s"`
}

func RedisConfigLoad(configPath string) (*redisConfig, error) {
	path, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	var redisCfg redisConfig

	// Читаем конфиг-файл и заполняем нашу структуру
	if err := cleanenv.ReadConfig(path, &redisCfg); err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if _, err = strconv.Atoi(redisCfg.Port); err != nil {
		return nil, fmt.Errorf("invalid database redis port: %s", err)
	}

	return &redisCfg, nil
}

func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}

func (cfg *redisConfig) ConnectionTimeout() time.Duration {
	return cfg.ConnectionTimeoutRedis
}

func (cfg *redisConfig) MaxIdle() int {
	return cfg.MaxIdleRedis
}

func (cfg *redisConfig) IdleTimeout() time.Duration {
	return cfg.IdleTimeoutRedis
}
