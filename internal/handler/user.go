package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"lost-items-service/internal/converter"
	"lost-items-service/internal/handler/dto"
	"lost-items-service/internal/handler/pkg/response"
	"lost-items-service/internal/model"
)

type UserService interface {
	GetOwnUser(ctx context.Context) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	InfoUsers(ctx context.Context, limits *model.InfoUsers) ([]*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, user *model.User) error
}

type UserHandlers struct {
	Service UserService
}

func NewUserHandler(service UserService) *UserHandlers {
	return &UserHandlers{
		Service: service,
	}
}

func (h *UserHandlers) InfoOwnUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.Service.GetOwnUser(r.Context())
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response.SuccessJSON(w, user, http.StatusOK)
}

func (h *UserHandlers) InfoUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")

	userID, err := converter.ToUUIDFromStringID(userIDStr)
	if err != nil {
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetUserByID(r.Context(), userID)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusNotFound)
		return
	}

	response.SuccessJSON(w, user, http.StatusOK)
}

func (h *UserHandlers) UsersInfo(w http.ResponseWriter, r *http.Request) {
	var reqQuery dto.InfoUsersRequestQuery
	var reqBody dto.InfoUsersRequestBody

	_ = getLogger(r)

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(false)

	if err := decoder.Decode(&reqQuery, r.URL.Query()); err != nil {
		response.WriteError(w, ErrQueryParameters, http.StatusBadRequest)
		//logger.InfoContext(r.Context(), ErrQueryParameters, slog.String(ErrorKey, err.Error()))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		return
	}

	infoUsers := converter.FromInfoUsersRequestToInfoUsersModel(&reqBody, &reqQuery)

	if err := validateRole(infoUsers.Role); err != nil {
		response.WriteError(w, ErrInvalidRole, http.StatusBadRequest)
		//logger.Info(ErrInvalidRole, slog.String(ErrorKey, err.Error()))
		return
	}

	//
	users, err := h.Service.InfoUsers(r.Context(), infoUsers)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusInternalServerError)
	}

	response.SuccessJSON(w, users, http.StatusOK)
}
