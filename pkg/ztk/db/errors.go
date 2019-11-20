package db

// The following covers a small set of the MySQL numerical error codes.
// See https://dev.mysql.com/doc/refman/8.0/en/server-error-reference.html for
// more information.
const (
	// ErrDuplicateEntry indicates a 'unique' constraint violation on insert.
	// Message: Duplicate entry '%s' for key %d
	ErrDuplicateEntry = 1062
	// ErrNoResults indicates the MySQL query returned no rows.
	// Message: Query was empty
	ErrNoResults = 1065
)
