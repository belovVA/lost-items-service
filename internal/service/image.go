package service

import (
	"context"
	"fmt"

	"lost-items-service/internal/model"
)

type ImageService struct {
	imageRepository ImageRepository
	annRepository   AnnRepository
}

func NewImageService(repoImg ImageRepository, repoAnn AnnRepository) *ImageService {
	return &ImageService{
		imageRepository: repoImg,
		annRepository:   repoAnn,
	}
}

func (s *ImageService) CreateImages(ctx context.Context, i *model.ImagesList) (*model.ImagesList, error) {
	if len(i.Images) == 0 {
		return nil, fmt.Errorf("empty images")
	}

	if _, err := s.annRepository.GetAnnByID(ctx, i.Images[0].AnnouncementID); err != nil {
		return nil, err
	}

	var err error
	for ind, _ := range i.Images {
		if i.Images[ind].ID, err = s.imageRepository.CreateImage(ctx, i.Images[ind]); err != nil {
			return nil, err
		}
	}

	return i, nil
}

func (s *ImageService) GetImages(ctx context.Context, ann *model.Announcement) ([]*model.Image, error) {
	if _, err := s.annRepository.GetAnnByID(ctx, ann.ID); err != nil {
		return nil, err
	}

	return s.imageRepository.GetImagesByAnnouncementID(ctx, ann.ID)
}
