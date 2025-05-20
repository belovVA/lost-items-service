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

func NewAnnService(repoAn AnnRepository, repoImg ImageRepository) *AnnService {
	return &AnnService{
		annRepository:   repoAn,
		imageRepository: repoImg,
	}
}

func (s *AnnService) CreateAnnouncement(ctx context.Context, ann *model.Announcement) (uuid.UUID, error) {
	return s.annRepository.CreateAnn(ctx, ann)
}

func (s *AnnService) GetAnn(ctx context.Context, id uuid.UUID) (*model.Announcement, error) {
	ann, err := s.annRepository.GetAnnByID(ctx, id)
	if err != nil {
		return nil, err
	}

	ann.Images, err = s.imageRepository.GetImagesByAnnouncementID(ctx, ann.ID)
	if err != nil {
		return nil, err
	}

	return ann, nil
}

func (s *AnnService) GetListAnn(ctx context.Context, i *model.InfoSetting) ([]*model.Announcement, error) {
	anns, err := s.annRepository.GetAnnsList(ctx, i)
	if err != nil {
		return nil, err
	}
	//
	for ind, _ := range anns {
		anns[ind].Images, err = s.imageRepository.GetImagesByAnnouncementID(ctx, anns[ind].ID)
		if err != nil {
			return nil, err
		}
	}

	return anns, nil
}
