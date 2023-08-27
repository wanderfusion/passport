package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `db:"id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	Username       string    `db:"username"`
	Email          string    `db:"email"`
	FullName       string    `db:"full_name"`
	DisplayPicture string    `db:"display_picture"`
}
