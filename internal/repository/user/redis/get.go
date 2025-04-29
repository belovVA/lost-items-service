package redis

import (
	"context"

	"github.com/google/uuid"
	"lost-items-service/internal/model"
)

func (r *userRepo) UserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return nil, nil
}
