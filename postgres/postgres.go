package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var (
	postgresSession *sql.DB
)

// Initialize the postgres session.
// URLformat: postgres://user:password@host:port/database
func InitPostgres(url string) {
	connStr := url

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	postgresSession = db
}

// Get the postgres session
func GetPostgresSession() *sql.DB {
	return postgresSession
}
