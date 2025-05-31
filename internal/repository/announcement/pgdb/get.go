package pgdb

import (
	"context"
	"fmt"
	"strconv"

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
			userIDColumn,
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
		&ann.UserID,
	)

	if err != nil {
		return nil, model.ErrorNotFound
	}

	return converter.FromRepoToAnnModel(&ann), nil
}

func (r *annRepo) GetListAnnouncement(ctx context.Context, info *model.InfoSetting) ([]*model.Announcement, error) {
	limits := converter.FromInfoModelToRepo(info)
	anns := make([]*model.Announcement, 0, info.Limit)

	req := sq.
		Select(
			annIDColumn,
			annTitleColumn,
			annDescColumn,
			annAddressColumn,
			annDateColumn,
			annContactsColumn,
			annSearchedStatusColumn,
			annModerationStatusColumn,
			userIDColumn,
		).
		From(annTable).
		PlaceholderFormat(sq.Dollar)
	if limits.ModerStatus != "" {
		req = req.Where(sq.Eq{annModerationStatusColumn: limits.ModerStatus})
	}
	if limits.FieldOrder != "" {
		boolVal, err := strconv.ParseBool(limits.FieldOrder)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean value for searched_status: %v", limits.FieldOrder)
		}
		req = req.Where(sq.Eq{annSearchedStatusColumn: boolVal})
	}
	if limits.Search != "" {
		searchQuery := fmt.Sprintf("plainto_tsquery('russian', unaccent('%s'))", limits.Search)
		req = req.Where("search_vector @@ " + searchQuery)
	}
	if limits.TimeRange != nil {
		req = req.Where(sq.GtOrEq{annDateColumn: limits.TimeRange})
	}

	req = req.Limit(limits.Limit).Offset(limits.Offset)

	query, args, err := req.ToSql()
	if err != nil {
		return nil, model.ErrorBuildQuery
	}
	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, model.ErrorExecuteQuery
	}
	defer rows.Close()
	for rows.Next() {
		var a modelpg.Announcement
		if err = rows.Scan(
			&a.ID,
			&a.Title,
			&a.Description,
			&a.Address,
			&a.Date,
			&a.Contacts,
			&a.SearchedStatus,
			&a.ModerationStatus,
			&a.UserID,
		); err != nil {
			return nil, model.ErrorScanRows
		}

		anns = append(anns, converter.FromRepoToAnnModel(&a))
	}

	if rows.Err() != nil {
		return nil, model.ErrorScanRows
	}
	return anns, nil
}

func (r *annRepo) GetAnnsByUserID(ctx context.Context, userID uuid.UUID, info *model.InfoSetting) ([]*model.Announcement, error) {
	limits := converter.FromInfoModelToRepo(info)
	anns := make([]*model.Announcement, 0, info.Limit)

	req := sq.
		Select(
			annIDColumn,
			annTitleColumn,
			annDescColumn,
			annAddressColumn,
			annDateColumn,
			annContactsColumn,
			annSearchedStatusColumn,
			annModerationStatusColumn,
			userIDColumn,
		).
		From(annTable).
		Where(sq.Eq{userIDColumn: userID}).
		PlaceholderFormat(sq.Dollar)
	if limits.FieldOrder != "" {
		boolVal, err := strconv.ParseBool(limits.FieldOrder)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean value for searched_status: %v", limits.FieldOrder)
		}
		req = req.Where(sq.Eq{annSearchedStatusColumn: boolVal})
	}
	if limits.Search != "" {
		req = req.Where(`
		search_vector @@ plainto_tsquery('russian', unaccent($1))
	`, limits.Search)
	}

	req = req.Limit(limits.Limit).Offset(limits.Offset)

	query, args, err := req.ToSql()
	if err != nil {
		return nil, model.ErrorBuildQuery
	}
	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, model.ErrorExecuteQuery
	}
	defer rows.Close()
	for rows.Next() {
		var a modelpg.Announcement
		if err = rows.Scan(
			&a.ID,
			&a.Title,
			&a.Description,
			&a.Address,
			&a.Date,
			&a.Contacts,
			&a.SearchedStatus,
			&a.ModerationStatus,
			&a.UserID,
		); err != nil {
			return nil, model.ErrorScanRows
		}

		anns = append(anns, converter.FromRepoToAnnModel(&a))
	}

	if rows.Err() != nil {
		return nil, model.ErrorScanRows
	}
	return anns, err
}
