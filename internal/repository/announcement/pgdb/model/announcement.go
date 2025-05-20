package modelpg

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID               uuid.UUID `db:"id"`
	Title            string    `db:"title"`
	Description      string    `db:"description"`
	Address          string    `db:"address"`
	Date             time.Time `db:"date"`
	Contacts         string    `db:"contacts"`
	ModerationStatus string    `db:"moderation_status"`
	SearchedStatus   bool      `db:"searched_status"`
	UserID           uuid.UUID `db:"owner_id"`
}

type LimitsAnn struct {
	FieldOrder string
	Search     string
	Limit      uint64
	Offset     uint64
	TimeRange  time.Time
}
