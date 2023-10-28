package auth

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/wanderfusion/passport/pkg/repositories"
)

func (db *Database) RegisterUser(email, hashedPassword, username, profilePicture string) error {
	query := `
		INSERT INTO public.users (email, hashed_password, username, profile_picture) VALUES ($1, $2, $3, $4)
	`

	tx := db.db.MustBegin()
	_, err := tx.Exec(query, email, hashedPassword, username, profilePicture)
	if err != nil {
		isViolated, err := repositories.CheckPGUniqueConstraintError(err)
		if isViolated {
			return errors.New("User already exists")
		}
		return err
	}
	return tx.Commit()
}

func (db *Database) FetchUserDataByEmail(email string) (*User, error) {
	user := User{}
	query := `
		SELECT * FROM public.users WHERE email = $1
	`

	// QueryRow still works, but now we're scanning into multiple variables
	err := db.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows returned means there's no user with that email
			return nil, errors.New("User not found")
		}
		// Some other cosmic-level error happened, man
		return nil, err
	}

	return &user, nil
}

func (db *Database) FetchUserDataByID(id uuid.UUID) (*User, error) {
	user := User{}
	query := `
		SELECT * FROM public.users WHERE id = $1
	`

	// QueryRow still works, but now we're scanning into multiple variables
	err := db.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows returned means there's no user with that email
			return nil, errors.New("User not found")
		}
		// Some other cosmic-level error happened, man
		return nil, err
	}

	return &user, nil
}

func (db *Database) FetchUserDataByUsername(username string) (*User, error) {
	user := User{}
	query := `
		SELECT (username, profile_picture, email) FROM public.users WHERE username = $1
	`

	// QueryRow still works, but now we're scanning into multiple variables
	err := db.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows returned means there's no user with that email
			return nil, errors.New("User not found")
		}
		// Some other cosmic-level error happened, man
		return nil, err
	}

	return &user, nil
}

func (db *Database) UpdateUserProfile(user User) error {
	query := `
        UPDATE public.users
        SET 
            username = CASE WHEN $2::text != '' THEN $2::text ELSE username END,
            profile_picture = CASE WHEN $3::text != '' THEN $3::text ELSE profile_picture END,
            updated_at = NOW()
        WHERE id = $1::uuid
    `

	tx := db.db.MustBegin()
	_, err := tx.Exec(query, user.ID, user.Username, user.ProfilePic)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *Database) FetchUsersUsingUUIDs(ids []uuid.UUID) ([]User, error) {
	users := []User{}
	query := `
        SELECT * FROM public.users WHERE id IN (?)
    `

	q, vs, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}

	q = db.db.Rebind(q)

	err = db.db.Select(&users, q, vs...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (db *Database) FetchUsersUsingUsernames(usernames []string) ([]User, error) {
	users := []User{}
	query := `
        SELECT * FROM public.users WHERE username IN (?)
    `

	q, vs, err := sqlx.In(query, usernames)
	if err != nil {
		return nil, err
	}

	q = db.db.Rebind(q)

	err = db.db.Select(&users, q, vs...)
	if err != nil {
		return nil, err
	}

	return users, nil
}
