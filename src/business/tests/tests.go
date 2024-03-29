package tests

import (
	"bytes"
	"database/sql"
	"fmt"
	//nolint:gosec // Need to find another workaround
	_ "github.com/mattn/go-sqlite3"
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

	initService := dbInitService.NewService(db, logger)

	isFirst, err := initService.IsFirstRun()
	if err != nil {
		t.Fatalf("Cannot identify if can run migrations for tests: %s", err)
	}

	if isFirst {
		if err := initService.Migrate(); err != nil {
			t.Fatalf("Running migrations: %s", err)
		}
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		//nolint:gosec //Not handle this error because this is teardown
		db.Close()
		//nolint:gosec //Not handle this error because this is teardown
		os.RemoveAll(dbPath)
		//nolint:gosec //Not handle this error because this is teardown
		w.Close()
		var buf bytes.Buffer
		//nolint:gosec //Not handle this error because this is teardown
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
	logger, db, teardown := NewUnit(t)

	test := Test{
		DB:       db,
		Logger:   logger,
		t:        t,
		Teardown: teardown,
	}

	return &test
}

// CreateObit test helper for obit creation
func CreateObit(t *testing.T, test *Test) {
	const q = `
		INSERT INTO 
		    gateway_view(
				obit_did, 
			 	usn, 
		 		serial_number_hash, 
		 		manufacturer, 
			 	part_number, 
			 	alternate_ids, 
			 	owner_did,
		 		obd_did,
			 	status,
		 		metadata,
		 		structured_data,
		 		documents,
				modified_on,
			 	checksum
			) 
		    VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := test.DB.Prepare(q)

	if err != nil {
		t.Fatal(err)
	}

	_, err = stmt.Exec(
		"did:obada:d7cf869423d12f623f5611e48d6f6665bbc4a270b6e09da2f4c32bcb1b949ecd",
		"test",
		"cae6b797ae2627d96689fed03adc28311d5f2175253c3a0e375301e225ddf44d",
		"SONY",
		"PN123456S",
		"[]",
		`did:obada:owner:123456`,
		"",
		"FUNCTIONAL",
		"[]",
		"[]",
		"{}",
		1624387537,
		"2eb12c48ad2f073c49b95fcf2190cec40548c69fdc6f49135dee0753020f1624",
	)

	if err != nil {
		t.Fatal(err)
	}
}

// CreateOwnerObits test helper that creates many obits for single owner
func CreateOwnerObits(t *testing.T, test *Test) {

	const q = `
		INSERT INTO 
		    gateway_view(
				obit_did, 
			 	usn, 
		 		serial_number_hash, 
		 		manufacturer, 
			 	part_number, 
			 	alternate_ids, 
			 	owner_did,
		 		obd_did,
			 	status,
		 		metadata,
		 		structured_data,
		 		documents,
				modified_on,
			 	checksum
			) 
		    VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := test.DB.Prepare(q)

	for i := 0; i < 150; i++ {
		if err != nil {
			t.Fatal(err)
		}

		_, err = stmt.Exec(
			"did:obada:d7cf869423d12f623f5611e48d6f6665bbc4a270b6e09da2f4c32bcb1b949ec"+fmt.Sprintf("%d", i),
			"usn"+fmt.Sprintf("%d", i),
			"cae6b797ae2627d96689fed03adc28311d5f2175253c3a0e375301e225ddf44"+fmt.Sprintf("%d", i),
			"SONY",
			"PN123456S",
			"[]",
			`did:obada:owner:678910`,
			"",
			"FUNCTIONAL",
			"[]",
			"[]",
			"{}",
			1624387537,
			"2eb12c48ad2f073c49b95fcf2190cec40548c69fdc6f49135dee0753020f1624",
		)

		if err != nil {
			t.Fatal(err)
		}
	}
}
