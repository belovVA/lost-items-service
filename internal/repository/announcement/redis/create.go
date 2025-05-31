package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/announcement/redis/converter"
)

func (r *annRepo) CreateAnn(ctx context.Context, ann *model.Announcement) (uuid.UUID, error) {
	a := converter.FromAnnModelToRedis(ann)

	idStr := a.ID

	if err := r.cl.HashSet(ctx, idStr, a); err != nil {
		return uuid.Nil, err
	}
	_ = r.cl.Expire(ctx, idStr, 10*time.Minute)

	return uuid.Parse(idStr)
}
