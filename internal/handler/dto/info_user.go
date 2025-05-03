package dto

type InfoUsersRequestBody struct {
	Role string `json:"role,omitempty"`
}

type InfoUsersResponse struct {
	Data []UserResponse
}

type InfoUsersRequestQuery struct {
	Page  int `schema:"page"      validate:"omitempty"`
	Limit int `schema:"limit"     validate:"omitempty"`
}
