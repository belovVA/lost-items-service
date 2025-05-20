package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"lost-items-service/internal/converter"
	"lost-items-service/internal/middleware"
	"lost-items-service/internal/model"
	"lost-items-service/internal/service/pkg/hash"
)

type UserService struct {
	userRepository UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		userRepository: repo,
	}
}

func (s *UserService) GetOwnUser(ctx context.Context) (*model.User, error) {
	userID, err := converter.ToUUIDFromStringID(ctx.Value(middleware.UserIDKey).(string))
	if err != nil {
		return nil, fmt.Errorf("failed to get uuid")
	}
	return s.userRepository.GetUserByID(ctx, userID)
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepository.GetUserByID(ctx, id)
}

func (s *UserService) InfoUsers(ctx context.Context, limits *model.InfoSetting) ([]*model.User, error) {
	return s.userRepository.GetUsers(ctx, limits)
}

func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	usr, err := s.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if usr.Email == user.Email && usr.ID != user.ID {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	if user.Password != "" {
		hashPass, err := hash.HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("failed to hash pass")
		}
		user.Password = hashPass
	}
	return s.userRepository.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, user *model.User) error {
	if _, err := s.userRepository.GetUserByID(ctx, user.ID); err != nil {
		return err
	}
	return s.userRepository.DeleteUser(ctx, user)
}
