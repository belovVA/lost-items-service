package redis

import (
	"context"

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
