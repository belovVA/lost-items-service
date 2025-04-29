package repository

import (
	"lost-items-service/internal/client/cache"
	"lost-items-service/internal/db/pgxdb"
	"lost-items-service/internal/repository/user"
)

type Repository struct {
	*user.UserRepository
	//Ad   *ad.Repository
}

func NewRepository(pg pgxdb.DB, redisClient cache.RedisClient) *Repository {
	return &Repository{
		UserRepository: user.NewRepository(pg, redisClient),
		//Ad:   ad.NewRepository(pg, redisClient),
	}
}
