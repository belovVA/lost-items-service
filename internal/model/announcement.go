package model

import "github.com/google/uuid"

type Announcement struct {
	ID          uuid.UUID
	Title       string
	Description string
	Address     string
	Contacts    string
	Images      []Image
	OwnerID     uuid.UUID
}
