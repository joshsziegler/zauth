package main

import (
	"fmt"
	"os"

	"github.com/ansel1/merry"
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/jmoiron/sqlx"
)

// DB is our shared database connection (handles connection pooling, and is
// goroutine-safe)
var DB *sqlx.DB

// connectToDB sets up the MySQL global database connection (i.e. DB).
//
// This sets several default parameters for the MySQL database:
//   - collation: Sets the charset, but avoids the additional queries of charset
//   - parseTime: changes the output type of DATE and DATETIME values to
//		 time.Time instead of []byte / string
//   - interpolateParams: Reduces the number of round trips required to
//       interpolate placeholders (i.e. ?)
//
//
// - username: Your database username
// - password: The password for this database user
// - address: The database address if not localhost (e.g. tcp(127.0.0.1))
// - name: The database to use
func connectToDB(username string, password string, address string, name string) {
	if DB != nil {
		fmt.Println("cannot connect to database (already defined)")
		return
	}

	address = fmt.Sprintf("tcp(%s)", address)
	// Create the Data Source Name (DSN), but print a password-masked version
	dsn_safe := fmt.Sprintf("%s:*****@%s/%s", username, address, name)
	log.Infof("Connecting to MySQL databse using DSN: %s", dsn_safe)
	dsn := fmt.Sprintf("%s:%s@%s/%s?%s%s%s", username, password, address, name,
		"collation=utf8mb4_general_ci&",
		"parseTime=true&",
		"interpolateParams=true")

	// SQLX connects AND pings the server, so we know the config is good or not
	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		err = merry.WithMessage(err, "could not open database connection")
		fmt.Println(err)
		os.Exit(1)
	}
}
