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

	"database/sql"
	"github.com/ardanlabs/conf"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/qldbsession"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	_ "github.com/mattn/go-sqlite3"
	"github.com/obada-foundation/node/app/node/handlers"
	dbInitService "github.com/obada-foundation/node/business/database"
	obitService "github.com/obada-foundation/node/business/obit"
	searchService "github.com/obada-foundation/node/business/search"
	pubsub "github.com/obada-foundation/node/business/pubsub/aws"
	"github.com/obada-foundation/sdkgo"
	"github.com/pkg/errors"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	logger := log.New(os.Stdout, "OBADA-NODE :", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(logger); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run(logger *log.Logger) error {

	var cfg struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		SQL struct {
			SqlitePath string `conf:"default:obada.db"`
		}
		Zipkin struct {
			ReporterURI string  `conf:"default:http://zipkin:9411/api/v2/spans"`
			ServiceName string  `conf:"default:obada-node"`
			Probability float64 `conf:"default:0.05"`
		}
		AWS struct {
			Region string `conf:"default:us-east-1"`
			Key    string `conf:"noprint"`
			Secret string `conf:"noprint"`
		}
		QLDB struct {
			Database string `conf:"default:obada"`
		}
		PUBSUB struct {
			Timeout  time.Duration `conf:"default:5s"`
			QueueURL string        `conf:"default:https://sqs.us-east-1.amazonaws.com/271164744603/obada-tradeloop.fifo"`
			TopicArn string        `conf:"default:arn:aws:sns:us-east-1:271164744603:obada.fifo"`
		}
	}
	cfg.Version.SVN = build
	cfg.Version.Desc = "(c) OBADA Foundation 2021"

	if err := conf.Parse(os.Args[1:], "", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, er := conf.Usage("", &cfg)
			if er != nil {
				return errors.Wrap(er, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, er := conf.VersionString("", &cfg)
			if er != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	expvar.NewString("build").Set(build)
	logger.Printf("main : Started : Application initializing : version %q", build)
	defer logger.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	logger.Printf("main: Config :\n%v\n", out)

	logger.Println("main: Initializing debugging support")

	go func() {
		logger.Printf("main: Debug Listening %s", cfg.Web.DebugHost)
		if er := http.ListenAndServe(cfg.Web.DebugHost, http.DefaultServeMux); er != nil {
			logger.Printf("main: Debug Listener closed : %v", er)
		}
	}()

	// Create AWS session
	awsConfig := aws.NewConfig().WithRegion(cfg.AWS.Region)
	awsConfig.Credentials = credentials.NewStaticCredentials(cfg.AWS.Key, cfg.AWS.Secret, "")
	awsSession := session.Must(session.NewSession(awsConfig))

	// =========================================================================
	// Start QLDB
	qldbSession := qldbsession.New(awsSession)

	qldb, err := qldbdriver.New(
		cfg.QLDB.Database,
		qldbSession,
		func(options *qldbdriver.DriverOptions) {
			options.LoggerVerbosity = qldbdriver.LogInfo
		})

	if err != nil {
		return errors.Wrap(err, "trying to configure QLDB")
	}

	defer func() {
		logger.Printf("main: QLDB Stopping database connection : %s", cfg.QLDB.Database)
		qldb.Shutdown(context.Background())
	}()

	if _, err = os.Stat(cfg.SQL.SqlitePath); os.IsNotExist(err) {
		file, er := os.Create(cfg.SQL.SqlitePath)

		if er != nil {
			return errors.Wrap(er, "Problem with creating sqlite db file")
		}

		file.Close()
	}

	// Initialize sqlite
	db, err := sql.Open("sqlite3", cfg.SQL.SqlitePath)
	defer func() {
		logger.Println("main: SQLite closing database connection")
		db.Close()
	}()

	if err != nil {
		return errors.Wrap(err, "initializing sqlite database")
	}

	initService := dbInitService.NewService(db, qldb, logger)

	isFirstRun, err := initService.IsFirstRun()

	if err != nil {
		return err
	}

	if isFirstRun {
		if er := initService.Migrate(); er != nil {
			return errors.Wrap(er, "Problem with running migrations")
		}
	}

	// Initialize OBADA SDK
	sdk, err := sdkgo.NewSdk(logger, true)

	if err != nil {
		return errors.Wrap(err, "initializing OBADA SDK")
	}

	ps := pubsub.NewClient(awsSession, cfg.PUBSUB.Timeout, cfg.PUBSUB.QueueURL, cfg.PUBSUB.TopicArn)

	// Initialize SearchService
	search := searchService.NewService(logger, db)

	// Initialize ObitService
	obit := obitService.NewObitService(sdk, logger, db, qldb, ps)

	// Start database sync
	go func() {
		for {
			if err := obit.Sync(context.Background()); err != nil {
				logger.Println(err)
			}
		}
	}()

	logger.Println("main: Initializing API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		Addr: cfg.Web.APIHost,
		Handler: handlers.API(handlers.APIConfig{
			Shutdown:      shutdown,
			Logger:        logger,
			ObitService:   obit,
			SearchService: search,
		}),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		logger.Printf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		logger.Printf("main: %v : Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			errClose := api.Close()

			if errClose != nil {
				return errors.Wrap(errClose, "could not stop server gracefully")
			}

			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
