package group

import (
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/jmoiron/sqlx"
)

// Group represents and LDAP group's attributes and members
type Group struct {
	ID          int64  `db:"ID"`
	Name        string `db:"Name"`
	Description string `db:"Description"`
	UnixGroupID int64  `db:"GroupID"`
	Members     []string
}

// GetGroups returns a slice of all Groups, sans Members attribute!
//
// If you need the Members attribute, please consider using user2group.GetAll()
// as this will likely be more efficient.
func GetGroups(tx *sqlx.Tx) (groups map[int64]*(Group), err error) {
	groups = make(map[int64]*(Group))
	rows, err := tx.Queryx("SELECT * FROM Groups")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		group := new(Group)
		if err := rows.StructScan(&group); err != nil {
			return nil, err
		}
		// Add group to the map by it's DB ID
		groups[group.ID] = group
	}
	return groups, nil
}
