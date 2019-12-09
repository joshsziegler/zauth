package group

import (
	"database/sql"
	"regexp"

	"github.com/ansel1/merry"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/jmoiron/sqlx"
)

// TODO: Restrict Group names to alphanumeric, with hyphens (no whitespace)
// TODO: Require group names to be unique
var (
	// reValidName represents the POSIX standard for valid user, group, and
	// file names. This definition comes from: https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap03.html#tag_03_282
	//
	// That information from that link is reproduced here (as of 2019.11.15):
	//    3.282 Portable Filename Character Set
	//
	//    The set of characters from which portable filenames are constructed.
	//
	//    A B C D E F G H I J K L M N O P Q R S T U V W X Y Z
	//    a b c d e f g h i j k l m n o p q r s t u v w x y z
	//    0 1 2 3 4 5 6 7 8 9 . _ -
	//
	//    The last three characters are the <period>, <underscore>, and
	//    <hyphen-minus> characters, respectively. See also Pathname.
	//
	reValidName = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,31}$`)
)

// Group represents and LDAP group's attributes and members
type Group struct {
	ID          int64  `db:"ID"`
	Name        string `db:"Name"`
	Description string `db:"Description"`
	Members     []string
}

// UnixGroupID is always their database ID + 100.
// This assumes that regular groups start at 100.
//
// ** Doesn't use a pointer to `u` so it can be use in HTML templates.
func (g Group) UnixGroupID() int64 {
	return g.ID + 100
}

func GetGroupsSliceWithoutUsers(tx *sqlx.Tx) (groups []*Group, err error) {
	err = tx.Select(&groups, "SELECT ID, Name, Description FROM Groups ORDER BY Name ASC")
	if err != nil {
		err = merry.WithMessage(err, "error retrieving groups list from database")
		return
	}
	return
}

// GetGroupsMapWithoutUsers returns a map of all Groups (using their DB ID as
// the key), sans Members attribute.
//
// If you need the Members attribute, please consider using user2group.GetAll()
// as this will likely be more efficient.
func GetGroupsMapWithoutUsers(tx *sqlx.Tx) (groups map[int64]*(Group), err error) {
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

// Add inserts a new group into the database.
//
// TODO: Add name validity checks? - JZ
func Add(tx *sqlx.Tx, name string, description string) error {
	_, err := tx.Exec("INSERT INTO Groups (Name, Description) VALUES (?,?);",
		name, description)
	if err != nil {
		sqlError, ok := err.(*mysql.MySQLError)
		if ok {
			if sqlError.Number == 1062 {
				return merry.New("a group with that name already exists")
			}
		}
		return merry.Wrap(err)
	}
	return nil
}

func Delete(tx *sql.Tx, name string) error {
	_, err := tx.Exec("DELETE FROM Groups WHERE Name IS ?;", name)
	if err != nil {
		return merry.Wrap(err)
	}
	return nil
}
