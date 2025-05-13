package redis

import (
	"context"

	"github.com/google/uuid"
	"lost-items-service/internal/client/cache"
	"lost-items-service/internal/model"
)

type UserRedisRepository interface {
	CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error)
	CreateUsersOrder(ctx context.Context, cachedKey string, users []*model.User) error

	GetUser(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUsers(ctx context.Context, limits *model.InfoUsers) ([]*model.User, error)

	SetEmailIndex(ctx context.Context, email string, id uuid.UUID) error

	Delete(ctx context.Context, key string) error
	DeleteUserPages(ctx context.Context, role string) error
	DeleteEmailIndex(ctx context.Context, email string) error
}

type userRepo struct {
	cl cache.RedisClient
}

func NewRepository(cl cache.RedisClient) UserRedisRepository {
	return &userRepo{cl: cl}
}
