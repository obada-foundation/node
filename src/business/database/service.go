package database

import (
	"context"
	"database/sql"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	"log"
)

// Service temp struct to handle dependencies
type Service struct {
	db     *sql.DB
	qldb   *qldbdriver.QLDBDriver
	logger *log.Logger
}

// NewService creates database service
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
//nolint:unused // Need this code for a future use
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
			id                 integer      PRIMARY KEY AUTOINCREMENT not null,
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
			modified_on        integer      null,
			checksum           varchar(64) not null,
			constraint gateway_view_obit_did_usn_serial_number_hash_unique unique (obit_did, usn, serial_number_hash)
		)`,
		`create table config (
			last_obit varchar(255) not null
		)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS 
			gateway_view_fts 
		USING 
			fts5(
				content=gateway_view, 
				content_rowid=id,
				tokenize="unicode61 tokenchars ';'",
				obit_did, 
				usn, 
				serial_number_hash, 
				manufacturer, 
				part_number,
				owner_did
			)`,
		`-- Triggers to keep the FTS index up to date.
		CREATE TRIGGER gateway_view_ai AFTER INSERT ON gateway_view BEGIN
		  INSERT INTO
			gateway_view_fts(rowid, obit_did, usn, serial_number_hash, manufacturer, part_number, owner_did)
			VALUES (
				new.id, 
				REPLACE(new.obit_did, ":", ""), 
				new.usn,
				new.serial_number_hash, 
				new.manufacturer, 
				new.part_number, 
				REPLACE(new.owner_did, ":", "")
			);
		END;`,
		`CREATE TRIGGER gateway_view_ad AFTER DELETE ON gateway_view BEGIN
		  INSERT INTO 
			gateway_view_fts(
				gateway_view_fts, rowid, obit_did, usn, serial_number_hash, manufacturer, part_number, owner_did
			) 
			VALUES(
				'delete', 
				old.id, 
				REPLACE(old.obit_did, ":", ""), 
				old.usn, 
				old.serial_number_hash, 
				old.manufacturer, 
				old.part_number,
				REPLACE(old.owner_did, ":", "")
			);
		END;`,
		`CREATE TRIGGER gateway_view_au AFTER UPDATE ON gateway_view BEGIN
		  INSERT INTO 
			gateway_view_fts(
				gateway_view_fts, rowid, obit_did, usn, serial_number_hash, manufacturer, part_number, owner_did
			) 
			VALUES(
				'delete', 
				old.id, 
				REPLACE(old.obit_did, ":", ""), 
				old.usn, 
				old.serial_number_hash, 
				old.manufacturer, 
				old.part_number,
				REPLACE(old.owner_did, ":", "")
			);
		  INSERT INTO
			gateway_view_fts(
				rowid, obit_did, usn, serial_number_hash, manufacturer, part_number, owner_did
			)
			VALUES (
				new.id, 
				REPLACE(new.obit_did, ":", ""), 
				new.usn, 
				new.serial_number_hash, 
				new.manufacturer, 
				new.part_number,
				REPLACE(new.owner_did, ":", "")
			);
		END`,
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

// Migrate attempts to bring the schema for QLDB up to date with the migrations
// defined in this package.
func (s Service) Migrate() error {

	if err := s.sqliteMigrate(); err != nil {
		return err
	}

	return nil
}
