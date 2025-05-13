package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/redis/converter"
)

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error) {
	var usr = converter.FromUserModelToRedis(user)
	idStr := usr.ID.String()

	if err := r.cl.HashSet(ctx, idStr, usr); err != nil {
		return uuid.Nil, err
	}

	// Индекс по email
	if err := r.SetEmailIndex(ctx, user.Email, user.ID); err != nil {
		// Можно залогировать, но не падать
		//log.Printf("failed to set email index: %v\n", err)
	}

	// TTL на основную запись (опционально)
	_ = r.cl.Expire(ctx, idStr, 10*time.Minute)

	return usr.ID, nil
}

func (r *userRepo) CreateUsersOrder(ctx context.Context, cachedKey string, users []*model.User) error {
	if data, err := json.Marshal(users); err == nil {
		// SET cacheKey data EX 300
		if err = r.cl.Set(ctx, cachedKey, data); err != nil {
			return err
		}
		_ = r.cl.Expire(ctx, cachedKey, 1*time.Minute)
	} else {
		// логируем, что не смогли закэшировать
		//log.Printf("failed to marshal users for cache: %v", err)
	}
	return nil
}
