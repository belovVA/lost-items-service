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

func ToAnnouncementResponseFromModel(a *model.Announcement) dto.AnnouncementResponse {
	imgs := make([]dto.ImageResponse, 0, len(a.Images))
	for _, img := range a.Images {
		imgs = append(imgs, ToImageResponseFromModel(img))
	}
	return dto.AnnouncementResponse{
		ID:               a.ID.String(),
		Title:            a.Title,
		Description:      a.Description,
		Address:          a.Address,
		Date:             a.Date,
		Images:           imgs,
		Contacts:         a.Contacts,
		ModerationStatus: a.ModerationStatus,
		SearchedStatus:   a.SearchedStatus,
		UserID:           a.UserID.String(),
	}
}

func ToAnnouncementModelFromUpdateRequest(a *dto.UpdateAnnouncementRequest) (*model.Announcement, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return nil, err
	}

	return &model.Announcement{
		ID:               id,
		Title:            a.Title,
		Description:      a.Description,
		Address:          a.Address,
		Date:             a.Date,
		Contacts:         a.Contacts,
		ModerationStatus: a.ModerationStatus,
		SearchedStatus:   a.SearchedStatus,
	}, nil
}

func ToAnnouncementModelFromUpdateMoserRequest(a *dto.UpdateModerationStatusRequest) (*model.Announcement, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return nil, err
	}

	return &model.Announcement{
		ID:               id,
		ModerationStatus: a.ModerationStatus,
	}, nil
}
