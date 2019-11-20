package main

import (
	"io/ioutil"

	"github.com/ansel1/merry"
	_ "github.com/go-sql-driver/mysql" // Blank import required for SQL drivers
	yaml "gopkg.in/yaml.v2"

	"github.com/joshsziegler/zauth/models/file"
	"github.com/joshsziegler/zauth/models/ldapserver"
	"github.com/joshsziegler/zauth/pkg/db"
)

const (
	configPath = `config.yml`
)

type httpConfig struct {
	ListenTo string
}

// Config stores the server options such as IP/NIC to listen on, port and
// Hash/Block Key (for secure cookies).
type Config struct {
	Production     bool
	Database       db.Config
	LDAP           ldapserver.LdapConfig
	HTTP           httpConfig
	SendGridAPIKey string
}

// MustLoadConfigs loads and returns our configuration from a YAML file or panic.
//
// - If no value is given in that file, defaults are used/created as needed.
// - The config is always written back to disk in order to preserve the Hash
//   and Block Keys if they are generated.
func MustLoadConfigs() (c Config) {
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
