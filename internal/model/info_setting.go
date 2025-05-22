package model

import "time"

type InfoSetting struct {
	OrderByField string // Role for Users, StatusModeration for ModAnn, StatusSearched for Anns
	Search       string
	Page         int
	Limit        int
	TimeOrder    time.Time
	ModerStatus  string
}
