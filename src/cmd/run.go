package cmd

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/qldbsession"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	//nolint:gosec // Need to find another workaround
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/obada-foundation/node/business/obit"
	searchService "github.com/obada-foundation/node/business/search"
	"github.com/obada-foundation/node/rest"
	"github.com/obada-foundation/sdkgo"
)

// RunCommand runs OBADA Node
type RunCommand struct {
	API  RestAPIGroup `group:"api" namespace:"api" env-namespace:"API"`
	AWS  AWSGroup     `group:"aws" namespace:"aws" env-namespace:"AWS"`
	QLDB QLDBGroup    `group:"qldb" namespace:"qldb" env-namespace:"QLDB"`
	SSL  SSLGroup     `group:"ssl" namespace:"ssl" env-namespace:"SSL"`
	SQL  SQL          `group:"sql" namespace:"sql" env-namespace:"SQL"`

	CommonOpts
}

// SQL defines options for SQLite
type SQL struct {
	SqlitePath string `long:"sql-path " default:"obada.db" description:"path to SQLite database"`
}

// RestAPIGroup defines options group for REST API params
type RestAPIGroup struct {
	Port    int    `long:"port" default:"3000" description:"port"`
	Address string `long:"address" default:"127.0.0.1" description:"listening address"`
}

// SSLGroup defines options group for server ssl params
type SSLGroup struct {
	Type         string `long:"type" env:"TYPE" description:"ssl (auto) support" choice:"none" choice:"static" choice:"auto" default:"none"` //nolint
	Port         int    `long:"port" env:"PORT" description:"port number for https server" default:"3443"`
	Cert         string `long:"cert" env:"CERT" description:"path to cert.pem file"`
	Key          string `long:"key" env:"KEY" description:"path to key.pem file"`
	ACMELocation string `long:"acme-location" env:"ACME_LOCATION" description:"dir where certificates will be stored by autocert manager" default:"./var/acme"`
	ACMEEmail    string `long:"acme-email" env:"ACME_EMAIL" description:"admin email for certificate notifications" default:"techops@obada.io"`
}

// AWSGroup defines options for AWS connection
type AWSGroup struct {
	Region string `long:"region" description:"AWS region" default:"us-east-1"`
	Key    string `long:"key" description:"AWS credential key" default:"us-east-1"`
	Secret string `long:"secret" description:"AWS credential secret" default:"us-east-1"`
}

// QLDBGroup defines options for QLDB database
type QLDBGroup struct {
	Database string `long:"database" description:"QLDB database name" default:"obada"`
}

type serverApp struct {
	*RunCommand

	db         *sql.DB
	qldb       *qldbdriver.QLDBDriver
	restSrv    *rest.Rest
	terminated chan struct{}
}

// Execute is the entry point for "run" command, called by flag parser
func (rc *RunCommand) Execute(_ []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		rc.Logger.Printf("interrupt signal")
		cancel()
	}()

	app, err := rc.newServerApp()
	if err != nil {
		return err
	}

	if err := app.run(ctx); err != nil {
		return err
	}

	return nil
}

func (rc *RunCommand) newServerApp() (*serverApp, error) {
	rc.Logger.Printf("BBBBB %s", rc.NodeURL)

	if !strings.HasPrefix(rc.NodeURL, "http://") && !strings.HasPrefix(rc.NodeURL, "https://") {
		return nil, errors.Errorf("invalid Node API url %s", rc.NodeURL)
	}
	rc.Logger.Printf("Node url=%s", rc.NodeURL)

	sslConfig, err := rc.makeSSLConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make config of ssl server params")
	}

	// Initialize sqlite
	db, err := sql.Open("sqlite3", rc.SQL.SqlitePath)
	if err != nil {
		return nil, errors.Wrap(err, "initializing sqlite database")
	}

	// Create AWS session
	awsConfig := aws.NewConfig().WithRegion(rc.AWS.Region)
	awsConfig.Credentials = credentials.NewStaticCredentials(rc.AWS.Key, rc.AWS.Secret, "")
	awsSession := session.Must(session.NewSession(awsConfig))

	// =========================================================================
	// Start QLDB
	qldbSession := qldbsession.New(awsSession)

	qldb, err := qldbdriver.New(
		rc.QLDB.Database,
		qldbSession,
		func(options *qldbdriver.DriverOptions) {
			options.LoggerVerbosity = qldbdriver.LogInfo
		})

	if err != nil {
		return nil, errors.Wrap(err, "trying to configure QLDB")
	}

	// Initialize OBADA SDK
	sdk, err := sdkgo.NewSdk(rc.Logger, true)
	if err != nil {
		return nil, errors.Wrap(err, "initializing OBADA SDK")
	}

	// Initialize SearchService
	search := searchService.NewService(rc.Logger, db)

	// Initialize ObitService
	obitService := obit.NewObitService(sdk, rc.Logger, db, qldb, nil)

	srv := &rest.Rest{
		Version:   rc.Version,
		Logger:    rc.Logger,
		NodeURL:   rc.NodeURL,
		SSLConfig: sslConfig,

		SearchService: search,
		ObitService:   obitService,
	}

	return &serverApp{
		RunCommand: rc,

		db:         db,
		qldb:       qldb,
		restSrv:    srv,
		terminated: make(chan struct{}),
	}, nil
}

func (rc *RunCommand) makeSSLConfig() (config rest.SSLConfig, err error) {
	switch rc.SSL.Type {
	case "none":
		config.SSLMode = rest.None
	case "static":
		if rc.SSL.Cert == "" {
			return config, errors.New("path to cert.pem is required")
		}
		if rc.SSL.Key == "" {
			return config, errors.New("path to key.pem is required")
		}
		config.SSLMode = rest.Static
		config.Port = rc.SSL.Port
		config.Cert = rc.SSL.Cert
		config.Key = rc.SSL.Key
	case "auto":
		config.SSLMode = rest.Auto
		config.Port = rc.SSL.Port
		config.ACMELocation = rc.SSL.ACMELocation
		config.ACMEEmail = rc.SSL.ACMEEmail
	}
	return config, err
}

func (a *serverApp) run(ctx context.Context) error {
	go func() {
		// shutdown on context cancellation
		<-ctx.Done()
		a.Logger.Print("server :: shutdown initiated")

		a.Logger.Println("main: SQLite closing database connection")
		if er := a.db.Close(); er != nil {
			a.Logger.Printf("main: Cannot close SQLite database: %s", er)
		}

		a.Logger.Printf("main: QLDB Stopping database connection : %s", a.QLDB.Database)
		a.qldb.Shutdown(ctx)

		a.restSrv.Shutdown()
	}()

	a.restSrv.Run(a.API.Address, a.API.Port)

	close(a.terminated)

	return nil
}

// Wait for application completion (termination)
func (a *serverApp) Wait() {
	<-a.terminated
}
