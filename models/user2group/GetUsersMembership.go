package user2group

import (
	"github.com/ansel1/merry"
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/jmoiron/sqlx"
)

// Group indicates the name and whether the User is a member or not.
type Group struct {
	Name   string `db:"Name"`
	Member bool   `db:"Member"`
}

// GetUsersMembership takes a User ID, and returns a slice of Groups, indicating
// whether that User is a member or not.
func GetUsersMembership(tx *sqlx.Tx, userID int64) (groups []Group, err error) {
	// This query uses a sub-select which returns True if the User is in the
	// Group, and False if not using COUNT(*)=1
	err = tx.Select(&groups, `SELECT Groups.Name AS Name,
							       (SELECT COUNT(*)=1
							        FROM User2Group
							        WHERE User2Group.GroupID=Groups.ID
							            AND User2Group.UserID=?) AS Member
							FROM Groups;`, userID)
	if err != nil {
		err = merry.Wrap(err)
		return
	}
	return
}
