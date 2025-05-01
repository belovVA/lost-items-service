package user

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"
	"lost-items-service/internal/client/cache"
	"lost-items-service/internal/db/pgxdb"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/pgdb"
	"lost-items-service/internal/repository/user/redis"
)

type UserRepository struct {
	Pg    pgdb.UserPGRepository
	Redis redis.UserRedisRepository
	group singleflight.Group
}

func NewRepository(pg pgxdb.DB, redisClient cache.RedisClient) *UserRepository {
	return &UserRepository{
		Pg:    pgdb.NewRepository(pg),
		Redis: redis.NewRepository(redisClient),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error) {
	id, err := r.Pg.AddUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}
	_, _ = r.Redis.CreateUser(ctx, user)

	return id, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	//Сначала пробуем достать из Redis
	user, err := r.Redis.GetUser(ctx, id)
	if err == nil {
		return user, nil
	}
	v, err, _ := r.group.Do(id.String(), func() (interface{}, error) {
		// Достаем из Postgres
		user, err = r.Pg.UserByID(ctx, id)
		if err != nil {
			return nil, err
		}

		_, _ = r.Redis.CreateUser(ctx, user)

		return user, nil
	})

	if err != nil {
		return nil, err
	}

	return v.(*model.User), nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	v, err, _ := r.group.Do(email, func() (interface{}, error) {
		// Достаем из Postgres
		user, err := r.Pg.UserByEmail(ctx, email)
		if err != nil {
			return nil, err
		}

		_, _ = r.Redis.CreateUser(ctx, user)

		return user, nil
	})

	if err != nil {
		return nil, err
	}

	return v.(*model.User), nil
}
