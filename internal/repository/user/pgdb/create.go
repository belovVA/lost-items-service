package pgdb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/pgdb/converter"
)

func (r *userRepo) CreateUser(ctx context.Context, model *model.User) (uuid.UUID, error) {
	var id uuid.UUID
	user := converter.FromUserModelToRepo(model)

	query, args, err := sq.
		Insert(usersTable).
		Columns(nameColumn, surnameColumn, emailColumn, phoneColumn, passwordColumn, roleColumn).
		Values(user.Name, user.Surname, user.Email, user.Phone, user.Password, user.Role).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING " + userIDColumn).
		ToSql()

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", FailedBuildQuery, err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %s", FailedCreateUser, err.Error())
	}

	return id, nil
}
