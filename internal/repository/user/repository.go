package user

import (
	"context"
	"fmt"

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

// 1. GetUserByID: пробует Redis.GetUser, на промахе — Postgres + кеш.
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	// 1.1) Попытка из кэша
	if user, err := r.Redis.GetUser(ctx, id); err == nil {
		return user, nil
	}

	// 1.2) Cache-miss: дедупликация запросов к БД
	v, err, _ := r.group.Do(id.String(), func() (interface{}, error) {

		// 1.2.1) Читаем из Postgres
		user, err := r.Pg.GetUserByID(ctx, id)
		if err != nil {
			return nil, err
		}

		go func(u *model.User) {
			// 1.2.2) Кэшируем хешом
			if _, err = r.Redis.CreateUser(ctx, user); err != nil {
				// log.Warnf("redis HSET failed: %v", err)
			}

		}(user)

		return user, nil
	})

	if err != nil {
		return nil, err
	}
	return v.(*model.User), nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error) {
	// 2.1) Сохраняем в БД
	id, err := r.Pg.AddUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}
	user.ID = id

	// 2.2) Кэшируем хешом
	if _, err = r.Redis.CreateUser(ctx, user); err != nil {
		// log.Warnf("redis HSET failed: %v", err)
	}

	return id, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	// 3.1) Попытка из кэша по email
	if user, err := r.Redis.GetUserByEmail(ctx, email); err == nil {
		return user, nil
	}

	// 3.2) Cache-miss: дедупликация по строковому ключу email
	v, err, _ := r.group.Do(email, func() (interface{}, error) {
		// 3.2.1) Читаем из Postgres
		user, err := r.Pg.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		// 3.2.2) Кэшируем хешом

		go func(u *model.User) {
			if _, err = r.Redis.CreateUser(ctx, user); err != nil {
				// log.Warnf("redis HSET failed: %v", err)
			}

		}(user)

		return user, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*model.User), nil
}

func (r *UserRepository) GetUsers(ctx context.Context, limits *model.InfoUsers) ([]*model.User, error) {
	if users, err := r.Redis.GetUsers(ctx, limits); err == nil {
		return users, nil
	}
	// 3.2) Cache-miss: дедупликация по строковому ключу email
	groupKey := fmt.Sprintf("%s:%d:%d", limits.Role, limits.Page, limits.Limit)
	v, err, _ := r.group.Do(groupKey, func() (interface{}, error) {
		// 3.2.1) Читаем из Postgres
		users, err := r.Pg.GetListUsers(ctx, limits)
		if err != nil {
			return nil, err
		}
		// 3.2.2) Кэшируем хешом

		go func(u []*model.User) {
			if err = r.Redis.CreateUsersOrder(ctx, groupKey, u); err != nil {
				// log.Warnf("redis HSET failed: %v", err)

			}
		}(users)

		return users, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]*model.User), nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	// Обновляем в БД
	if err := r.Pg.UpdateUserByID(ctx, user); err != nil {
		return err
	}

	// Удаляем кэш по ID (основная запись)
	_ = r.Redis.Delete(ctx, user.ID.String())

	// Удаляем индекс по email
	_ = r.Redis.DeleteEmailIndex(ctx, user.Email)

	// Удаляем все связанные страницы кэша — нужен ключевой паттерн
	_ = r.Redis.DeleteUserPages(ctx, user.Role)

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, user *model.User) error {
	// Удаляем из Postgres
	if err := r.Pg.DeleteUserByID(ctx, user); err != nil {
		return err
	}

	// Чистим кэш
	_ = r.Redis.Delete(ctx, user.ID.String())
	_ = r.Redis.DeleteEmailIndex(ctx, user.Email)
	_ = r.Redis.DeleteUserPages(ctx, user.Role)

	return nil
}
