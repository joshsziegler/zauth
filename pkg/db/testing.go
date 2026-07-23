//go:build !test

// The line above will exclude this file from builds, except during testing.
// This file contains methods useful for doing *real* database tests, rather
// than using mocks.
//

package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joshsziegler/zgo/pkg/log"
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

// SetupTestingDatabase loads the given SQL scripts (in order) into the
// database and returns a connection to it. Connection settings come from the
// ZAUTH_DB_* environment variables (see ConfigFromEnv), so tests run against
// a local or containerized MariaDB/MySQL server without code changes.
func SetupTestingDatabase(t *testing.T, scriptPaths ...string) *sqlx.DB {
	config := ConfigFromEnv()
	// Use a _test suffix so tests never touch the development data
	config.DBName += "_test"
	// Use a dedicated connection with multiStatements enabled so each schema
	// file can run in one Exec. We close it afterward rather than returning
	// it, since multiStatements increases SQL-injection risk. Connect without
	// a database name since the test database may not exist yet.
	dsn := getDSN(config.Username, config.Password,
		fmt.Sprintf("tcp(%s)", config.Address), "") +
		"&multiStatements=true"
	schemaConn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		t.Fatalf("error connecting to testing database: %s\n", err)
	}
	defer schemaConn.Close()
	// One connection only, so the USE below applies to every script we run
	schemaConn.SetMaxOpenConns(1)
	// Recreate the test database so every run starts from a clean slate
	_, err = schemaConn.Exec(fmt.Sprintf(
		"DROP DATABASE IF EXISTS `%s`; CREATE DATABASE `%s`; USE `%s`;",
		config.DBName, config.DBName, config.DBName))
	if err != nil {
		t.Fatalf("error recreating test database %s: %s\n", config.DBName, err)
	}
	for _, p := range scriptPaths {
		script, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("error reading SQL script %s: %s\n", p, err)
		}
		_, err = schemaConn.Exec(string(script))
		if err != nil {
			t.Fatalf("error running SQL script %s: %s\n", p, err)
		}
	}
	return MustConnect(config)
}
