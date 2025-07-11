package pgdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lost-items-service/internal/model"
	"lost-items-service/internal/repository/user/pgdb/converter"
	modelpg "lost-items-service/internal/repository/user/pgdb/model"
)

func (r *userRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user modelpg.User

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
		return nil, model.ErrorNotFound
	}

	return converter.FromModelRepoToUser(&user), nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user modelpg.User

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
	if err != nil {
		return nil, model.ErrorBuildQuery
	}
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
		return nil, model.ErrorNotFound
	}

	return converter.FromModelRepoToUser(&user), nil
}

func (r *userRepo) GetListUsers(ctx context.Context, info *model.InfoSetting) ([]*model.User, error) {
	limits := converter.FromInfoUsersToLimitsUsers(info)
	users := make([]*model.User, 0, info.Limit)

	//// -------- COUNT QUERY --------
	//countReq := sq.Select("COUNT(*)").From(usersTable).PlaceholderFormat(sq.Dollar)
	//if limits.OrderByField != "" {
	//	countReq = countReq.Where(sq.Eq{roleColumn: limits.OrderByField})
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
		PlaceholderFormat(sq.Dollar)
	if limits.Role != "" {
		req = req.Where(sq.Eq{roleColumn: limits.Role})
	}
	if limits.Search != "" {
		searchPattern := "%" + limits.Search + "%"
		req = req.Where(sq.Or{
			sq.ILike{nameColumn: searchPattern},
			sq.ILike{surnameColumn: searchPattern},
			sq.ILike{emailColumn: searchPattern},
			sq.ILike{phoneColumn: searchPattern},
			sq.ILike{roleColumn: searchPattern},
		})
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
		var user modelpg.User
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
