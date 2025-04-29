package modelRepo

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Surname  string    `db:"surname"`
	Email    string    `db:"email"`
	Phone    string    `db:"phone"`
	Password string    `db:"password"`
	Role     string    `db:"role"`
}
