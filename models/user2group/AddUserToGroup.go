package user2group

import (
	"github.com/ansel1/merry"
	"github.com/jmoiron/sqlx"
)

// setUserGroupMembership is a helper function for AddUserToGroup and
// RemoveUserFromGroup. If add is true, it adds the user to the group. If add is
// false, it removes the user from the group.
func setUserGroupMembership(tx *sqlx.Tx, user string, group string, add bool) error {
	// 1. Get the UserID and GroupID here once, instead of doing two SQL JOINS
	var userID, groupID uint64
	err := tx.Get(&userID, `SELECT ID FROM Users WHERE Username=?;`, user)
	if err != nil {
		return merry.Wrap(err)
	}
	err = tx.Get(&groupID, `SELECT ID FROM Groups WHERE Name=?;`, group)
	if err != nil {
		return merry.Wrap(err)
	}
	// 2. Figure out if the user is currently in the group
	var inGroup bool
	err = tx.Get(&inGroup, `SELECT (COUNT(*)=1)
						    FROM User2Group 
						    WHERE UserID=? AND GroupID=?;`, userID, groupID)
	if err != nil {
		return merry.Wrap(err)
	}
	// 3. Add or Remove them from the group
	if add && !inGroup {
		_, err = tx.Exec(`INSERT INTO User2Group (UserID, GroupID) 
						   VALUES (?, ?);`, userID, groupID)
		if err != nil {
			return merry.Wrap(err)
		}
	} else if !add && inGroup {
		_, err = tx.Exec(`DELETE FROM User2Group 
						  WHERE UserID=? AND GroupID=?;`, userID, groupID)
		if err != nil {
			return merry.Wrap(err)
		}
	}
	return nil
}

// AddUserToGroup adds the User to a Group.
func AddUserToGroup(tx *sqlx.Tx, user string, group string) error {
	return setUserGroupMembership(tx, user, group, true)
}

// RemoveUserFromGroup removes the User from a Group.
func RemoveUserFromGroup(tx *sqlx.Tx, user string, group string) error {
	return setUserGroupMembership(tx, user, group, false)
}
