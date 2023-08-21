package waitlist

import "github.com/akxcix/passport/pkg/repositories"

func (db *Database) AddUser(mail, name string) error {
	query := `
		INSERT INTO public.waitlist (mail, name) VALUES ($1, $2)
	`

	tx := db.db.MustBegin()
	_, err := tx.Exec(query, mail, name)
	if err != nil {
		return repositories.CheckPGUniqueConstraintError(err)
	}
	return tx.Commit()
}
