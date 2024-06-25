package main

import (
	"io/ioutil"

	"github.com/ansel1/merry"
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	"github.com/go-yaml/yaml"
	"github.com/jmoiron/sqlx"

	"github.com/joshsziegler/zauth/models/httpserver"
	"github.com/joshsziegler/zauth/pkg/db"
	"github.com/joshsziegler/zauth/pkg/email"
	"github.com/joshsziegler/zauth/pkg/ldap"
	"github.com/joshsziegler/zgo/pkg/file"
	"github.com/joshsziegler/zgo/pkg/log"
)

const (
	configPath  = `config.yml`
	programName = `zauth`
)

var (
	// Build and Version info (given at compile time via build.go)
	Version   string
	BuildDate string
	// We use a global config, because it should be read-only after initial loading
	config Config
	// DB is our shared database connection (handles connection pooling, and is
	// goroutine-safe)
	DB *sqlx.DB
)

type httpConfig struct {
	ListenTo string
}

// Config stores the all server options.
type Config struct {
	Production     bool
	Database       db.Config
	LDAP           ldap.Config
	HTTP           httpConfig
	SendGridAPIKey string
}

// mustLoadConfig loads and returns our configuration from a YAML file or panic.
func mustLoadConfig() (c Config) {
	// Read the existing config file from disk
	if file.Exists(configPath) {
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(merry.Prepend(err, "error reading "+configPath))
		}
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			panic(merry.Prepend(err, "error parsing YAML from "+configPath))
		}
	}

	return c
}

func main() {
	log.Infof("%s %s (Built: %s)", programName, Version, BuildDate)
	config = mustLoadConfig()
	DB = db.MustConnect(config.Database)
	email.Init(config.SendGridAPIKey)
	go httpserver.Listen(DB, config.HTTP.ListenTo, config.Production)
	ldap.Listen(DB, config.LDAP) // blocking
}
