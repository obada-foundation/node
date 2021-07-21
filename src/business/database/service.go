package database

import (
	"context"
	"database/sql"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	"log"
)

type Service struct {
	db     *sql.DB
	qldb   *qldbdriver.QLDBDriver
	logger *log.Logger
}

func NewService(db *sql.DB, qldb *qldbdriver.QLDBDriver, logger *log.Logger) Service {
	return Service{
		db:     db,
		qldb:   qldb,
		logger: logger,
	}
}

// IsFirstRun checks if service runs first time on machine
func (s Service) IsFirstRun() (bool, error) {
	var cnt int

	q := `
		SELECT 
		    count(*) as cnt 
		FROM 
		    sqlite_master 
		WHERE 
		    type='table' AND name='gateway_view';
	`

	row := s.db.QueryRow(q)

	if err := row.Scan(&cnt); err != nil {
		return false, err
	}

	return cnt == 0, nil
}

// Migrate attempts to bring the schema for qldb up to date with the migrations
// defined in this package.
func (s Service) qldbMigrate() error {
	s.logger.Println("Running QLDB migrations")

	_, err := s.qldb.Execute(context.Background(), func(txn qldbdriver.Transaction) (interface{}, error) {
		_, err := txn.Execute("CREATE TABLE Obits")
		if err != nil {
			return nil, err
		}

		// When working with QLDB, it's recommended to create an index on fields we're filtering on.
		// This reduces the chance of OCC conflicts exceptions with large datasets.
		_, err = txn.Execute("CREATE INDEX ON Obits (ObitDID)")
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

// Migrate attempts to bring the schema for sqlite up to date with the migrations
// defined in this package.
func (s Service) sqliteMigrate() error {
	queries := []string{
		`create table gateway_view (
			obit_did           varchar(255) not null,
			usn                varchar(255) not null,
			serial_number_hash varchar(255) not null,
			manufacturer       varchar(255) not null,
			part_number        varchar(255) not null,
			alternate_ids	   json         null,
			owner_did          varchar(255) null,
			obd_did            varchar(255) null,
			status        	   varchar(10)  null,
			metadata           json         null,
			structured_data    json         null,
			documents          json         null,
			modified_on        int          null,
			root_hash          varchar(255) not null,
			constraint gateway_view_obit_did_usn_serial_number_hash_unique unique (obit_did, usn, serial_number_hash)
		)`,
		`create table config (
			last_obit varchar(255) not null
		)`,
	}

	for _, q := range queries {
		stmt, err := s.db.Prepare(q)

		if err != nil {
			return err
		}

		if _, err := stmt.Exec(); err != nil {
			return err
		}
	}

	return nil
}

// Migrate attempts to bring the schema for qldb up to date with the migrations
// defined in this package.
func (s Service) Migrate() error {

	if err := s.sqliteMigrate(); err != nil {
		return err
	}

	//if err := s.qldbMigrate(); err != nil {
	//	return err
	//}

	return nil
}