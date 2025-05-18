package service

import (
	"context"

	"github.com/google/uuid"
	"lost-items-service/internal/model"
)

type AnnService struct {
	annRepository   AnnRepository
	imageRepository ImageRepository
}

func NewAnnService(repo AnnRepository) *AnnService {
	return &AnnService{
		annRepository: repo,
	}
}

func (s *AnnService) CreateAnnouncement(ctx context.Context, ann *model.Announcement) (uuid.UUID, error) {
	return s.annRepository.CreateAnn(ctx, ann)
}

func (s *AnnService) GetAnn(ctx context.Context, id uuid.UUID) (*model.Announcement, error) {
	return s.annRepository.GetAnnByID(ctx, id)
}
