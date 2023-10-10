package db

import (
	"fmt"

	"github.com/ansel1/merry"
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/jmoiron/sqlx"

	"github.com/joshsziegler/zgo/pkg/log"
)

// Config holds the database configuration, required for connecting to and
// using the database.
type Config struct {
	// Username: MySQL user to login as.
	Username string
	// Password: The MySQL user's password.
	Password string
	// Address: Either 'localhost' or an IP such as '192.168.1.100'
	Address string
	// DBName: The name of the MySQL database to use.
	DBName string
}

// MustConnect connects and returns a database connection, or calls os.Exit(1).
//
// This sets several default parameters for the MySQL database:
//   - collation: Sets the charset, but avoids the additional queries of charset
//   - parseTime: changes the output type of DATE and DATETIME values to
//     time.Time instead of []byte / string
//   - interpolateParams: Reduces the number of round trips required to
//     interpolate placeholders (i.e. ?)
func MustConnect(config Config) *sqlx.DB {
	config.Address = fmt.Sprintf("tcp(%s)", config.Address)
	// Create the Data Source Name (DSN), but print a password-masked version
	dsn_safe := getDSN(config.Username, "*****", config.Address, config.DBName)
	log.Infof("Connecting to MySQL databse using DSN: %s", dsn_safe)
	dsn := getDSN(config.Username, config.Password, config.Address, config.DBName)
	// SQLX connects AND pings the server, so we know the config is good or not
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		err = merry.WithMessage(err, "could not open database connection")
		log.Fatal(err)
		return nil
	}
	return db
}

// getDSN returns the full MySQL DSN.
//
// This was factored out of MustConnect so we could use it for the logger-safe
// version (with the password masked out) and the real Data Source Name (DSN).
func getDSN(username, password, address, dbName string) string {
	return fmt.Sprintf("%s:%s@%s/%s?%s%s%s", username, password, address, dbName,
		"collation=utf8mb4_general_ci&",
		"parseTime=true&",
		"interpolateParams=true")
}
