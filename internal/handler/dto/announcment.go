package dto

import "time"

type CreateAnnouncementRequest struct {
	ID               string    `json:"id"`
	Title            string    `json:"title" validate:"required"`
	Description      string    `json:"description" validate:"required"`
	Address          string    `json:"address" validate:"required"`
	Date             time.Time `json:"date" validate:"required"`
	Contacts         string    `json:"contacts" validate:"required"`
	ModerationStatus string    `json:"moderation_status" validate:"required"`
	SearchedStatus   bool      `json:"searched_status"`
	UserID           string    `json:"user_id" validate:"required"`
}

type IDResponse struct {
	ID string `json:"id"`
}

type InfoAnnRequestBody struct {
	SearchedStatus string    `json:"searched_status,omitempty"`
	Search         string    `json:"search,omitempty"`
	TimeRange      time.Time `json:"time_range,omitempty"`
}

type AnnouncementResponse struct {
	ID               string          `json:"id"`
	Title            string          `json:"title"`
	Description      string          `json:"description"`
	Address          string          `json:"address"`
	Date             time.Time       `json:"date"`
	Images           []ImageResponse `json:"images"`
	Contacts         string          `json:"contacts"`
	ModerationStatus string          `json:"moderation_status"`
	SearchedStatus   bool            `json:"searched_status"`
	UserID           string          `json:"user_id"`
}

type IDRequest struct {
	ID string `json:"id" validate:"required"`
}
