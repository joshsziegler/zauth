package mysql

import (
	"github.com/ansel1/merry"
)

// The following covers a small set of the MySQL numerical error codes.
// See https://dev.mysql.com/doc/refman/8.0/en/server-error-reference.html for
// more information.
const (
	// Indicates a 'unique' constraint violation (during an insert)
	// Message: Duplicate entry '%s' for key %d
	errDuplicateEntry = 1062
	// Indicates the query returned no rows.
	// Message: Query was empty
	errResultSetEmpty = 1065
)

var (
	// ErrDuplicateEntry indicates a 'unique' constraint violation (typically
	// during an insert).
	ErrDuplicateEntry = merry.New("duplicate sql entry")
	// ErrResultSetEmpty indicates the query returned no rows (was empty).
	ErrResultSetEmpty = merry.New("query returned no rows")
)

// func Error(err *error) bool {
// 	if err != nil {
// 		return true
// 	}
// 	sqlError, ok := err.(*mysql.MySQLError)
// 	if ok {
// 		switch sqlError.Number {
// 		case errDuplicateEntry:
// 			return merry.WithMessage(merry.Here(ErrDuplicateEntry, err))
// 		case errResultSetEmpty:
// 			return merry.WithMessage(merry.Here(ErrResultSetEmpty, err))
// 		default:
// 			// Unknown and unhandled SQL error
// 		}
// 	}
// 	return merry.Wrap(err)
//
// }
