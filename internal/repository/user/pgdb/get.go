package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/pgdb/converter"
	modelRepo "lost-items-service/internal/repository/user/pgdb/model"
)

func (r *userRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
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

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
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

func (r *userRepo) GetListUsers(ctx context.Context, info *model.InfoUsers) ([]*model.User, error) {
	limits := converter.FromInfoUsersToLimitsUsers(info)
	users := make([]*model.User, 0, info.Limit)

	//// -------- COUNT QUERY --------
	//countReq := sq.Select("COUNT(*)").From(usersTable).PlaceholderFormat(sq.Dollar)
	//if limits.Role != "" {
	//	countReq = countReq.Where(sq.Eq{roleColumn: limits.Role})
	//}
	//
	//countQuery, countArgs, err := countReq.ToSql()
	//if err != nil {
	//	return nil, 0, fmt.Errorf("%s: %w", FailedBuildQuery, err)
	//}
	//
	//var totalCount int
	//if err := r.DB.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount); err != nil {
	//	return nil, 0, fmt.Errorf("%s: %w", FailedExecuteQuery, err)
	//}

	req := sq.
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
		PlaceholderFormat(sq.Dollar).
		OrderBy(userIDColumn)
	if limits.Role != "" {
		req = req.Where(sq.Eq{roleColumn: limits.Role})
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
		var user modelRepo.User
		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Phone,
			&user.Password,
			&user.Role,
		); err != nil {
			return nil, model.ErrorScanRows
		}

		users = append(users, converter.FromModelRepoToUser(&user))
	}

	if rows.Err() != nil {
		return nil, model.ErrorScanRows
	}
	return users, nil
}
