package repositories

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
)

func FormatPostgresDSN(user, pwd, host, port, db string) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", user, pwd, host, port, db)
}

func CheckPGUniqueConstraintError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" { // Unique constraint violation
			return errors.New("email already exists in waitlist")
		}
	}
	return err
}
