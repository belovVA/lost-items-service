package converter

import (
	"fmt"

	"lost-items-service/internal/handler/dto"
	"lost-items-service/internal/model"
)

func ToImageModelFromCreateRequest(i *dto.CreateImagesRequest) (*model.ImagesList, error) {
	if len(i.Images) == 0 {
		return nil, fmt.Errorf("images list is empty")
	}
	images := make([]*model.Image, 0, len(i.Images))
	for _, v := range i.Images {
		img := &model.Image{
			Bytes:          v,
			AnnouncementID: i.AnnouncementID,
		}
		images = append(images, img)
	}

	return &model.ImagesList{
		AnnID:  i.AnnouncementID,
		Images: images,
	}, nil
}

func ToResponseFromImagesModel(i *model.ImagesList) *dto.ListImageResponse {
	if len(i.Images) == 0 {
		return &dto.ListImageResponse{
			AnnouncementID: i.AnnID,
			Images:         make([]dto.ImageResponse, 0),
		}
	}

	resp := &dto.ListImageResponse{
		AnnouncementID: i.AnnID,
		Images:         make([]dto.ImageResponse, 0, len(i.Images)),
	}

	for _, img := range i.Images {
		resp.Images = append(resp.Images, dto.ImageResponse{
			ID:    img.ID.String(),
			Bytes: img.Bytes,
			AnnId: img.AnnouncementID.String(),
		})
	}

	return resp
}

func ToImageResponseFromModel(a *model.Image) dto.ImageResponse {
	return dto.ImageResponse{
		ID:    a.ID.String(),
		Bytes: a.Bytes,
		AnnId: a.AnnouncementID.String(),
	}
}
