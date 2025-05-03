package converter

import (
	"github.com/google/uuid"
	"lost-items-service/internal/handler/dto"
	"lost-items-service/internal/model"
)

func ToUserFromCreateUserRequest(request *dto.RegisterRequest) *model.User {
	return &model.User{
		ID:       uuid.Nil,
		Name:     request.Name,
		Surname:  request.Surname,
		Email:    request.Email,
		Phone:    request.Phone,
		Password: request.Password,
		Role:     request.Role,
	}
}

func ToRegisterResponseFromUser(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
	}
}

func ToUserFromLoginRequest(request *dto.LoginRequest) *model.User {
	return &model.User{
		ID:       uuid.Nil,
		Email:    request.Email,
		Password: request.Password,
	}
}

func ToUUIDFromStringID(idStr string) (uuid.UUID, error) {
	return uuid.Parse(idStr)
}
