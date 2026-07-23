// Command seed fills the database with development data. It uses the same
// code paths as the application (user.NewUser, user.SetUserPassword, etc.) so
// passwords are hashed and all validation runs.
//
// Connection settings come from the ZAUTH_DB_* environment variables (see
// db.ConfigFromEnv). Run it against a fresh database; it does not try to be
// idempotent. To start over: docker compose down -v && docker compose up -d
package main

import (
	"github.com/jmoiron/sqlx"

	"github.com/joshsziegler/zauth/pkg/db"
	"github.com/joshsziegler/zauth/pkg/user"
	"github.com/joshsziegler/zgo/pkg/log"
)

// devPassword is the password for every seeded user. Dev use only.
const devPassword = "correct-horse-battery-staple"

type seedUser struct {
	First  string
	Last   string
	Email  string
	Groups []string
}

func main() {
	database := db.MustConnect(db.ConfigFromEnv())

	groups := map[string]string{
		"admin": "Administrators (members get admin rights in zauth)",
		"staff": "Example non-admin group",
	}
	users := []seedUser{
		{"Alice", "Zephyr", "alice.zephyr@example.com", []string{"admin", "staff"}},
		{"Bob", "Yonder", "bob.yonder@example.com", []string{"staff"}},
		{"Carol", "Xylo", "carol.xylo@example.com", nil},
	}

	err := seed(database, groups, users)
	if err != nil {
		log.Fatalf("seeding failed (already seeded?): %+v", err)
	}
	log.Infof("seeded %d groups and %d users (password: %s)",
		len(groups), len(users), devPassword)
}

func seed(database *sqlx.DB, groups map[string]string, users []seedUser) error {
	tx, err := database.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for name, desc := range groups {
		err = user.AddGroup(tx, name, desc)
		if err != nil {
			return err
		}
	}
	for _, s := range users {
		u, err := user.NewUser(tx, s.First, s.Last, s.Email)
		if err != nil {
			return err
		}
		err = user.SetUserPassword(tx, u.Username, devPassword)
		if err != nil {
			return err
		}
		for _, g := range s.Groups {
			err = user.AddUserToGroup(tx, u.Username, g)
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}
