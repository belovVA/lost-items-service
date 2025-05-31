package redis

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/announcement/redis/converter"
	modelredis "lost-items-service/internal/repository/announcement/redis/model"
)

func (r *annRepo) GetAnn(ctx context.Context, id uuid.UUID) (*model.Announcement, error) {
	idStr := id.String()
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrorNotFound
	}

	var ann modelredis.Announcement
	if err = redigo.ScanStruct(values, &ann); err != nil {
		return nil, err
	}

	return converter.FromRedisToAnnModel(&ann)
}
