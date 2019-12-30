package user

import (
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/jmoiron/sqlx"
)

// db is the shared (via connection pooling) database connection; goroutine-safe
var db *sqlx.DB

// Init sets the logger, and database connection object for internal use
func Init(database *sqlx.DB) {
	db = database
}
