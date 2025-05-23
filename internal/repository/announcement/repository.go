package announcement

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"
	"lost-items-service/internal/client/cache"
	"lost-items-service/internal/db/pgxdb"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/announcement/pgdb"
	"lost-items-service/internal/repository/announcement/redis"
)

type AnnouncementRepository struct {
	Pg    pgdb.AnnPGRepository
	Redis redis.AnnRedisRepository
	group *singleflight.Group
}

func NewRepository(pg pgxdb.DB, redisClient cache.RedisClient) *AnnouncementRepository {
	return &AnnouncementRepository{
		Pg:    pgdb.NewRepository(pg),
		Redis: redis.NewRepository(redisClient),
		group: &singleflight.Group{},
	}
}

func (r *AnnouncementRepository) CreateAnn(ctx context.Context, ann *model.Announcement) (uuid.UUID, error) {
	// 2.1) Сохраняем в БД
	id, err := r.Pg.AddAnn(ctx, ann)
	if err != nil {
		return uuid.Nil, err
	}
	ann.ID = id

	// 2.2) Кэшируем хешом
	if _, err = r.Redis.CreateAnn(ctx, ann); err != nil {
		// log.Warnf("redis HSET failed: %v", err)
	}

	return id, nil
}

func (r *AnnouncementRepository) GetAnnByID(ctx context.Context, id uuid.UUID) (*model.Announcement, error) {
	// 1.1) Попытка из кэша
	if ann, err := r.Redis.GetAnn(ctx, id); err == nil {
		return ann, nil
	}

	// 1.2) Cache-miss: дедупликация запросов к БД
	v, err, _ := r.group.Do(id.String(), func() (interface{}, error) {

		// 1.2.1) Читаем из Postgres
		ann, err := r.Pg.GetAnnByID(ctx, id)
		if err != nil {
			return nil, err
		}

		go func(u *model.Announcement) {
			// 1.2.2) Кэшируем хешом
			if _, err = r.Redis.CreateAnn(ctx, ann); err != nil {
				// log.Warnf("redis HSET failed: %v", err)
			}

		}(ann)

		return ann, nil
	})

	if err != nil {
		return nil, err
	}
	return v.(*model.Announcement), nil
}

func (r *AnnouncementRepository) GetAnnsList(ctx context.Context, info *model.InfoSetting) ([]*model.Announcement, error) {
	// 1.2) Cache-miss: дедупликация запросов к БД
	groupKey := fmt.Sprintf("%s:%s:%d:%d:%s", info.OrderByField, info.Search, info.Page, info.Limit, info.TimeOrder)
	v, err, _ := r.group.Do(groupKey, func() (interface{}, error) {
		// 1.2.1) Читаем из Postgres
		ann, err := r.Pg.GetListAnnouncement(ctx, info)
		if err != nil {
			return nil, err
		}

		//go func(u *model.Announcement) {
		//	// 1.2.2) Кэшируем хешом
		//	if _, err = r.Redis.CreateAnn(ctx, ann); err != nil {
		//		// log.Warnf("redis HSET failed: %v", err)
		//	}
		//
		//}(ann)

		return ann, nil
	})

	if err != nil {
		return nil, err
	}
	return v.([]*model.Announcement), nil
}

func (r *AnnouncementRepository) GetUserAnns(ctx context.Context, userID uuid.UUID, info *model.InfoSetting) ([]*model.Announcement, error) {
	groupKey := fmt.Sprintf("%s:%s:%d:%d:%s", info.OrderByField, info.Search, info.Page, info.Limit, info.TimeOrder)
	v, err, _ := r.group.Do(groupKey, func() (interface{}, error) {
		// 1.2.1) Читаем из Postgres
		ann, err := r.Pg.GetListAnnouncement(ctx, info)
		if err != nil {
			return nil, err
		}

		//go func(u *model.Announcement) {
		//	// 1.2.2) Кэшируем хешом
		//	if _, err = r.Redis.CreateAnn(ctx, ann); err != nil {
		//		// log.Warnf("redis HSET failed: %v", err)
		//	}
		//
		//}(ann)

		return ann, nil
	})

	if err != nil {
		return nil, err
	}
	return v.([]*model.Announcement), nil
}

func (r *AnnouncementRepository) UpdateAnnouncement(ctx context.Context, ann *model.Announcement) error {
	// 1.2) Cache-miss: дедупликация запросов к БД
	_, err, _ := r.group.Do("update"+ann.ID.String(), func() (interface{}, error) {

		// 1.2.1) Читаем из Postgres
		if err := r.Pg.UpdateFields(ctx, ann); err != nil {
			return nil, err
		}

		//go func(u *model.Announcement) {
		//	// 1.2.2) Кэшируем хешом
		//	if err := r.Redis.Delete(ctx, ann.ID); err != nil {
		//		// log.Warnf("redis HSET failed: %v", err)
		//	}
		//
		//	if _, err := r.Redis.CreateAnn(ctx, ann); err != nil {
		//		// log.Warnf("redis HSET failed: %v", err)
		//	}
		//
		//}(ann)

		return nil, nil
	})

	return err
}

func (r *AnnouncementRepository) DeleteAnnByID(ctx context.Context, id uuid.UUID) error {
	// 1.2) Cache-miss: дедупликация запросов к БД
	_, err, _ := r.group.Do("delete"+id.String(), func() (interface{}, error) {

		// 1.2.1) Читаем из Postgres
		err := r.Pg.DeleteAnn(ctx, id)
		if err != nil {
			return nil, err
		}

		go func(id uuid.UUID) {
			// 1.2.2) Кэшируем хешом
			if err = r.Redis.Delete(ctx, id); err != nil {
				// log.Warnf("redis HSET failed: %v", err)
			}

		}(id)

		return nil, nil
	})

	return err
}

func (r *AnnouncementRepository) UpdateModerationStatusAnnouncement(ctx context.Context, ann *model.Announcement) error {
	// 1.2) Cache-miss: дедупликация запросов к БД
	_, err, _ := r.group.Do("update"+ann.ID.String(), func() (interface{}, error) {

		// 1.2.1) Читаем из Postgres
		if err := r.Pg.UpdateModerationStatus(ctx, ann); err != nil {
			return nil, err
		}

		//go func(u *model.Announcement) {
		//	// 1.2.2) Кэшируем хешом
		//	if err := r.Redis.Delete(ctx, ann.ID); err != nil {
		//		// log.Warnf("redis HSET failed: %v", err)
		//	}
		//
		//	if _, err := r.Redis.CreateAnn(ctx, ann); err != nil {
		//		// log.Warnf("redis HSET failed: %v", err)
		//	}
		//
		//}(ann)

		return nil, nil
	})

	return err
}
