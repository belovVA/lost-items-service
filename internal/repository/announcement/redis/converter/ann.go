package converter

import (
	"lost-items-service/internal/model"
	modelredis "lost-items-service/internal/repository/announcement/redis/model"
)

func FromAnnModelToRedis(a *model.Announcement) *modelredis.Announcement {
	return &modelredis.Announcement{
		ID:               a.ID,
		Title:            a.Title,
		Description:      a.Description,
		Address:          a.Address,
		Date:             a.Date,
		Contacts:         a.Contacts,
		ModerationStatus: a.ModerationStatus,
		SearchedStatus:   a.SearchedStatus,
		OwnerID:          a.OwnerID,
	}
}

func FromRedisToAnnModel(a *modelredis.Announcement) *model.Announcement {
	return &model.Announcement{
		ID:               a.ID,
		Title:            a.Title,
		Description:      a.Description,
		Address:          a.Address,
		Date:             a.Date,
		Contacts:         a.Contacts,
		ModerationStatus: a.ModerationStatus,
		SearchedStatus:   a.SearchedStatus,
		OwnerID:          a.OwnerID,
	}
}
