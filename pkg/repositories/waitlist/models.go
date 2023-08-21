package waitlist

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Mail      string    `db:"mail"`
	Name      string    `db:"name"`
}
