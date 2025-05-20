package converter

import (
	"lost-items-service/internal/handler/dto"
	"lost-items-service/internal/model"
)

func FromInfoAnnRequestToModel(body *dto.InfoAnnRequestBody, query *dto.InfoRequestQuery) *model.InfoSetting {
	InfoModel := model.InfoSetting{
		OrderByField: body.SearchedStatus,
		Search:       body.Search,
		Page:         query.Page,
		Limit:        query.Limit,
		TimeOrder:    body.TimeRange,
	}
	setDefaultsPagination(&InfoModel)
	return &InfoModel
}

func FromUserToUserShortResponse(user *model.User) dto.UserShortResponse {
	return dto.UserShortResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
	}
}
