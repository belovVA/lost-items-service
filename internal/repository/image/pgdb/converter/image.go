package converter

import (
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	modelpg "lost-items-service/internal/repository/image/pgdb/model"
)

func FromImageModelToRepo(img *model.Image) *modelpg.Image {
	return &modelpg.Image{
		ID:    uuid.Nil,
		Bytes: img.Bytes,
		AnnID: img.AnnouncementID,
	}
}

func FromImageRepoToModel(i *modelpg.Image) *model.Image {
	return &model.Image{
		ID:             i.ID,
		Bytes:          i.Bytes,
		AnnouncementID: i.AnnID,
	}
}
