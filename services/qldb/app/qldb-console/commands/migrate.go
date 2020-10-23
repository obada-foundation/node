package commands

import (
	"context"
	"fmt"

	"github.com/obada-protocol/server-gateway/services/qldb/business/data/schema"
	"github.com/obada-protocol/server-gateway/services/qldb/foundation/database"
	"github.com/pkg/errors"
)

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

// Migrate creates the schema in the database.
func Migrate(cfg database.Config) error {
	qldb, err := database.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect database")
	}
	defer qldb.Close(context.Background())

	if err := schema.Migrate(qldb); err != nil {
		return errors.Wrap(err, "migrate database")
	}

	fmt.Println("migrations complete")
	return nil
}