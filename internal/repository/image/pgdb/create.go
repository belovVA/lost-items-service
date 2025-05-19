package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/image/pgdb/converter"
)

func (r *imageRepo) CreateImage(ctx context.Context, image *model.Image) (uuid.UUID, error) {
	var id uuid.UUID
	img := converter.FromImageModelToRepo(image)

	query, args, err := sq.
		Insert(imageTable).
		Columns(bytesColumn, annIDColumn).
		Values(img.Bytes, img.AnnID).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING " + imageIDColumn).
		ToSql()
	if err != nil {
		return uuid.Nil, model.ErrorBuildQuery
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, model.ErrorScanRows
	}

	return id, nil
}
