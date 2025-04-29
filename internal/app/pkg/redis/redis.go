package redis

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"lost-items-service/internal/config"
)

func InitRedisPool(ctx context.Context, cfg config.RedisConfig) (*redigo.Pool, error) {
	pool := &redigo.Pool{
		MaxIdle:     cfg.MaxIdle(),
		IdleTimeout: cfg.IdleTimeout(),
		DialContext: func(ctx context.Context) (redigo.Conn, error) {
			return redigo.DialContext(ctx, "tcp", cfg.Address())
		},
	}

	// Пытаемся сразу проверить соединение
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection from redis pool: %w", err)
	}
	defer conn.Close()

	_, err = conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return pool, nil
}
