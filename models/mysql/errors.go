package mysql

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
