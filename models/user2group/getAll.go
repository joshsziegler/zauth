package user2group

import (
	"github.com/ansel1/merry"
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/jmoiron/sqlx"

	mGroup "github.com/joshsziegler/zauth/models/group"
	mUser "github.com/joshsziegler/zauth/models/user"
)

// GetAll Users and Groups, WITH membership info populated.
//
// This exists because it *should* be more efficient for populating group
// membership info IF AND ONLY IF you need all or most of the users and groups.
func GetAll(tx *sqlx.Tx) (users map[int64]*(mUser.User),
	groups map[int64]*(mGroup.Group), err error) {

	// Get all Users (does NOT pull group membership)
	users, err = mUser.GetUsersMapWithoutGroups(tx)
	if err != nil {
		err = merry.Append(err, "error getting users")
		return
	}

	// Get all Groups (does NOT pull members)
	groups, err = mGroup.GetGroupsMapWithoutUsers(tx)
	if err != nil {
		err = merry.Append(err, "error getting groups")
		return
	}

	// Map Users to Groups and vice versa using the User2Group table directly
	// This *should* be more efficient than query this table for each User and
	// Group individually
	u2g := user2Group{}
	rows, err := tx.Queryx("SELECT * FROM User2Group")
	if err != nil {
		err = merry.Append(err, "error querying User2Group table")
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&u2g)
		if err != nil {
			err = merry.Append(err, "error scanning into User2Group")
			return
		}
		// Update the User record
		users[u2g.UserID].Groups = append(users[u2g.UserID].Groups,
			groups[u2g.GroupID].Name)
		// Update the Group record
		groups[u2g.GroupID].Members = append(groups[u2g.GroupID].Members,
			users[u2g.UserID].Username)
	}
	return
}

// Represents the Many-to-Many DB table mapping Users to Groups
type user2Group struct {
	UserID  int64 `db:"UserID"`
	GroupID int64 `db:"GroupID"`
}
