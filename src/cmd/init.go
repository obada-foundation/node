package cmd

import (
	"database/sql"
	dbInitService "github.com/obada-foundation/node/business/database"
	"github.com/pkg/errors"
	"os"
)

// InitCommand initialize OBADA Node database files and configurations
type InitCommand struct {
	SQL SQL `group:"sql" namespace:"sql" env-namespace:"SQL"`

	CommonOpts
}

// Execute is the entry point for "init" command, called by flag parser
func (ic *InitCommand) Execute(_ []string) error {
	ic.Logger.Println("Creating required files")

	if _, err := os.Stat(ic.SQL.SqlitePath); os.IsNotExist(err) {
		file, er := os.Create(ic.SQL.SqlitePath)

		if er != nil {
			return errors.Wrap(er, "Problem with creating sqlite db file")
		}

		if er := file.Close(); er != nil {
			return errors.Wrap(er, "Problem with closing SQLite file")
		}
	}

	// Initialize sqlite
	db, err := sql.Open("sqlite3", ic.SQL.SqlitePath)
	defer func() {
		ic.Logger.Println("main: SQLite closing database connection")
		if er := db.Close(); er != nil {
			ic.Logger.Printf("main: Cannot close SQLite database: %s", err)
		}
	}()
	if err != nil {
		return errors.Wrap(err, "initializing sqlite database")
	}

	ic.Logger.Println("Running database migration")
	initService := dbInitService.NewService(db, ic.Logger)
	if er := initService.Migrate(); er != nil {
		return errors.Wrap(er, "Problem with running migrations")
	}

	ic.Logger.Println("Node initialization completed")

	return nil
}
