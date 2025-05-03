package service

import (
	"context"

	"github.com/google/uuid"
	"lost-items-service/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUsers(ctx context.Context, limits *model.InfoUsers) ([]*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, user *model.User) error
}

type Repository interface {
	UserRepository
}

type Service struct {
	*AuthService
	*UserService
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		AuthService: &AuthService{repo, jwtSecret},
		UserService: &UserService{repo},
	}
}
