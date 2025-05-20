package converter

import (
	"lost-items-service/internal/handler/dto"
	"lost-items-service/internal/model"
)

func FromInfoUsersRequestToInfoUsersModel(body *dto.InfoUsersRequestBody, query *dto.InfoRequestQuery) *model.InfoSetting {
	InfoModel := model.InfoSetting{
		OrderByField: body.Role,
		Search:       body.Search,
		Page:         query.Page,
		Limit:        query.Limit,
	}
	setDefaultsPagination(&InfoModel)
	return &InfoModel
}

func setDefaultsPagination(q *model.InfoSetting) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 30 {
		q.Limit = 10
	}
}
