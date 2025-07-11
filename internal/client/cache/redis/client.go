package redisclient

import (
	"context"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	cache "lost-items-service/internal/client/cache"
	"lost-items-service/internal/config"
)

var _ cache.RedisClient = (*client)(nil)

type handler func(ctx context.Context, conn redis.Conn) error

type client struct {
	pool   *redis.Pool
	config config.RedisConfig
}

func NewClient(pool *redis.Pool, config config.RedisConfig) *client {
	return &client{
		pool:   pool,
		config: config,
	}
}

func (c *client) HashSet(ctx context.Context, key string, values interface{}) error {
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, err := conn.Do("HSET", redis.Args{key}.AddFlat(values)...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Set(ctx context.Context, key string, value interface{}) error {
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, err := conn.Do("SET", redis.Args{key}.Add(value)...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) HGetAll(ctx context.Context, key string) ([]interface{}, error) {
	var values []interface{}
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		var errEx error
		values, errEx = redis.Values(conn.Do("HGETALL", key))
		if errEx != nil {
			return errEx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (c *client) Get(ctx context.Context, key string) (interface{}, error) {
	var value interface{}
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		var errEx error
		value, errEx = conn.Do("GET", key)
		if errEx != nil {
			return errEx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (c *client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, err := conn.Do("EXPIRE", key, int(expiration.Seconds()))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Ping(ctx context.Context) error {
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, err := conn.Do("PING")
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) execute(ctx context.Context, handler handler) error {
	conn, err := c.getConnect(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("failed to close redis connection: %v\n", err)
		}
	}()

	err = handler(ctx, conn)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) getConnect(ctx context.Context) (redis.Conn, error) {
	getConnTimeoutCtx, cancel := context.WithTimeout(ctx, c.config.ConnectionTimeout())
	defer cancel()

	conn, err := c.pool.GetContext(getConnTimeoutCtx)
	if err != nil {
		log.Printf("failed to get redis connection: %v\n", err)

		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}

func (c *client) Close() error {
	return c.pool.Close()
}

func (c *client) Del(ctx context.Context, key string) error {
	return c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, err := conn.Do("DEL", key)
		return err
	})
}

func (c *client) DeleteByPattern(ctx context.Context, pattern string) error {
	return c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		iter := 0
		for {
			reply, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern, "COUNT", 100))
			if err != nil {
				return err
			}

			var keys []string
			_, err = redis.Scan(reply, &iter, &keys)
			if err != nil {
				return err
			}

			if len(keys) > 0 {
				args := make([]interface{}, len(keys))
				for i, k := range keys {
					args[i] = k
				}
				_, _ = conn.Do("DEL", args...)
			}

			if iter == 0 {
				break
			}
		}
		return nil
	})
}
