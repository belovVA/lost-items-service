package modelredis

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID               uuid.UUID `redis:"id"`
	Title            string    `redis:"title"`
	Description      string    `redis:"description"`
	Address          string    `redis:"address"`
	Date             time.Time `redis:"date"`
	Contacts         string    `redis:"contacts"`
	ModerationStatus string    `redis:"moderation_status"`
	SearchedStatus   bool      `redis:"searched_status"`
	UserID           uuid.UUID `redis:"owner_id"`
}
