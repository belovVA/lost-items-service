package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"lost-items-service/internal/converter"
	"lost-items-service/internal/handler/dto"
	"lost-items-service/internal/handler/pkg/response"
	"lost-items-service/internal/model"
)

type AnnouncementService interface {
	CreateAnnouncement(ctx context.Context, ann *model.Announcement) (uuid.UUID, error)
	GetAnn(ctx context.Context, id uuid.UUID) (*model.Announcement, error)
	GetListAnn(ctx context.Context, i *model.InfoSetting) ([]*model.Announcement, error)
	GetListAnnByUser(ctx context.Context, i *model.InfoSetting) ([]*model.Announcement, error)
	UpdateAnn(ctx context.Context, a *model.Announcement) error
	UpdateMoserStatusAnn(ctx context.Context, a *model.Announcement) error
	DeleteAnn(ctx context.Context, id uuid.UUID) error
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
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		logger.Info("error create Ann", slog.String(ErrorKey, err.Error()), slog.String("userID", a.UserID.String()))
		return
	}

	logger.InfoContext(r.Context(), "success create announcement", slog.String("ID", annID.String()))
	response.SuccessJSON(w, converter.ToIDResponse(annID), http.StatusCreated)
}

func checkStatus(status string) error {
	switch status {
	case "watching", "canceled", "accepted":
		return nil
	}
	return fmt.Errorf("invalid moderation status")
}

func (h *AnnHandlers) AnnsInfo(w http.ResponseWriter, r *http.Request) {
	var reqQuery dto.InfoRequestQuery
	var reqBody dto.InfoAnnRequestBody

	logger := getLogger(r)

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(false)

	if err := decoder.Decode(&reqQuery, r.URL.Query()); err != nil {
		response.WriteError(w, ErrQueryParameters, http.StatusBadRequest)
		logger.InfoContext(r.Context(), "AnnsInfo "+ErrQueryParameters, slog.String(ErrorKey, err.Error()))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Println(r.Body)
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.InfoContext(r.Context(), "AnnsInfo "+ErrBodyRequest, slog.String(ErrorKey, err.Error()))

		return
	}

	infoAnns := converter.FromInfoAnnRequestToModel(&reqBody, &reqQuery)

	if infoAnns.OrderByField != "" {
		if infoAnns.OrderByField != "true" && infoAnns.OrderByField != "false" {
			response.WriteError(w, "invalid status", http.StatusBadRequest)
			logger.InfoContext(r.Context(), "invalid searched status", slog.String(ErrorKey, infoAnns.OrderByField))
			return
		}
	}

	//
	anns, err := h.Service.GetListAnn(r.Context(), infoAnns)
	if err != nil {
		logger.InfoContext(r.Context(), "AnnsInfo service", slog.String(ErrorKey, err.Error()))
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info("success info anns", slog.String("role order", infoAnns.OrderByField), slog.String("search word", infoAnns.Search))
	annsResp := make([]dto.AnnouncementResponse, 0, len(anns))
	for _, a := range anns {
		annsResp = append(annsResp, converter.ToAnnouncementResponseFromModel(a))
	}

	response.SuccessJSON(w, annsResp, http.StatusOK)
}

func (h *AnnHandlers) GetAnnouncement(w http.ResponseWriter, r *http.Request) {
	var req dto.IDRequest

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

	annID, err := converter.ToUUIDFromStringID(req.ID)
	if err != nil {
		logger.InfoContext(r.Context(), "GetAnnouncement "+ErrUUIDParsing, slog.String(ErrorKey, err.Error()))
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		return
	}
	ann, err := h.Service.GetAnn(r.Context(), annID)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusNotFound)
		logger.Info("error get Ann", slog.String(ErrorKey, err.Error()), slog.String("AnnID", annID.String()))
		return
	}
	log.Println(ann)

	response.SuccessJSON(w, converter.ToAnnouncementResponseFromModel(ann), http.StatusCreated)

}

