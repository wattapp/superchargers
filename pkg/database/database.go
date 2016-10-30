package database

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/mgutz/dat.v1"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
)

var db runner.Connection

func Conn() runner.Connection {
	return db
}

// Connect establishes a connection with the configured DATABASE_URL. You can
// keep a copy of the returned connection or fetch it using database.Conn()
func Connect() (runner.Connection, error) {
	driver, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	// ensures the database can be pinged with an exponential backoff (15 min)
	runner.MustPing(driver)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	d := os.Getenv("DATABASE_SLOW_QUERY_DURATION")
	dur, err := time.ParseDuration(d)
	if err != nil {
		return nil, err
	}
	runner.LogQueriesThreshold = dur

	db = runner.NewDB(driver, "postgres")
	return db, nil
}
