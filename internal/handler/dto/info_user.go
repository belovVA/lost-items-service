package dto

import "github.com/google/uuid"

type InfoUsersRequestBody struct {
	Role string `json:"role,omitempty"`
}

type InfoUsersResponse struct {
	Data []UserShortResponse
}

type InfoUsersRequestQuery struct {
	Page  int `schema:"page"      validate:"omitempty"`
	Limit int `schema:"limit"     validate:"omitempty"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
}

type UpdateRequest struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Surname  string    `json:"surname,omitempty"`
	Phone    string    `json:"phone,omitempty"`
	Email    string    `json:"email,omitempty"`
	Password string    `json:"password,omitempty"`
	Role     string    `json:"role,omitempty"`
}
