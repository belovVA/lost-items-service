package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/announcement/pgdb/converter"
)

func (r *annRepo) UpdateFields(ctx context.Context, ann *model.Announcement) error {
	a := converter.FromAnnModelToRepo(ann)
	q := sq.
		Update(annTable).
		Set(annTitleColumn, a.Title).
		Set(annDescColumn, a.Description).
		Set(annAddressColumn, a.Address).
		Set(annDateColumn, a.Date).
		Set(annContactsColumn, a.Contacts).
		Set(annSearchedStatusColumn, a.SearchedStatus)
	if a.ModerationStatus != "" {
		q = q.Set(annModerationStatusColumn, a.ModerationStatus)
	}
	query, args, err := q.Where(sq.Eq{annIDColumn: a.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return model.ErrorBuildQuery
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return model.ErrorExecuteQuery
	}

	return nil
}

func (r *annRepo) UpdateModerationStatus(ctx context.Context, ann *model.Announcement) error {
	a := converter.FromAnnModelToRepo(ann)

	builder := sq.Update(annTable).
		Set(annModerationStatusColumn, a.ModerationStatus).
		Where(sq.Eq{annIDColumn: a.ID}).
		PlaceholderFormat(sq.Dollar)

	if a.UserID != uuid.Nil {
		builder = builder.Set(userIDColumn, a.UserID)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return model.ErrorBuildQuery
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return model.ErrorExecuteQuery
	}

	return nil
}
