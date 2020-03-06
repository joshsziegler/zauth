package user

import (
	"testing"

	"github.com/joshsziegler/zauth/pkg/db"
)

func TestNewUserDuplicate(t *testing.T) {
	database := db.SetupTestingDatabase(t, db.Config{
		Username: "joshz",
		Password: "Manikin06!",
		Address:  "localhost",
		DBName:   "zauth"},
		"../../db-schema-v2.0.sql")

	// NewUser(tx *sqlx.Tx, firstName string, lastName string, email string) (user User, err error)
	tx := db.GetTxOrFailTesting(t, database)
	_, err := NewUser(tx, "first", "last", "first.last@email.com")
	if err != nil {
		t.Errorf("Creating a valid user failed: \n%+v", err)
	}
	tx.Commit()
	// Now create the same user again...
	tx = db.GetTxOrFailTesting(t, database)
	_, err = NewUser(tx, "first", "last", "first.last@email.com")
	if err == nil {
		t.Errorf("Creating a duplicate user didn't return an error: \n%+v", err)
	}
	tx.Commit()

	// Create a non-duplicate, just to be sure
	tx = db.GetTxOrFailTesting(t, database)
	_, err = NewUser(tx, "John", "Doe", "doe@email.com")
	if err != nil {
		t.Errorf("Creating a valid user failed: \n%+v", err)
	}
	tx.Commit()
}

func TestDuplicateEmail(t *testing.T) {
	database := db.SetupTestingDatabase(t, db.Config{
		Username: "joshz",
		Password: "Manikin06!",
		Address:  "localhost",
		DBName:   "zauth"},
		"../../db-schema-v2.0.sql")

	// Create a non-duplicate, just to be sure
	tx := db.GetTxOrFailTesting(t, database)
	_, err := NewUser(tx, "John", "Doe", "doe@email.com")
	if err != nil {
		t.Errorf("Creating a valid user failed: \n%+v", err)
	}
	tx.Commit()

	// Emails do NOT need to be unique, only usernames...
	tx = db.GetTxOrFailTesting(t, database)
	_, err = NewUser(tx, "Jane", "Doe", "doe@email.com")
	if err != nil {
		t.Errorf("Creating a valid user failed: \n%+v", err)
	}
	tx.Commit()
}
