package user

import (
	"lost-items-service/internal/client/cache"
	"lost-items-service/internal/db/pgxdb"
	"lost-items-service/internal/repository/user/pgdb"
	"lost-items-service/internal/repository/user/redis"
)

type UserRepository struct {
	Pg    pgdb.UserPGRepository
	Redis redis.UserRedisRepository
}

func NewRepository(pg pgxdb.DB, redisClient cache.RedisClient) *UserRepository {
	return &UserRepository{
		Pg:    pgdb.NewRepository(pg),
		Redis: redis.NewRepository(redisClient),
	}
}
