package config

import (
	"io/ioutil"
	"log"

	"github.com/ansel1/merry"
	"github.com/go-yaml/yaml"
	"github.com/joshsziegler/zgo/pkg/file"
)

const (
	configPath = `config.yml`
)

type Config struct {
	// Production: If true, enables recommended security settings.
	Production bool
	Database   struct {
		// Username: MySQL user to login as.
		Username string
		// Password: The MySQL user's password.
		Password string
		// Address: Either 'localhost' or an IP such as '192.168.1.100'
		Address string
		// DBName: The name of the MySQL database to use.
		DBName string
	}
	HTTP struct {
		// ListenTo: The IP and port for the HTTP server to listen to.
		ListenTo string
	}
	LDAP struct {
		// BaseDN: The base designated name for this LDAP server.
		BaseDN string
		// UserOU: The organizational unit for users.
		UserOU string
		// GroupOU: The organizational unit for groups.
		GroupOU string
		// ListenTo: The IP and port for the LDAP server to listen to.
		ListenTo string
	}
	// SendGridAPIKey: The API key required for ending emails through SendGrid.
	SendGridAPIKey string
}

// MustLoad loads and returns our configuration from a YAML file or panic.
func MustLoad() (c Config) {
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

	// Handle required values, and defaults
	if c.Database.Username == "" {
		panic("Database Username required in " + configPath)
	}
	if c.Database.Password == "" {
		panic("Database Password required in " + configPath)
	}
	if c.Database.Address == "" {
		panic("Database Address required in " + configPath + " (e.g. localhost or 192.168.1.100)")
	}
	if c.Database.DBName == "" {
		panic("Database Name required in " + configPath)
	}
	if c.HTTP.ListenTo == "" {
		panic("HTTP ListenTo required in " + configPath + " (e.g. localhost:8080)")
	}
	// TODO: LDAP defaults and/or required args
	if c.LDAP.BaseDN == "" {
		panic("LDAP BaseDN required in " + configPath + " (e.g. dc=example,dc=com)")
	}
	if c.LDAP.UserOU == "" { // default
		c.LDAP.UserOU = "ou=People,"
	}
	if c.LDAP.GroupOU == "" { // default
		c.LDAP.GroupOU = "ou=Group,"
	}
	if c.LDAP.ListenTo == "" {
		panic("LDAP ListenTo required in " + configPath + " (e.g. localhost:3389)")
	}
	if c.SendGridAPIKey == "" {
		log.Fatal("SendGridAPIKey required in " + configPath)
	}

	// Save the resulting config to disk to preserve defaults and correct keys
	data, err := yaml.Marshal(&c)
	if err != nil {
		panic(merry.Prepend(err, "error marshalling YAML config"))
	}
	err = ioutil.WriteFile(configPath, data, 0600)
	if err != nil {
		panic(merry.Prepend(err, "error writing to "+configPath))
	}

	return c
}
