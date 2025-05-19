package dto

import "github.com/google/uuid"

type CreateImagesRequest struct {
	AnnouncementID uuid.UUID `json:"announcement_id" validate:"required"`
	Images         [][]byte  `json:"images,omitempty"`
}

type ListImageResponse struct {
	AnnouncementID uuid.UUID       `json:"announcement_id"`
	Images         []ImageResponse `json:"images"`
}

type ImageResponse struct {
	ID    string `json:"id"`
	Bytes []byte `json:"bytes"`
	AnnId string `json:"announcement_id"`
}
