package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"lost-items-service/internal/config"
)

var (
	ErrConnectionFailed = errors.New("postgres: connection failed")
)

func InitDBPool(ctx context.Context, cfg config.PGConfig) (*pgxpool.Pool, error) {
	dsn := cfg.DSN()
	bgCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.Connect(bgCtx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnectionFailed, err)
	}

	// Проверим, что соединение установлено
	if err := pingDBPool(ctx, pool); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnectionFailed, err)
	}

	return pool, nil
}

func pingDBPool(ctx context.Context, pool *pgxpool.Pool) error {
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		return fmt.Errorf("%w: %s", ErrConnectionFailed, err)
	}
	return nil
}
