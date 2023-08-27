package auth

import (
	"errors"

	"github.com/akxcix/passport/pkg/repositories"
)

func (db *Database) RegisterUser(username, mail string) error {
	query := `
		INSERT INTO public.users (username, email) VALUES ($1, $2)
	`

	tx := db.db.MustBegin()
	_, err := tx.Exec(query, username, mail)
	if err != nil {
		isViolated, err := repositories.CheckPGUniqueConstraintError(err)
		if isViolated {
			return errors.New("User already exists")
		}
		return err
	}
	return tx.Commit()
}

// func (db *Database) UpdateToken(userId uuid.UUID, mail string) error {
// 	query := `
// 		INSERT INTO public.user_tokens (user_id, token) VALUES ($1, $2)
// 	`

// 	tx := db.db.MustBegin()
// 	_, err := tx.Exec(query, mail)
// 	if err != nil {
// 		return err
// 	}

// 	return tx.Commit()
// }
