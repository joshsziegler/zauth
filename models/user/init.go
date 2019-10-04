package user

import (
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
)

// db is the shared (via connection pooling) database connection; goroutine-safe
var db *sqlx.DB
var log *logging.Logger

// Init sets the logger, and database connection object for internal use
func Init(logger *logging.Logger, database *sqlx.DB) {
	log = logger
	db = database
}
