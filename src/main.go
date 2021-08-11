package main

import (
	"fmt"
	"github.com/obada-foundation/node/cmd"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

// Opts with all cli commands and flags
type Opts struct {
	InitCommand cmd.InitCommand `command:"init"`
	RunCommand  cmd.RunCommand  `command:"run"`

	NodeURL    string `long:"url" env:"NODE_URL" required:"true" description:"url to OBADA API Node"`
}

var version = "undefined"

func main() {
	logger := log.New(os.Stdout, "OBADA-NODE :: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	fmt.Printf("OBADA-NODE %s\n", version)

	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		c := command.(cmd.CommonOptionsCommander)

		c.SetCommon(cmd.CommonOpts{
			Version: version,
			Logger:  logger,
			NodeURL: opts.NodeURL,
		})

		err := c.Execute(args)
		if err != nil {
			logger.Printf("main :: failed with %+v", err)
		}
		return err
	}

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
