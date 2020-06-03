package httpserver

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joshsziegler/zauth/pkg/db"
)

type App struct {
	ProgramName string
	// DB is our database connection (handles connection pooling, and is goroutine-safe)
	DB *sqlx.DB
}

func NewApp(programName string, dbConfig db.Config) (*App, error) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)
	app := App{
		ProgramName: programName,
	}
	// app.DB = db.MustConnect(app.Log, dbConfig)
	return &app, nil
}
