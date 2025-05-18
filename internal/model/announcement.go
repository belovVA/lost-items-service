package model

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID               uuid.UUID
	Title            string
	Description      string
	Address          string
	Date             time.Time
	Contacts         string
	ModerationStatus string
	SearchedStatus   bool
	Images           []Image
	UserID           uuid.UUID
}
