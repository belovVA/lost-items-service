package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/announcement/pgdb/converter"
	modelpg "lost-items-service/internal/repository/announcement/pgdb/model"
)

func (r *annRepo) GetAnnByID(ctx context.Context, id uuid.UUID) (*model.Announcement, error) {
	var ann modelpg.Announcement

	query, args, err := sq.
		Select(
			annIDColumn,
			annTitleColumn,
			annDescColumn,
			annAddressColumn,
			annDateColumn,
			annContactsColumn,
			annSearchedStatusColumn,
			annModerationStatusColumn,
			ownerIDColumn,
		).
		From(annTable).
		Where(sq.Eq{annIDColumn: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&ann.ID,
		&ann.Title,
		&ann.Description,
		&ann.Address,
		&ann.Date,
		&ann.Contacts,
		&ann.SearchedStatus,
		&ann.ModerationStatus,
		&ann.OwnerID,
	)

	if err != nil {
		return nil, model.ErrorNotFound
	}

	return converter.FromRepoToAnnModel(&ann), nil
}
