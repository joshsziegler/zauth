package user

import (
	"time"

	"github.com/ansel1/merry"
	"github.com/jmoiron/sqlx"

	"github.com/joshsziegler/zauth/pkg/log"
	pw "github.com/joshsziegler/zauth/pkg/password"
)

// Login returns nil IFF the account is not disabled AND the password is correct
func Login(tx *sqlx.Tx, username string, password string) (err error) {
	var correctPasswordHash string
	var disabled bool
	err = tx.QueryRowx(`SELECT PasswordHash, Disabled
						FROM Users
						WHERE Username=?`,
		username).Scan(&correctPasswordHash, &disabled)
	if err != nil {
		return merry.Wrap(err)
	}

	if disabled {
		return ErrorLoginDisabled.Here().WithMessagef("user '%s' is disabled", username)
	}

	valid, insecure, err := pw.Valid(password, correctPasswordHash)
	if err != nil {
		return merry.Wrap(err)
	}
	if !valid {
		return ErrorLoginPassword.Here().WithMessagef("wrong password for '%s'", username)
	}

	// Update LastLogin
	_, err = tx.Exec(`UPDATE Users
		 			  SET LastLogin=?
		 			  WHERE Username=?`, time.Now(), username)
	if err != nil {
		return merry.Wrap(err)
	}

	// Update PasswordHash IFF it's using an insecure hashing method (e.g. MD5)
	if insecure {
		err = setUserPassword(tx, username, password)
		if err != nil {
			return err // already wrapped
		}
		log.Infof("upgraded password hash for: %s", username)
	}
	return nil
}
