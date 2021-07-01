package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/obada-foundation/node/app/node/handlers"
	obitService "github.com/obada-foundation/node/business/obit"
	sdk "github.com/obada-foundation/sdkgo"
	"github.com/pkg/errors"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	log := log.New(os.Stdout, "OBADA-NODE :", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	var cfg struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		Zipkin struct {
			ReporterURI string  `conf:"default:http://zipkin:9411/api/v2/spans"`
			ServiceName string  `conf:"default:obada-node"`
			Probability float64 `conf:"default:0.05"`
		}
		QLDB struct {
			Region   string `conf:"default:us-east-1,env:REGION"`
			Database string `conf:"default:obada,env:DATABASE"`
			Key      string `conf:"env:KEY,noprint"`
			Secret   string `conf:"env:SECRET,noprint"`
		}
	}
	cfg.Version.SVN = build
	cfg.Version.Desc = "(c) OBADA Foundation 2021"

	if err := conf.Parse(os.Args[1:], "OBADA-NODE", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("OBADA-NODE", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("OBADA-NODE", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	expvar.NewString("build").Set(build)
	log.Printf("main : Started : Application initializing : version %q", build)
	defer log.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config :\n%v\n", out)

	log.Println("main: Initializing debugging support")

	go func() {
		log.Printf("main: Debug Listening %s", cfg.Web.DebugHost)
		if err := http.ListenAndServe(cfg.Web.DebugHost, http.DefaultServeMux); err != nil {
			log.Printf("main: Debug Listener closed : %v", err)
		}
	}()

	// Initialize OBADA SDK
	sdk, err := sdk.NewSdk(log, true)

	if err != nil {
		return errors.Wrap(err, "initializing OBADA SDK")
	}

	// Initialize ObitService
	obitService := obitService.NewObitService(sdk)

	log.Println("main: Initializing API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      handlers.API(build, shutdown, log, obitService),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main: %v : Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
