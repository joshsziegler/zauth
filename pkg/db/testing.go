//+build !test

// The line above will exclude this file from builds, except during testing.
// This file contains methods useful for doing *real* database tests, rather
// than using mocks.
//

package db

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joshsziegler/zauth/pkg/log"
)

// GetTxOrFailTesting creates and returns a database transaction, or will end
// the unit test by calling t.Fatal()
func GetTxOrFailTesting(t *testing.T, db *sqlx.DB) *sqlx.Tx {
	tx, err := db.Beginx()
	if err != nil {
		log.Fatalf("error creating DB transaction: %v", err)
		return nil
	}
	return tx
}

// SetupTestingDatabase creates and returns an empty database, using the the
// database name foo_test, where 'foo' is the name in the config.
func SetupTestingDatabase(t *testing.T, config Config, scriptPath string) *sqlx.DB {
	// 1. Use the default database name but add _test to avoid overwriting data
	config.DBName += "_test"
	// 2. Load our schema file from file to create an empty database
	cmd := exec.Command("mysql", config.DBName)
	script, err := os.Open(scriptPath)
	if err != nil {
		t.Fatalf("error reading MySQL script: %s\n", err)
	}
	// Push the script in via standard input
	cmd.Stdin = script
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// 2. Run the schema file against the database
	err = cmd.Run()
	if err != nil {
		t.Fatalf("error creating testing database from schema:\n%s\n%s\n%s\n",
			err, stdout.String(), stderr.String())
	}
	// 3. Connect to the newly created database
	database := MustConnect(config)
	return database
}
