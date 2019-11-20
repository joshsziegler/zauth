package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	logging "github.com/op/go-logging"

	"github.com/joshsziegler/zauth/models/email"
	"github.com/joshsziegler/zauth/models/httpserver"
	"github.com/joshsziegler/zauth/models/ldapserver"
	"github.com/joshsziegler/zauth/models/user"
	"github.com/joshsziegler/zauth/pkg/ztk/db"
)

const programName = `zauth`

// Build and Version info (given at compile time via build.go)
var Version string
var BuildDate string

// We use a global config, because it should be read-only after initial loading
var config Config
var log = logging.MustGetLogger(programName)

// DB is our shared database connection (handles connection pooling, and is
// goroutine-safe)
var DB *sqlx.DB

// initLogging sets up logging to stdout
func initLogging() {
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
	logging.SetLevel(logging.INFO, programName)
	format := "%{level:.4s} ▶ %{message}"
	logging.SetFormatter(logging.MustStringFormatter(format))
}

func main() {
	initLogging()
	log.Info(fmt.Sprintf("%s %s (Built: %s)", programName, Version, BuildDate))
	config = MustLoadConfigs()
	DB = db.MustConnect(log, config.Database)
	user.Init(log, DB)
	email.Init(config.SendGridAPIKey)
	go httpserver.Listen(log, DB, config.HTTP.ListenTo, config.Production)
	ldapserver.Listen(log, DB, config.LDAP) // blocking
}
