package user

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ansel1/merry"
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/jmoiron/sqlx"
)

const (
	// MySQL doesn't like Go's Zero-value for Dates (0000-00-00 00:00:00) so we
	// use this instead - JZ 2019.08.22
	mysqlZeroDate = `0001-01-01 00:00:00`
)

var (
	reBadChars         = regexp.MustCompile("[^a-zA-Z0-9]+")
	ErrorLogin         = merry.New("login error")
	ErrorLoginDisabled = merry.WithMessage(ErrorLogin, "account disabled")
	ErrorLoginPassword = merry.WithMessage(ErrorLogin, "wrong password")
)

// User represents an LDAP user's attributes and group membership
//
// Assumptions:
//   - Once created, a user's username and ID will NEVER change.
//   - Only admins can create new users, change groups, and enable/disable users
//   - Enabled means that user can perform LDAP BIND operations. Disabled users
//     can still login to this website to see and change their info however.
//   - A user's UnixUserID and UnixGroupID are ALWAYS their DB ID + 1000.
type User struct {
	ID       int64  `db:"ID"` // Database ID
	Username string `db:"Username"`
	// FirstName represents the user's first name. In LDAP it's referred to as
	// their given name (givenName).
	FirstName string `db:"FirstName"`
	// LastName represents the user's last (or family) name. In LDAP it's
	// referred to as their surname (sn).
	LastName     string `db:"LastName"`
	Email        string `db:"Email"`
	PasswordHash string `db:"PasswordHash"` // SQL Default: '-'
	// Date and time when was this password last set or changed.
	PasswordSet time.Time `db:"PasswordSet"` // SQL Default: 0001-01-01 00:00:00
	// Date and time when this user last logged in.
	LastLogin time.Time `db:"LastLogin"` // SQL Default: 0001-01-01 00:00:00
	// If disabled, LDAP binds for this account will fail. Logins to zauth's
	// user management page will continue to work however!
	Disabled bool `db:"Disabled"` // If true, don't allow to login
	Groups   []string
}

// CommonName is the user's full name (returns the first and last names).
//
// The name of this function is a reference to LDAP's terminology 'cn' for the
// full name of a user (LDAP uses 'sn' or Surname, and 'givenname' as the first
// name).
//
// ** Doesn't use a pointer to `u` so it can be use in HTML templates.
func (u User) CommonName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// UnixUserID returns their Unix ID, which is always their database ID + 1000.
// This assumes that regular user accounts start at 1000.
//
// ** Doesn't use a pointer to `u` so it can be use in HTML templates.
func (u User) UnixUserID() int64 {
	return u.ID + 1000
}

// UnixGroupID returns the same value as UnixUserID, which assumes that they
// belong to their own group.
//
// ** Doesn't use a pointer to `u` so it can be use in HTML templates.
func (u User) UnixGroupID() int64 {
	return u.ID + 1000
}

// HomeDirectory returns their Unix directory as "/home/username"
//
// ** Doesn't use a pointer to `u` so it can be use in HTML templates.
func (u User) HomeDirectory() string {
	return fmt.Sprintf("/home/%s", u.Username)
}

// IsAdmin returns true if this User belongs to a group named 'admin'.
//
// ** Doesn't use a pointer to `u` so it can be use in HTML templates.
func (u User) IsAdmin() bool {
	for _, name := range u.Groups {
		if strings.ToLower(name) == "admin" {
			return true
		}
	}
	return false
}

// userSetEnable is a helper function for UserEnable and UserDisable.
//
// Note that isEnabled is flipped because the database uses Disabled!
func userSetEnable(isEnabled bool, username string) error {
	_, err := db.Exec(`UPDATE Users
		 			  SET Disabled=?
		 			  WHERE Username=?`, !isEnabled, username)
	if err != nil {
		merry.Wrap(err)
	}
	return nil
}

func UserEnable(username string) (err error) {
	return userSetEnable(true, username)
}

func UserDisable(username string) (err error) {
	return userSetEnable(false, username)
}

// GetUserWithGroups returns a single User struct, including the groups they
// belong to (in alphabetical ascending order by name).
func GetUserWithGroups(tx *sqlx.Tx, username string) (user User, err error) {
	err = tx.QueryRowx(`SELECT * FROM Users WHERE Username=?`, username).StructScan(&user)
	if err != nil {
		return User{}, merry.Wrap(err)
	}
	// Get the name of each group this user belongs to
	rows, err := tx.Queryx(`SELECT Groups.Name 
							FROM Groups 
							INNER JOIN User2Group 
								ON Groups.ID=User2Group.GroupID 
							WHERE User2Group.UserID=?
							ORDER BY Groups.Name ASC;`, user.ID)
	if err != nil {
		return User{}, merry.Wrap(err)
	}
	defer rows.Close()
	var groupName string
	for rows.Next() {
		err = rows.Scan(&groupName)
		if err != nil {
			return User{}, merry.Wrap(err)
		}
		user.Groups = append(user.Groups, groupName)
	}
	return
}

// GetUsersWithoutGroups returns a map of users, stored by their database ID.
func GetUsersMapWithoutGroups(tx *sqlx.Tx) (users map[int64]*(User), err error) {
	users = make(map[int64]*(User))
	rows, err := tx.Queryx(`SELECT * FROM Users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(User)
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
		// Add user to the map by it's DB ID
		users[user.ID] = user
	}
	return users, nil
}

// GetGroupsNotMemberOf returns a slice of all Group names this user is NOT a
// member of.
func (u *User) GetGroupsNotMemberOf(tx *sqlx.Tx) (groups []string, err error) {
	groups = make([]string, 0)
	var rows *sqlx.Rows
	// If they aren't in any groups, the other query will fail
	if len(u.Groups) < 1 {
		rows, err = tx.Queryx(`SELECT Name
								FROM Groups
								ORDER BY Name ASC;`)

	} else {
		// Create a query with an arbitrary number of params using sqlx.In()
		query, qArgs, err := sqlx.In(`SELECT Name 
									  FROM Groups
									  WHERE Name NOT IN(?)
									  ORDER BY Name ASC;`, u.Groups)
		if err != nil {
			return groups, merry.Wrap(err)
		}
		// Run the actual query
		rows, err = tx.Queryx(query, qArgs...)
		if err != nil {
			return groups, merry.Wrap(err)
		}
	}
	defer rows.Close()
	// Parse each Group name
	var groupName string
	for rows.Next() {
		err = rows.Scan(&groupName)
		if err != nil {
			return groups, merry.Wrap(err)
		}
		groups = append(groups, groupName)
	}
	return groups, nil
}
