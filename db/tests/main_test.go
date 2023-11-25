package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:testing@localhost:5432/spectrumshelf_db?sslmode=disable"
)

var testQueries *db.Queries

func TestMain(m *testing.M) {
	// This function is to perform the main test.
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("There was an error connecting to database: %v", err)
	}
	testQueries = db.New(conn)
	os.Exit(m.Run())
}
