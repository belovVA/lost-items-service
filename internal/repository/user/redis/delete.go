package redis

import (
	"context"
	"fmt"
)

func (r *userRepo) Delete(ctx context.Context, key string) error {

	return r.cl.Del(ctx, key)
}

func (r *userRepo) DeleteUserPages(ctx context.Context, role string) error {
	pattern := fmt.Sprintf("users:role:%s:*", role)
	return r.cl.DeleteByPattern(ctx, pattern)
}

func (r *userRepo) DeleteEmailIndex(ctx context.Context, email string) error {
	cacheKey := fmt.Sprintf("user:email:%s", email)
	return r.cl.Del(ctx, cacheKey)
}
