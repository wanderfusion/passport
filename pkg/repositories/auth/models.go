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
	Username       string `db:"username"`
	HashedPassword string `db:"hashed_password"`
}
