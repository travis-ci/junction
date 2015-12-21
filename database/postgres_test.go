package database

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestPostgres(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	require.Nil(t, err)
	defer db.Close()

	// Clear the database before and after the tests
	db.Exec("DELETE FROM junction.workers")
	defer db.Exec("DELETE FROM junction.workers")

	pd := &Postgres{db: db}

	testDatabase(t, pd)
}
