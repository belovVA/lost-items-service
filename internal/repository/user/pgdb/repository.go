package pgdb

import (
	"context"

	"github.com/google/uuid"
	"lost-items-service/internal/db/pgxdb"
	"lost-items-service/internal/model"
)

const (
	usersTable     = "users"
	userIDColumn   = "id"
	nameColumn     = "name"
	surnameColumn  = "surname"
	emailColumn    = "email"
	phoneColumn    = "phone"
	passwordColumn = "password"
	roleColumn     = "role"
)

const (
	FailedBuildQuery = "failed Build Query"
	FailedCreateUser = "failed to Create User"
	UserNotFound     = "user not found"
)

type UserPGRepository interface {
	CreateUser(ctx context.Context, model *model.User) (uuid.UUID, error)
	UserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type userRepo struct {
	DB pgxdb.DB
}

func NewRepository(db pgxdb.DB) UserPGRepository {
	return &userRepo{
		DB: db,
	}
}
