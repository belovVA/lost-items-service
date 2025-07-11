package converter

import (
	"lost-items-service/internal/model"
	modelpg "lost-items-service/internal/repository/user/pgdb/model"
)

func FromUserModelToRepo(user *model.User) *modelpg.User {
	return &modelpg.User{
		ID:       user.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Role:     user.Role,
	}
}

func FromModelRepoToUser(user *modelpg.User) *model.User {
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

func FromInfoUsersToLimitsUsers(info *model.InfoSetting) *modelpg.LimitsUsers {
	return &modelpg.LimitsUsers{
		Role:   info.OrderByField,
		Search: info.Search,
		Limit:  uint64(info.Limit),
		Offset: uint64((info.Page - 1) * info.Limit),
	}
}
