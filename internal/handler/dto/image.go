package dto

import "github.com/google/uuid"

type CreateImagesRequest struct {
	AnnouncementID uuid.UUID `json:"announcement_id"`
	Images         [][]byte  `json:"images"`
}
