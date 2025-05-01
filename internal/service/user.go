package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"lost-items-service/internal/converter"
	"lost-items-service/internal/middleware"
	"lost-items-service/internal/model"
)

type UserService struct {
	userRepository UserRepository
}

func NewUserService(repo UserRepository) *AuthService {
	return &AuthService{
		userRepository: repo,
	}
}

func (s *UserService) GetOwnUser(ctx context.Context) (*model.User, error) {
	userID, err := converter.ToUUIDFromStringID(ctx.Value(middleware.UserIDKey).(string))
	if err != nil {
		return nil, fmt.Errorf("failed to get uuid")
	}
	log.Println("aba")
	return s.userRepository.GetUserByID(ctx, userID)
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepository.GetUserByID(ctx, id)
}
