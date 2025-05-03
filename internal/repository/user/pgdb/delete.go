package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/pgdb/converter"
)

func (r *userRepo) DeleteUserByID(ctx context.Context, user *model.User) error {

	usr := converter.FromUserModelToRepo(user)

	query, args, err := sq.
		Delete(usersTable).
		Where(sq.Eq{userIDColumn: usr.ID}).
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
		return model.ErrorUserNotFound
	}

	return nil
}
