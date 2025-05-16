package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/pgdb/converter"
)

func (r *userRepo) UpdateUserByID(ctx context.Context, user *model.User) error {

	usr := converter.FromUserModelToRepo(user)

	builder := sq.Update(usersTable).
		Set(nameColumn, usr.Name).
		Set(surnameColumn, usr.Surname).
		Set(emailColumn, usr.Email).
		Set(phoneColumn, usr.Phone)

	// Изменяем password если не пустой
	if usr.Password != "" {
		builder = builder.Set(passwordColumn, usr.Password)
	}

	// Изменяем роль при наличии
	if usr.Role != "" {
		builder = builder.Set(roleColumn, usr.Role)
	}

	builder = builder.Where(sq.Eq{userIDColumn: usr.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return model.ErrorBuildQuery
	}

	cmdTag, err := r.DB.Exec(ctx, query, args...)
	if err != nil {
		return model.ErrorExecuteQuery
	}

	if cmdTag.RowsAffected() != 1 {
		return model.ErrorNotFound
	}

	return nil

}
