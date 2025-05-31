package converter

import (
	"time"

	"github.com/google/uuid"
	"lost-items-service/internal/model"
	modelredis "lost-items-service/internal/repository/announcement/redis/model"
)

func FromAnnModelToRedis(a *model.Announcement) *modelredis.Announcement {
	return &modelredis.Announcement{
		ID:               a.ID.String(),
		Title:            a.Title,
		Description:      a.Description,
		Address:          a.Address,
		Date:             a.Date.Format(time.RFC3339),
		Contacts:         a.Contacts,
		ModerationStatus: a.ModerationStatus,
		SearchedStatus:   a.SearchedStatus,
		UserID:           a.UserID.String(),
	}
}

func FromRedisToAnnModel(a *modelredis.Announcement) (*model.Announcement, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(a.UserID)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse(time.RFC3339, a.Date)
	if err != nil {
		return nil, err
	}
	return &model.Announcement{
		ID:               id,
		Title:            a.Title,
		Description:      a.Description,
		Address:          a.Address,
		Date:             date,
		Contacts:         a.Contacts,
		ModerationStatus: a.ModerationStatus,
		SearchedStatus:   a.SearchedStatus,
		UserID:           userID,
	}, nil
}
