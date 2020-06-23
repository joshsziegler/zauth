package user

import (
	"strconv"
	"time"

	"github.com/ansel1/merry"
	"github.com/dchest/passwordreset"
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/jmoiron/sqlx"

	"github.com/joshsziegler/zauth/pkg/email"
	pw "github.com/joshsziegler/zauth/pkg/password"
	"github.com/joshsziegler/zauth/pkg/secrets"
	"github.com/joshsziegler/zgo/pkg/log"
)

// setUserPassword hashes the given cleartext password and updates the database.
//
// Warning: This does NOT check the password for strength!
func setUserPassword(tx *sqlx.Tx, username string, password string) error {
	newPasswordHash, err := pw.Hash(password)
	if err != nil {
		return merry.Wrap(err)
	}
	_, err = tx.Exec(`UPDATE Users
					  SET PasswordHash=?,
					      PasswordSet=?
					  WHERE Username=?`, newPasswordHash, time.Now(), username)
	if err != nil {
		return merry.Wrap(err)
	}

	return nil
}

// SetUserPassword checks the password's strength, and if ok, updates the
// database.
func SetUserPassword(tx *sqlx.Tx, username string, password string) error {
	// Get first and last name so we can pass to CheckPasswordRules()
	var firstName, lastName string
	err := tx.QueryRowx(`SELECT FirstName, LastName
					FROM Users
					WHERE Username=?`,
		username).Scan(&firstName, &lastName)
	if err != nil {
		return merry.Wrap(err)
	}
	// Check the user's password against our rules to see if it's too weak
	err = pw.CheckPasswordRules(username, firstName, lastName, password)
	if err != nil {
		return err
	}
	// Everything is ok, so change the password hash in the database
	err = setUserPassword(tx, username, password)
	if err != nil {
		return err
	}
	log.Infof("changed password for %s", username)
	return nil
}

type passwordResetHelper struct {
	Tx   *sqlx.Tx
	User *User
}

// Returns the concatenated password hash and password set time, for token-based
// password resets. This has been factored out to reduce the chance for errors
// when creating and checking password reset tokens.
//
// We use BOTH the password hash AND the set time to guarantee the result will
// change if the user changes their password. This will then invalidate any old
// password reset tokens still "out there."
//
// From the documentation for github.com/dchest/passwordreset:
//   Create a function that will query your users database and return some
//   password-related value for the given login. A password-related value means
//   some value that will change once a user changes their password, for
//   example: a password hash, a random salt used to generate it, or time of
//   password creation. This value, mixed with app-specific secret key, will be
//   used as a key for password reset token, thus it will be kept secret.
func getPasswordResetValue(passwordHash string, passwordSet time.Time) []byte {
	return []byte(passwordHash + passwordSet.String())
}

// GetPasswordResetValue return the password reset value for THIS user.
//
// Use user.GetPasswordResetValue(username) if you don't already have the user
// in memory.
func (u *User) GetPasswordResetValue() []byte {
	return getPasswordResetValue(u.PasswordHash, u.PasswordSet)
}

// GetPasswordResetValue return the password reset value for a user given their
// username.
//
// This queries the database, to avoid overhead. Use
// User.GetPasswordResetValue() if you already have the user in memory.
func (h *passwordResetHelper) GetPasswordResetValue(username string) (resetValue []byte, err error) {
	var passwordHash string
	var passwordSet time.Time
	err = h.Tx.QueryRowx(`SELECT PasswordHash, PasswordSet
						FROM Users
						WHERE Username=?`,
		username).Scan(&passwordHash, &passwordSet)
	if err != nil {
		err = merry.Wrap(err)
		return
	}
	resetValue = getPasswordResetValue(passwordHash, passwordSet)
	return
}

// GetPasswordResetToken returns a new token allowing the user to authenticate
// and reset their password for a limited time.
//
// The token will expire in the number of hours specified at creation.
func (u *User) GetPasswordResetToken(hours int64) string {
	expireIn := time.Duration(hours) * time.Hour
	token := passwordreset.NewToken(u.Username, expireIn,
		u.GetPasswordResetValue(), secrets.PasswordResetSecret())
	return token
}

// ValidatePasswordResetToken returns the username this password reset token
// belongs to if and only if it is valid. Otherwise it will return an error.
func ValidatePasswordResetToken(tx *sqlx.Tx, token string) (username string, err error) {
	helper := passwordResetHelper{Tx: tx, User: nil}
	username, err = passwordreset.VerifyToken(token,
		helper.GetPasswordResetValue, secrets.PasswordResetSecret())
	if err != nil {
		// Token verification failed, don't allow password reset
		return "", err
	}
	// OK, reset password for login (e.g. allow to change it)
	return
}

// SendPasswordResetEmail uses `GetPasswordResetToken` to create and send a 
// password reset link.
// 
// This uses the configured site name, URI, reply email, and reset timeout to 
// create the email. If these are incorrectly configured, this may not work!
func (u *User) SendPasswordResetEmail() error {
	// TODO: Allow these to be set in the config.yml
	siteName := "MindModeling"
	siteURI := "https://user.mindmodeling.org"
	replyEmail := "no-reply@mindmodeling.org"
	linkTime := int64(8) // In Hours

	// Generate link and send the email
	resetLink := u.GetPasswordResetToken(linkTime)
	err := email.Send(siteName, replyEmail, u.CommonName(), u.Email,
		"Your New "+siteName+" Account",
		"A new account has been created",
		`<p>Hello `+u.CommonName()+`,</p>
		<p>A new account has been created for you. Your username is <b>`+u.Username+`</b>. To complete
		the setup, you need to
		<a href="`+siteURI+`/reset-password/`+resetLink+`">
		set your password here</a>. This link is valid for the next `+strconv.FormatInt(linkTime, 10)+`
		hours.</p>`)
	return err
}
