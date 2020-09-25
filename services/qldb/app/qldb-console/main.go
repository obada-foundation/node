// This program performs console tasks for the qldb service.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/obada-protocol/server-gateway/services/qldb/app/qldb-console/commands"
	"github.com/obada-protocol/server-gateway/services/qldb/foundation/database"
	"github.com/pkg/errors"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	log := log.New(os.Stdout, "CONSOLE : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		if errors.Cause(err) != commands.ErrHelp {
			log.Printf("error: %s", err)
		}
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	// =========================================================================
	// Configuration
	var cfg struct {
		conf.Version
		Args conf.Args
		QLDB struct {
			Region    string `conf:default:us-east-1`
			Database  string `conf:default:obada`
			Key       string `conf`
			Secret    string `conf`
		}
	}
	cfg.Version.SVN = build
	cfg.Version.Desc = "(c) OBADA 2020"

	const prefix = "QLDB"
	if err := conf.Parse(os.Args[1:], prefix, &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config :\n%v\n", out)

	// =========================================================================
	// Commands

	qldbConfig := database.Config{
		Region:     cfg.QLDB.Region,
		Database:   cfg.QLDB.Database,
		Key:        cfg.QLDB.Key,
		Secret:     cfg.QLDB.Secret,
	}

	switch cfg.Args.Num(0) {
	case "migrate":
		if err := commands.Migrate(qldbConfig); err != nil {
			return errors.Wrap(err, "migrating database")
		}

	case "seed":
		if err := commands.Seed(qldbConfig); err != nil {
			return errors.Wrap(err, "seeding database")
		}

	default:
		fmt.Println("migrate: create the schema in the database")
		fmt.Println("seed: add data to the database")
		fmt.Println("provide a command to get more help.")
		return commands.ErrHelp
	}

	return nil
}