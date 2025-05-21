package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
)

func (r *annRepo) DeleteAnn(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.
		Delete(annTable).
		Where(sq.Eq{annIDColumn: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return model.ErrorBuildQuery
	}

	cmdTag, err := r.DB.Exec(ctx, query, args...)
	if err != nil {
		return model.ErrorExecuteQuery
	}

	// Проверяем, что была затронута хотя бы одна строка
	if cmdTag.RowsAffected() != 1 {
		return model.ErrorNotFound
	}

	return nil
}
