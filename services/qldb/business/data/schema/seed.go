package schema

import (
	"context"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
)

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(qldb *qldbdriver.QLDBDriver) error {
	_, err := qldb.Execute(context.Background(), func(txn qldbdriver.Transaction) (interface{}, error) {
		_, err := txn.Execute(seeds)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

// seeds is a string constant containing all of the queries needed to get the
// db seeded to a useful state for development.
const seeds = `
INSERT INTO Obits {};
`

// DeleteAll runs the set of Drop-table queries against qldb. The queries are ran in a
// transaction and rolled back if any fail.
func DeleteAll(qldb *qldbdriver.QLDBDriver) error {
	_, err := qldb.Execute(context.Background(), func(txn qldbdriver.Transaction) (interface{}, error) {
		_, err := txn.Execute(deleteAll)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

// deleteAll is used to clean the database between tests.
const deleteAll = `
DELETE FROM Obits;
`
