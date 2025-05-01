package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"lost-items-service/internal/converter"
	"lost-items-service/internal/handler/pkg/response"
	"lost-items-service/internal/model"
)

type UserService interface {
	GetOwnUser(ctx context.Context) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
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
