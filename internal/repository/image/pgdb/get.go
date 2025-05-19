package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/image/pgdb/converter"
	modelpg "lost-items-service/internal/repository/image/pgdb/model"
)

func (r *imageRepo) GetImagesByAnnID(ctx context.Context, annID uuid.UUID) ([]*model.Image, error) {
	images := make([]*model.Image, 0, 3)
	query, args, err := sq.
		Select(
			imageIDColumn,
			bytesColumn,
			annIDColumn,
		).
		From(imageTable).
		Where(sq.Eq{annIDColumn: annID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, model.ErrorBuildQuery
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, model.ErrorExecuteQuery
	}

	defer rows.Close()

	for rows.Next() {
		var img modelpg.Image
		if err = rows.Scan(
			&img.ID,
			&img.Bytes,
			&img.AnnID,
		); err != nil {
			return nil, model.ErrorScanRows
		}
		images = append(images, converter.FromImageRepoToModel(&img))
	}

	if rows.Err() != nil {
		return nil, model.ErrorScanRows
	}

	return images, nil
}
