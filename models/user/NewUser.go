package user

import (
	"fmt"
	"strings"

	"github.com/ansel1/merry"
	"github.com/badoux/checkmail"
	"github.com/jmoiron/sqlx"
)

// NewUser creates a new user (if details are valid), and send them an email so
// they can set their initial password.
func NewUser(DB *sqlx.DB, firstName string, lastName string, email string) (user User, err error) {
	// 0. Create temp user object to hold our values pre-Db-insert
	firstName = strings.Trim(firstName, " ")
	lastName = strings.Trim(lastName, " ")
	email = strings.Trim(email, " ")

	// 1. Validate inputs
	if len(firstName) < 1 || len(lastName) < 1 {
		err = merry.New("FirstName or LastName < 1 character").
			WithUserMessage("First and last name are required.")
		return
	}
	err = checkmail.ValidateFormat(email)
	if err != nil {
		err = merry.Wrap(err).WithUserMessage("Email must be valid.")
		return
	}

	// 2. Create username and home directory (based on first and last name)
	firstName = reBadChars.ReplaceAllString(firstName, "")
	lastName = reBadChars.ReplaceAllString(lastName, "")
	username := strings.ToLower(fmt.Sprintf("%s.%s", firstName, lastName))

	// 3. Insert the user into the DB
	_, err = DB.Exec(`INSERT INTO Users (Username, FirstName, LastName, Email)
					  VALUES (?,?,?,?)`, username, firstName, lastName, email)
	if err != nil {
		// TODO: Handle duplicate username errors differently? - JZ 2019.08.23
		err = merry.Wrap(err).WithUserMessage("Database insertion failed.")
		return
	}

	// 4. Get and return the user
	user, err = GetUserWithGroups(DB, username)
	if err != nil {
		err = merry.Wrap(err).WithUserMessage("Retrieving new user failed.")
		return
	}

	return
}
