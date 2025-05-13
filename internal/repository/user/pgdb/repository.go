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

type UserPGRepository interface {
	AddUser(ctx context.Context, user *model.User) (uuid.UUID, error)

	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetListUsers(ctx context.Context, info *model.InfoUsers) ([]*model.User, error)

	UpdateUserByID(ctx context.Context, user *model.User) error

	DeleteUserByID(ctx context.Context, user *model.User) error
}

type userRepo struct {
	DB pgxdb.DB
}

func NewRepository(db pgxdb.DB) UserPGRepository {
	return &userRepo{
		DB: db,
	}
}
