package redis

import (
	"context"

	"github.com/google/uuid"
	"lost-items-service/internal/client/cache"
	"lost-items-service/internal/model"
)

type AnnRedisRepository interface {
	CreateAnn(ctx context.Context, ann *model.Announcement) (uuid.UUID, error)
	GetAnn(ctx context.Context, id uuid.UUID) (*model.Announcement, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type annRepo struct {
	cl cache.RedisClient
}

func NewRepository(cl cache.RedisClient) AnnRedisRepository {
	return &annRepo{cl: cl}
}
