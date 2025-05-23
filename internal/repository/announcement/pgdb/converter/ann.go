package converter

import (
	"lost-items-service/internal/model"
	modelpg "lost-items-service/internal/repository/announcement/pgdb/model"
)

func FromAnnModelToRepo(a *model.Announcement) *modelpg.Announcement {
	return &modelpg.Announcement{
		ID:               a.ID,
		Title:            a.Title,
		Description:      a.Description,
		Address:          a.Address,
		Date:             a.Date,
		Contacts:         a.Contacts,
		ModerationStatus: a.ModerationStatus,
		SearchedStatus:   a.SearchedStatus,
		UserID:           a.UserID,
	}
}

func FromRepoToAnnModel(a *modelpg.Announcement) *model.Announcement {
	return &model.Announcement{
		ID:               a.ID,
		Title:            a.Title,
		Description:      a.Description,
		Address:          a.Address,
		Date:             a.Date,
		Contacts:         a.Contacts,
		ModerationStatus: a.ModerationStatus,
		SearchedStatus:   a.SearchedStatus,
		UserID:           a.UserID,
	}
}

func FromInfoModelToRepo(info *model.InfoSetting) *modelpg.LimitsAnn {
	return &modelpg.LimitsAnn{
		FieldOrder:  info.OrderByField,
		Search:      info.Search,
		Limit:       uint64(info.Limit),
		Offset:      uint64((info.Page - 1) * info.Limit),
		TimeRange:   &info.TimeOrder,
		ModerStatus: info.ModerStatus,
	}
}
