package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"lost-items-service/internal/converter"
	"lost-items-service/internal/handler/dto"
	"lost-items-service/internal/handler/pkg/response"
	"lost-items-service/internal/model"
)

type AnnouncementService interface {
	CreateAnnouncement(ctx context.Context, ann *model.Announcement) (uuid.UUID, error)
	GetAnn(ctx context.Context, id uuid.UUID) (*model.Announcement, error)
}

type AnnHandlers struct {
	Service AnnouncementService
}

func NewAnnHandler(service AnnouncementService) *AnnHandlers {
	return &AnnHandlers{
		Service: service,
	}
}

func (h *AnnHandlers) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAnnouncementRequest

	logger := getLogger(r)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.Info("CreateAnnouncement "+ErrBodyRequest, slog.String(ErrorKey, err.Error()))
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		response.WriteError(w, ErrRequestFields, http.StatusBadRequest)
		logger.Info("CreateAnnouncement "+ErrRequestFields, slog.String(ErrorKey, err.Error()))
		return
	}

	a, err := converter.ToAnnouncementModelFromRequest(&req)
	if err != nil {
		logger.InfoContext(r.Context(), "CreateAnnouncement "+ErrUUIDParsing, slog.String(ErrorKey, err.Error()))
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		return
	}

	if err = checkStatus(a.ModerationStatus); err != nil {
		logger.InfoContext(r.Context(), "CreateAnnouncement "+ErrInvalidStatus, slog.String(ErrorKey, err.Error()))
		response.WriteError(w, ErrInvalidStatus, http.StatusBadRequest)
		return
	}
	annID, err := h.Service.CreateAnnouncement(r.Context(), a)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusInternalServerError)
		logger.Info("error create Ann", slog.String(ErrorKey, err.Error()), slog.String("userID", a.UserID.String()))
		return
	}

	logger.InfoContext(r.Context(), "success create announcement", slog.String("ID", annID.String()))
	response.SuccessJSON(w, converter.ToIDResponse(annID), http.StatusCreated)
}

func checkStatus(status string) error {
	switch status {
	case "in_progress", "canceled", "accepted":
		return nil
	}
	return fmt.Errorf("invalid moderation status")
}
