package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"lost-items-service/internal/converter"
	"lost-items-service/internal/handler/dto"
	"lost-items-service/internal/handler/pkg/response"
	"lost-items-service/internal/model"
)

type ImageService interface {
	CreateImages(ctx context.Context, i *model.ImagesList) (*model.ImagesList, error)
	GetImages(ctx context.Context, ann *model.Announcement) ([]*model.Image, error)
}

type ImageHandlers struct {
	Service ImageService
}

func NewImageHandler(service ImageService) *ImageHandlers {
	return &ImageHandlers{
		Service: service,
	}
}

func (h *ImageHandlers) AddImagesToAnnouncement(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateImagesRequest
	logger := getLogger(r)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.Info(ErrBodyRequest, slog.String(ErrorKey, err.Error()))
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		response.WriteError(w, ErrRequestFields, http.StatusBadRequest)
		logger.Info(ErrRequestFields, slog.String(ErrorKey, err.Error()))
		return
	}

	imgs, err := converter.ToImageModelFromCreateRequest(&req)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		logger.Info("error to convert ImageList model", slog.String(ErrorKey, err.Error()))
		return
	}

	res, err := h.Service.CreateImages(r.Context(), imgs)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		logger.Info("error to add images", slog.String(ErrorKey, err.Error()))
		return
	}

	resp := converter.ToResponseFromImagesModel(res)
	logger.Info("successful add images", slog.String(AnnIDKey, res.AnnID.String()))

	response.SuccessJSON(w, resp, http.StatusCreated)
}
