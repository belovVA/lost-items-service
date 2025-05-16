package handler

import (
	"context"
	"encoding/json"
	"log/slog"
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
	logger := getLogger(r)
	if err != nil {

		response.WriteError(w, err.Error(), http.StatusUnauthorized)
		logger.InfoContext(r.Context(), "error get info user", slog.String(ErrorKey, err.Error()))
		return
	}

	logger.InfoContext(r.Context(), "success get info user")
	response.SuccessJSON(w, converter.ToUserResponseFromUserModel(user), http.StatusOK)
}

func (h *UserHandlers) InfoUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	logger := getLogger(r)

	userID, err := converter.ToUUIDFromStringID(userIDStr)
	if err != nil {
		logger.InfoContext(r.Context(), "InfoUser"+ErrUUIDParsing, slog.String(ErrorKey, err.Error()))
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetUserByID(r.Context(), userID)
	if err != nil {
		logger.InfoContext(r.Context(), "error get info other user", slog.String("otherID", userIDStr))
		response.WriteError(w, err.Error(), http.StatusNotFound)
		return
	}

	logger.InfoContext(r.Context(), "success get info other user", slog.String("otherID", userIDStr))
	response.SuccessJSON(w, converter.ToUserResponseFromUserModel(user), http.StatusOK)
}

func (h *UserHandlers) UsersInfo(w http.ResponseWriter, r *http.Request) {
	var reqQuery dto.InfoUsersRequestQuery
	var reqBody dto.InfoUsersRequestBody

	logger := getLogger(r)

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(false)

	if err := decoder.Decode(&reqQuery, r.URL.Query()); err != nil {
		response.WriteError(w, ErrQueryParameters, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrQueryParameters, slog.String(ErrorKey, err.Error()))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrBodyRequest, slog.String(ErrorKey, err.Error()))

		return
	}

	infoUsers := converter.FromInfoUsersRequestToInfoUsersModel(&reqBody, &reqQuery)

	if err := validateRole(infoUsers.Role); err != nil {
		response.WriteError(w, ErrInvalidRole, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrInvalidRole, slog.String(ErrorKey, err.Error()))
		return
	}

	//
	users, err := h.Service.InfoUsers(r.Context(), infoUsers)
	if err != nil {
		logger.InfoContext(r.Context(), ErrBodyRequest, slog.String(ErrorKey, err.Error()))
		response.WriteError(w, err.Error(), http.StatusInternalServerError)
	}

	logger.Info("success info users", slog.String("role order", infoUsers.Role))
	response.SuccessJSON(w, users, http.StatusOK)
}

func (h *UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	logger := getLogger(r)

	userID, err := converter.ToUUIDFromStringID(userIDStr)
	if err != nil {
		logger.InfoContext(r.Context(), "UpdateUser"+ErrUUIDParsing, slog.String(ErrorKey, err.Error()))

		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		return
	}

	var req dto.UpdateRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.Info("UpdateUser"+ErrBodyRequest, slog.String(ErrorKey, err.Error()))
		return
	}

	v := getValidator(r)
	if err = v.Struct(req); err != nil {
		response.WriteError(w, ErrRequestFields, http.StatusBadRequest)
		logger.Info("UpdateUser"+ErrRequestFields, slog.String(ErrorKey, err.Error()))
		return
	}

	userModel := converter.ToUserFromUpdateUserRequest(&req)
	if userModel.Role != "" {
		if err = validateRole(userModel.Role); err != nil {
			response.WriteError(w, ErrInvalidRole, http.StatusBadRequest)
			logger.Info("UpdateUser "+ErrInvalidRole, slog.String(ErrorKey, err.Error()))
			return
		}
	}
	userModel.ID = userID

	err = h.Service.UpdateUser(r.Context(), userModel)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusConflict)
		logger.Info("UpdateUser", slog.String(ErrorKey, err.Error()))
		return
	}

	logger.Info("UpdateUser success", slog.String(UserIDKey, userID.String()))
	response.Success(w, http.StatusNoContent)
}

func (h *UserHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	logger := getLogger(r)

	userID, err := converter.ToUUIDFromStringID(userIDStr)
	if err != nil {
		logger.InfoContext(r.Context(), "UpdateUser"+ErrUUIDParsing, slog.String(ErrorKey, err.Error()))

		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		return
	}

	var user model.User
	user.ID = userID

	err = h.Service.DeleteUser(r.Context(), &user)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusInternalServerError)
		logger.Info("DeleteUser", slog.String(ErrorKey, err.Error()))
		return
	}

	logger.Info("DeleteUser success", slog.String(UserIDKey, userID.String()))
	response.Success(w, http.StatusNoContent)
}
