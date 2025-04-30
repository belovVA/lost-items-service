package service

import (
	"context"

	"github.com/google/uuid"
	"lost-items-service/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error)
	UserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UserByEmail(ctx context.Context, email string) (*model.User, error)
}

type Repository interface {
	UserRepository
}

type Service struct {
	*AuthService
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		AuthService: &AuthService{repo, jwtSecret},
	}
}
