package image

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"
	"lost-items-service/internal/db/pgxdb"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/image/pgdb"
)

type ImageRepository struct {
	Pg    pgdb.ImagePGRepository
	group *singleflight.Group
}

func NewRepository(pg pgxdb.DB) *ImageRepository {
	return &ImageRepository{
		Pg:    pgdb.NewRepository(pg),
		group: &singleflight.Group{},
	}
}

func (r *ImageRepository) CreateImage(ctx context.Context, image *model.Image) (uuid.UUID, error) {
	return r.Pg.CreateImage(ctx, image)
}

func (r *ImageRepository) GetImagesByAnnouncementID(ctx context.Context, annID uuid.UUID) ([]*model.Image, error) {
	v, err, _ := r.group.Do("images"+annID.String(), func() (interface{}, error) {
		return r.Pg.GetImagesByAnnID(ctx, annID)
	})

	if err != nil {
		return nil, err
	}

	return v.([]*model.Image), nil
}
