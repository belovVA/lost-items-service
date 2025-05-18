package dto

import "time"

type CreateAnnouncementRequest struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Address          string    `json:"address"`
	Date             time.Time `json:"date"`
	Contacts         string    `json:"contacts"`
	ModerationStatus string    `json:"moderation_status"`
	SearchedStatus   bool      `json:"searched_status"`
	UserID           string    `json:"user_id"`
}

type IDResponse struct {
	ID string `json:"id"`
}
