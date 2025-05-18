package model

import "github.com/google/uuid"

type Image struct {
	ID             uuid.UUID
	ImageBytes     []byte
	AnnouncementID uuid.UUID
}
