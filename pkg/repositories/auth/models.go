package auth

import (
	"time"

	"github.com/google/uuid"
)

type BaseTable struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type User struct {
	BaseTable
	Email          string  `db:"email"`
	HashedPassword string  `db:"hashed_password"`
	Username       *string `db:"username"`
	ProfilePic     *string `db:"profile_picture"`
}
