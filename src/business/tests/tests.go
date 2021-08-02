package tests

import (
	"bytes"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	dbInitService "github.com/obada-foundation/node/business/database"
	"io"
	"log"
	"os"
	"testing"
)

// Test owns state for running and shutting down tests.
type Test struct {
	DB       *sql.DB
	Logger   *log.Logger
	Teardown func()

	t *testing.T
}

var dbPath = "/tmp/nodetest"

// NewUnit creates a test database. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewUnit(t *testing.T) (*log.Logger, *sql.DB, func()) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	logger := log.New(os.Stdout, "", 0)

	initService := dbInitService.NewService(db, nil, logger)

	initService.Migrate()

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
		os.RemoveAll(dbPath)
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old
		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}

	return logger, db, teardown
}

// NewIntegration creates a database, seeds it, constructs an authenticator.
func NewIntegration(t *testing.T) *Test {
	log, db, teardown := NewUnit(t)

	test := Test{
		DB:       db,
		Logger:   log,
		t:        t,
		Teardown: teardown,
	}

	return &test
}