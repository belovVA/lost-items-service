package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/pgdb/converter"
)

func (r *userRepo) UpdateUserByID(ctx context.Context, user *model.User) error {

	usr := converter.FromUserModelToRepo(user)

	query, args, err := sq.
		Update(usersTable).
		Set(nameColumn, usr.Name).
		Set(surnameColumn, usr.Surname).
		Set(emailColumn, usr.Email).
		Set(phoneColumn, usr.Phone).
		Set(passwordColumn, usr.Password).
		Set(roleColumn, usr.Role).
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
