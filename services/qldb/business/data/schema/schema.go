// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"context"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
)

// Migrate attempts to bring the schema for qldb up to date with the migrations
// defined in this package.
func Migrate(qldb *qldbdriver.QLDBDriver) error {
	_, err := qldb.Execute(context.Background(), func(txn qldbdriver.Transaction) (interface{}, error) {
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