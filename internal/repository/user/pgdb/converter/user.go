package converter

import (
	"lost-items-service/internal/model"
	modelRepo "lost-items-service/internal/repository/user/pgdb/model"
)

func FromUserModelToRepo(user *model.User) *modelRepo.User {
	return &modelRepo.User{
		ID:       user.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Role:     user.Role,
	}
}

func fromModelRepoToUser(user *modelRepo.User) *model.User {
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
