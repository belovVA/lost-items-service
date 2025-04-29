package env

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"lost-items-service/internal/config"
)

const (
	redisHostEnvName              = "REDIS_HOST"
	redisPortEnvName              = "REDIS_PORT"
	redisConnectionTimeoutEnvName = "REDIS_CONNECTION_TIMEOUT_SEC"
	redisMaxIdleEnvName           = "REDIS_MAX_IDLE"
	redisIdleTimeoutEnvName       = "REDIS_IDLE_TIMEOUT_SEC"
)

type redisConfig struct {
	host string `env:"redis_host" env-required:"true"`
	port string `env:"redis_port" env-required:"true"`

	connectionTimeout time.Duration `yaml:"redis_connection_timeout_sec" env-default:"5s"`
	maxIdle           int           `yaml:"redis_max_idle" env-default:"5"`

	idleTimeout time.Duration `yaml:"redis_idle_timeout_sec" env-default:"60s"`
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

	if _, err = strconv.Atoi(redisCfg.port); err != nil {
		return nil, fmt.Errorf("invalid database redis port: %s", err)
	}

	return &redisCfg, nil
}

func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *redisConfig) ConnectionTimeout() time.Duration {
	return cfg.connectionTimeout
}

func (cfg *redisConfig) MaxIdle() int {
	return cfg.maxIdle
}

func (cfg *redisConfig) IdleTimeout() time.Duration {
	return cfg.idleTimeout
}
