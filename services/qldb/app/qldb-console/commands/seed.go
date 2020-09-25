package commands

import (
	"context"
	"fmt"

	"github.com/obada-protocol/server-gateway/services/qldb/business/data/schema"
	"github.com/obada-protocol/server-gateway/services/qldb/foundation/database"
	"github.com/pkg/errors"
)

// Seed loads test data into the database.
func Seed(cfg database.Config) error {
	qldb, err := database.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect database")
	}
	defer qldb.Close(context.Background())

	if err := schema.Seed(qldb); err != nil {
		return errors.Wrap(err, "seed database")
	}

	fmt.Println("seed data complete")
	return nil
}