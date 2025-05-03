package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (r *userRepo) SetEmailIndex(ctx context.Context, email string, id uuid.UUID) error {
	key := "user:email:" + email
	if err := r.cl.Set(ctx, key, id.String()); err != nil {
		return err
	}
	// Можно установить TTL (например, 10 минут)
	return r.cl.Expire(ctx, key, 10*time.Second)
}
