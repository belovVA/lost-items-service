package pgdb

import (
	"context"

	"github.com/google/uuid"

	"lost-items-service/internal/db/pgxdb"
	"lost-items-service/internal/model"
)

const (
	imageTable    = "images"
	imageIDColumn = "id"
	bytesColumn   = "bytes"
	annIDColumn   = "announcement_id"
)

type ImagePGRepository interface {
	CreateImage(ctx context.Context, image *model.Image) (uuid.UUID, error)
	GetImagesByAnnID(ctx context.Context, annID uuid.UUID) ([]*model.Image, error)
}

type imageRepo struct {
	DB pgxdb.DB
}

func NewRepository(db pgxdb.DB) ImagePGRepository {
	return &imageRepo{
		DB: db,
	}
}
