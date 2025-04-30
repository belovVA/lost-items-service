package redis

import (
	"context"

	"github.com/google/uuid"
	"lost-items-service/internal/client/cache"
	"lost-items-service/internal/model"
)

type UserRedisRepository interface {
	CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type userRepo struct {
	cl cache.RedisClient
}

func NewRepository(cl cache.RedisClient) UserRedisRepository {
	return &userRepo{cl: cl}
}
