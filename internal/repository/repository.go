package repository

import (
	"lost-items-service/internal/client/cache"
	"lost-items-service/internal/db/pgxdb"
	"lost-items-service/internal/repository/announcement"
	"lost-items-service/internal/repository/user"
)

type Repository struct {
	*user.UserRepository
	*announcement.AnnouncementRepository
}

func NewRepository(pg pgxdb.DB, redisClient cache.RedisClient) *Repository {
	return &Repository{
		UserRepository:         user.NewRepository(pg, redisClient),
		AnnouncementRepository: announcement.NewRepository(pg, redisClient),
	}
}
