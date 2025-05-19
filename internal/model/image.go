package model

import "github.com/google/uuid"

type Image struct {
	ID             uuid.UUID
	Bytes          []byte
	AnnouncementID uuid.UUID
}

type ImagesList struct {
	AnnID  uuid.UUID
	Images []*Image
}
