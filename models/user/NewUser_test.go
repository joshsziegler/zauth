package user

import (
	"testing"

	dbLib "github.com/joshsziegler/zauth/pkg/db"
)

func TestNewUserDuplicate(t *testing.T) {
	// TODO: Create DB if doesn't exist, with '_test' postfix
	database := dbLib.MustConnect(dbLib.Config{
		Username: "metis",
		Password: "metis",
		Address:  "localhost",
		DBName:   "zauth_test"})
	Init(database)

	// NewUser(tx *sqlx.Tx, firstName string, lastName string, email string) (user User, err error)
	tx := dbLib.GetTxOrFailTesting(t, database)
	_, err := NewUser(tx, "first", "last", "josh.s.ziegler@gmail.com")
	if err != nil {
		t.Errorf("Creating a valid user failed: \n%+v", err)
	}

}
