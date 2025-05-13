package converter

import (
	"lost-items-service/internal/model"
	modelredis "lost-items-service/internal/repository/user/redis/model"
)

func FromUserModelToRedis(user *model.User) modelredis.User {
	return modelredis.User{
		ID:       user.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Role:     user.Role,
	}
}

func FromModelRedisToUser(user *modelredis.User) *model.User {
	return &model.User{
		ID:       user.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Role:     user.Role,
	}
}
