package converter

import (
	"lost-items-service/internal/handler/dto"
	"lost-items-service/internal/model"
)

func FromInfoUsersRequestToInfoUsersModel(body *dto.InfoUsersRequestBody, query *dto.InfoUsersRequestQuery) *model.InfoUsers {
	InfoModel := model.InfoUsers{
		Role:  body.Role,
		Page:  query.Page,
		Limit: query.Limit,
	}
	setDefaultsPagination(&InfoModel)
	return &InfoModel
}

func setDefaultsPagination(q *model.InfoUsers) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 30 {
		q.Limit = 10
	}
}
