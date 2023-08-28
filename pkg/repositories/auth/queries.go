package auth

import (
	"database/sql"
	"errors"

	"github.com/akxcix/passport/pkg/repositories"
	"github.com/google/uuid"
)

func (db *Database) RegisterUser(email, hashedPassword string) error {
	query := `
		INSERT INTO public.users (email, hashed_password) VALUES ($1, $2)
	`

	tx := db.db.MustBegin()
	_, err := tx.Exec(query, email, hashedPassword)
	if err != nil {
		isViolated, err := repositories.CheckPGUniqueConstraintError(err)
		if isViolated {
			return errors.New("User already exists")
		}
		return err
	}
	return tx.Commit()
}

func (db *Database) FetchUserDataByEmail(email string) (uuid.UUID, string, error) {
	var passwordHash string
	var id uuid.UUID
	query := `
		SELECT id, hashed_password FROM public.users WHERE email = $1
	`

	// QueryRow still works, but now we're scanning into multiple variables
	err := db.db.QueryRow(query, email).Scan(&id, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows returned means there's no user with that email
			return uuid.UUID{}, "", errors.New("User not found")
		}
		// Some other cosmic-level error happened, man
		return uuid.UUID{}, "", err
	}

	return id, passwordHash, nil
}
