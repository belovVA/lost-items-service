package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/pgdb/converter"
	modelRepo "lost-items-service/internal/repository/user/pgdb/model"
)

func (r *userRepo) UserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user modelRepo.User

	query, args, err := sq.
		Select(
			userIDColumn,
			nameColumn,
			surnameColumn,
			emailColumn,
			phoneColumn,
			passwordColumn,
			roleColumn,
		).
		From(usersTable).
		Where(sq.Eq{userIDColumn: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return nil, model.ErrorUserNotFound
	}

	return converter.FromModelRepoToUser(&user), nil
}

func (r *userRepo) UserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user modelRepo.User

	query, args, err := sq.
		Select(
			userIDColumn,
			nameColumn,
			surnameColumn,
			emailColumn,
			phoneColumn,
			passwordColumn,
			roleColumn,
		).
		From(usersTable).
		Where(sq.Eq{emailColumn: email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return nil, model.ErrorUserNotFound
	}

	return converter.FromModelRepoToUser(&user), nil
}
