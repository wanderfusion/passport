package auth

import (
	"database/sql"
	"errors"

	"github.com/akxcix/passport/pkg/repositories"
)

func (db *Database) RegisterUser(username, hashedPassword string) error {
	query := `
		INSERT INTO public.users (username, hashed_password) VALUES ($1, $2)
	`

	tx := db.db.MustBegin()
	_, err := tx.Exec(query, username, hashedPassword)
	if err != nil {
		isViolated, err := repositories.CheckPGUniqueConstraintError(err)
		if isViolated {
			return errors.New("User already exists")
		}
		return err
	}
	return tx.Commit()
}

func (db *Database) FetchHashByUsername(username string) (string, error) {
	var passwordHash string
	query := `
		SELECT hashed_password FROM public.users WHERE username = $1
	`

	// You can use QueryRow when expecting a single row
	err := db.db.QueryRow(query, username).Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows returned means there's no user with that username
			return "", errors.New("User not found")
		}
		// Some other error occurred
		return "", err
	}

	return passwordHash, nil
}
