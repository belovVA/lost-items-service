package redis

import (
	"context"

	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/redis/converter"
	modelredis "lost-items-service/internal/repository/user/redis/model"
)

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error) {
	var usr modelredis.User = converter.FromUserModelToRedis(user)
	idStr := usr.ID.String()

	if err := r.cl.HashSet(ctx, idStr, usr); err != nil {
		return uuid.Nil, err
	}

	return usr.ID, nil
}
