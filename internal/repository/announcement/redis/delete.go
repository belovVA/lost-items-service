package redis

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *annRepo) Delete(ctx context.Context, id uuid.UUID) error {
	cacheKey := fmt.Sprintf("announcement:%s", id.String())
	return r.cl.Del(ctx, cacheKey)
}
