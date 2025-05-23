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
	GetUsers(ctx context.Context, limits *model.InfoSetting) ([]*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, user *model.User) error
}

type AnnRepository interface {
	CreateAnn(ctx context.Context, ann *model.Announcement) (uuid.UUID, error)
	GetAnnByID(ctx context.Context, id uuid.UUID) (*model.Announcement, error)
	GetAnnsList(ctx context.Context, info *model.InfoSetting) ([]*model.Announcement, error)
	GetUserAnns(ctx context.Context, userID uuid.UUID, info *model.InfoSetting) ([]*model.Announcement, error)
	UpdateAnnouncement(ctx context.Context, ann *model.Announcement) error
	UpdateModerationStatusAnnouncement(ctx context.Context, ann *model.Announcement) error
	DeleteAnnByID(ctx context.Context, id uuid.UUID) error
}

type ImageRepository interface {
	CreateImage(ctx context.Context, image *model.Image) (uuid.UUID, error)
	GetImagesByAnnouncementID(ctx context.Context, annID uuid.UUID) ([]*model.Image, error)
}

type Repository interface {
	UserRepository
	AnnRepository
	ImageRepository
}

type Service struct {
	*AuthService
	*UserService
	*AnnService
	*ImageService
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		AuthService:  NewAuthService(repo, jwtSecret),
		UserService:  NewUserService(repo),
		AnnService:   NewAnnService(repo, repo),
		ImageService: NewImageService(repo, repo),
	}
}
