package modelredis

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `redis:"id"`
	Name     string    `redis:"name"`
	Surname  string    `redis:"surname"`
	Email    string    `redis:"email"`
	Phone    string    `redis:"phone"`
	Password string    `redis:"password"`
	Role     string    `redis:"role"`
}