func (h *AnnHandlers) GetUserAnnouncements(w http.ResponseWriter, r *http.Request) {
	var reqQuery dto.InfoRequestQuery
	var reqBody dto.InfoAnnRequestBody

	logger := getLogger(r)

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(false)

	if err := decoder.Decode(&reqQuery, r.URL.Query()); err != nil {
		response.WriteError(w, ErrQueryParameters, http.StatusBadRequest)
		logger.InfoContext(r.Context(), "AnnsInfo "+ErrQueryParameters, slog.String(ErrorKey, err.Error()))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Println(r.Body)
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.InfoContext(r.Context(), "AnnsInfo "+ErrBodyRequest, slog.String(ErrorKey, err.Error()))

		return
	}

	infoAnns := converter.FromInfoAnnRequestToModel(&reqBody, &reqQuery)

	if infoAnns.OrderByField != "" {
		if infoAnns.OrderByField != "true" && infoAnns.OrderByField != "false" {
			response.WriteError(w, "invalid status", http.StatusBadRequest)
			logger.InfoContext(r.Context(), "invalid searched status", slog.String(ErrorKey, infoAnns.OrderByField))
			return
		}
	}

	//
	anns, err := h.Service.GetListAnnByUser(r.Context(), infoAnns)
	if err != nil {
		logger.InfoContext(r.Context(), "AnnsInfo service", slog.String(ErrorKey, err.Error()))
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info("success info user anns", slog.String("role order", infoAnns.OrderByField), slog.String("search word", infoAnns.Search))
	annsResp := make([]dto.AnnouncementResponse, 0, len(anns))
	for _, a := range anns {
		annsResp = append(annsResp, converter.ToAnnouncementResponseFromModel(a))
	}

	response.SuccessJSON(w, annsResp, http.StatusOK)
}

func (h *AnnHandlers) UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateAnnouncementRequest

	logger := getLogger(r)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.Info("UpdateAnnouncement "+ErrBodyRequest, slog.String(ErrorKey, err.Error()))
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		response.WriteError(w, ErrRequestFields, http.StatusBadRequest)
		logger.Info("UpdateAnnouncement "+ErrRequestFields, slog.String(ErrorKey, err.Error()))
		return
	}

	a, err := converter.ToAnnouncementModelFromUpdateRequest(&req)
	if err != nil {
		logger.InfoContext(r.Context(), "UpdateAnnouncement "+ErrUUIDParsing, slog.String(ErrorKey, err.Error()))
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		return
	}

	if err = checkStatus(a.ModerationStatus); err != nil {
		logger.InfoContext(r.Context(), "UpdateAnnouncement "+ErrInvalidStatus, slog.String(ErrorKey, err.Error()))
		response.WriteError(w, ErrInvalidStatus, http.StatusBadRequest)
		return
	}

	if err = h.Service.UpdateAnn(r.Context(), a); err != nil {
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		logger.Info("error update Ann", slog.String(ErrorKey, err.Error()), slog.String("userID", a.UserID.String()))
		return
	}

	logger.InfoContext(r.Context(), "success update announcement", slog.String("ID", a.ID.String()))
	response.Success(w, http.StatusCreated)
}

func (h *AnnHandlers) DeleteAnn(w http.ResponseWriter, r *http.Request) {
	var req dto.IDRequest
	logger := getLogger(r)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.Info("DeleteAnn "+ErrBodyRequest, slog.String(ErrorKey, err.Error()))
		return
	}

	annID, err := converter.ToUUIDFromStringID(req.ID)
	if err != nil {
		logger.InfoContext(r.Context(), "DeleteAnn"+ErrUUIDParsing, slog.String(ErrorKey, err.Error()))

		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteAnn(r.Context(), annID)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusInternalServerError)
		logger.Info("DeleteAnn", slog.String(ErrorKey, err.Error()))
		return
	}

	logger.Info("DeleteAnn success", slog.String(AnnIDKey, annID.String()))
	response.Success(w, http.StatusOK)
}

func (h *AnnHandlers) UpdateModerationStatusAnnouncement(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateModerationStatusRequest

	logger := getLogger(r)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.Info("UpdateModerationAnnouncement "+ErrBodyRequest, slog.String(ErrorKey, err.Error()))
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		response.WriteError(w, ErrRequestFields, http.StatusBadRequest)
		logger.Info("UpdateModerationAnnouncement "+ErrRequestFields, slog.String(ErrorKey, err.Error()))
		return
	}

	a, err := converter.ToAnnouncementModelFromUpdateMoserRequest(&req)
	if err != nil {
		logger.InfoContext(r.Context(), "UpdateModerationAnnouncement "+ErrUUIDParsing, slog.String(ErrorKey, err.Error()))
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		return
	}

	if err = checkStatus(a.ModerationStatus); err != nil {
		logger.InfoContext(r.Context(), "UpdateModerationAnnouncement "+ErrInvalidStatus, slog.String(ErrorKey, err.Error()))
		response.WriteError(w, ErrInvalidStatus, http.StatusBadRequest)
		return
	}

	if err = h.Service.UpdateMoserStatusAnn(r.Context(), a); err != nil {
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		logger.Info("error update MStatus Ann", slog.String(ErrorKey, err.Error()), slog.String("userID", a.UserID.String()))
		return
	}

	logger.InfoContext(r.Context(), "success update announcement", slog.String("ID", a.ID.String()))
	response.Success(w, http.StatusCreated)
}
