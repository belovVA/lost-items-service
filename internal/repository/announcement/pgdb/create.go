package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/announcement/pgdb/converter"
)

func (r *annRepo) AddAnn(ctx context.Context, ann *model.Announcement) (uuid.UUID, error) {
	var id uuid.UUID
	a := converter.FromAnnModelToRepo(ann)

	query, args, err := sq.
		Insert(annTable).
		Columns(
			annTitleColumn,
			annDescColumn,
			annAddressColumn,
			annDateColumn,
			annContactsColumn,
			annSearchedStatusColumn,
			annModerationStatusColumn,
			ownerIDColumn,
		).
		Values(
			a.Title,
			a.Description,
			a.Address,
			a.Date,
			a.Contacts,
			a.SearchedStatus,
			a.ModerationStatus,
			a.OwnerID).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING " + annIDColumn).
		ToSql()

	if err != nil {
		return uuid.Nil, model.ErrorBuildQuery
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, model.ErrorBuildQuery
	}

	return id, nil
}
