package redis

import (
	"context"
	"encoding/json"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/redis/converter"
	modelredis "lost-items-service/internal/repository/user/redis/model"
)

func (r *userRepo) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	idStr := id.String()
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, model.ErrorUserNotFound
	}

	var user modelredis.User
	if err = redigo.ScanStruct(values, &user); err != nil {
		return nil, err
	}

	return converter.FromModelRedisToUser(&user), nil
}

// GetUserByEmail: сначала GET email→ID, затем HGETALL по ID
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	// 1) получаем ID по email
	raw, err := r.cl.Get(ctx, "user:email:"+email)
	if err != nil {
		return nil, err
	}
	idStr, ok := raw.(string)
	if !ok {
		return nil, fmt.Errorf("invalid ID type: %T", raw)
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	// 2) читаем самого пользователя по ID
	return r.GetUser(ctx, id)
}

func (r *userRepo) GetUsers(ctx context.Context, limits *model.InfoUsers) ([]*model.User, error) {
	cacheKey := fmt.Sprintf(
		"users:role:%s:page:%d:limit:%d",
		limits.Role, limits.Page, limits.Limit,
	)
	// 3) Попытка взять из Redis
	if raw, err := r.cl.Get(ctx, cacheKey); err == nil {
		// Ожидаем, что мы сохранили JSON-байты
		if data, ok := raw.([]byte); ok {
			var cached []*model.User
			if err := json.Unmarshal(data, &cached); err == nil {
				return cached, nil
			}
			// Если JSON битый — удалим ключ, чтобы не мусорить
			_ = r.cl.Del(ctx, cacheKey)
		}
	}
	return nil, model.ErrorUserNotFound
}
