package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/pgdb/converter"
)

func (r *userRepo) AddUser(ctx context.Context, user *model.User) (uuid.UUID, error) {
	var id uuid.UUID
	usr := converter.FromUserModelToRepo(user)

	query, args, err := sq.
		Insert(usersTable).
		Columns(nameColumn, surnameColumn, emailColumn, phoneColumn, passwordColumn, roleColumn).
		Values(usr.Name, usr.Surname, usr.Email, usr.Phone, usr.Password, usr.Role).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING " + userIDColumn).
		ToSql()

	if err != nil {
		return uuid.Nil, model.ErrorFailedBuildQuery
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, model.ErrorFailedBuildQuery
	}

	return id, nil
}
