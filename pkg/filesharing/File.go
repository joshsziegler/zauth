package filesharing

import (
	"regexp"
	"time"

	"github.com/ansel1/merry"
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/jmoiron/sqlx"
)

var (
	reBadChars         = regexp.MustCompile("[^a-zA-Z0-9]+")
	ErrorLogin         = merry.New("login error")
	ErrorLoginDisabled = merry.WithMessage(ErrorLogin, "account disabled")
	ErrorLoginPassword = merry.WithMessage(ErrorLogin, "wrong password")
)

// File represents an single shared file.
//
// Assumptions:
//   - Filenames MUST be unique for the folder they are current in.
type File struct {
	ID        int64     `db:"ID"`   // Database ID
	Name      string    `db:"Name"` // Name is the user-supplied filename, which may be unsafe!
	FileSize  int64     `db:"FileSize"`
	Digest    string    `db:"Digest"`    // Hash of this file to verify content
	CreatedAt time.Time `db:"CreatedAt"` // SQL Default: Current Time
	CreatedBy string    `db:"CreatedBy"` // SQL stored the ID, we get the Username
	// UpdatedAt time.Time `db:"UpdatedAt"` // SQL Default: Current Time; Set on Updates
	// UpdatedBy string // SQL stored the ID, we fetch the User's name instead
}

// GetFiles returns a sorted list of files in the Group -- or a specific folder
// -- sorted by Name.
func GetFiles(tx *sqlx.Tx, groupID int64) (files [](File), err error) {
	files = make([](File), 0)
	rows, err := tx.Queryx(`SELECT f.ID, f.Name, f.FileSize, f.Digest, f.CreatedAt, 
		                           u.Username AS CreatedBy
		                    FROM Files f 
		                    LEFT JOIN Users u ON f.CreatedBy=u.ID
		                    WHERE groupID=?
		                    ORDER BY Name ASC`, groupID)
	if err != nil {
		return files, err
	}
	defer rows.Close()

	for rows.Next() {
		f := File{}
		if err := rows.StructScan(&f); err != nil {
			return files, err
		}
		files = append(files, f)
	}
	return files, nil
}
