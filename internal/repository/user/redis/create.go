package redis

import (
	"context"

	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/pgdb/converter"
)

func (r *userRepo) CreateUser(ctx context.Context, model *model.User) (uuid.UUID, error) {
	user := converter.FromUserModelToRepo(model)

	return user.ID, nil
}
