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

func FromModelRepoToUser(user *modelRepo.User) *model.User {
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

func FromInfoUsersToLimitsUsers(info *model.InfoUsers) *modelRepo.LimitsUsers {
	return &modelRepo.LimitsUsers{
		Role:   info.Role,
		Limit:  uint64(info.Limit),
		Offset: uint64((info.Page - 1) * info.Limit),
	}
}
