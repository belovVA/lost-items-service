package converter

import (
	"github.com/google/uuid"
	"lost-items-service/internal/handler/dto"
	"lost-items-service/internal/model"
)

func ToAnnouncementModelFromRequest(a *dto.CreateAnnouncementRequest) (*model.Announcement, error) {
	userID, err := uuid.Parse(a.UserID)
	if err != nil {
		return nil, err
	}

	return &model.Announcement{
		Title:            a.Title,
		Description:      a.Description,
		Address:          a.Address,
		Date:             a.Date,
		Contacts:         a.Contacts,
		ModerationStatus: a.ModerationStatus,
		SearchedStatus:   a.SearchedStatus,
		UserID:           userID,
	}, nil
}

func ToIDResponse(id uuid.UUID) dto.IDResponse {
	return dto.IDResponse{
		ID: id.String()}
}
