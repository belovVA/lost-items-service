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
		OwnerID:          a.OwnerID,
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
		OwnerID:          a.OwnerID,
	}
}
