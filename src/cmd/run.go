package cmd

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/qldbsession"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	_ "github.com/mattn/go-sqlite3"
	"github.com/obada-foundation/node/business/obit"
	searchService "github.com/obada-foundation/node/business/search"
	"github.com/obada-foundation/node/rest"
	"github.com/obada-foundation/sdkgo"
	"github.com/pkg/errors"
)

type RunCommand struct {
	API RestAPIGroup `group:"api" namespace:"api" env-namespace:"API"`
	AWS AWSGroup `group:"aws" namespace:"aws" env-namespace:"AWS"`
	QLDB QLDBGroup `group:"qldb" namespace:"qldb" env-namespace:"QLDB"`
	SSL SSLGroup `group:"ssl" namespace:"ssl" env-namespace:"SSL"`
	SQL SQL `group:"sql" namespace:"sql" env-namespace:"SQL"`

	CommonOpts
}

type SQL struct {
	SqlitePath string `default:"obada.db" description:"path to SQLite database"`
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
	ACMEEmail    string `long:"acme-email" env:"ACME_EMAIL" description:"admin email for certificate notifications"`
}

type AWSGroup struct {
	Region string `long:"region" description:"AWS region" default:"us-east-1"`
	Key string `long:"key" description:"AWS credential key" default:"us-east-1"`
	Secret string `long:"secret" description:"AWS credential secret" default:"us-east-1"`
}

type QLDBGroup struct {
	Database string `long:"database" description:"QLDB database name" default:"obada"`
}

type serverApp struct {
	*RunCommand

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

	if err = app.run(ctx); err != nil {
		return err
	}

	return nil
}

func (rc *RunCommand) newServerApp() (*serverApp, error) {
	// Initialize sqlite
	db, err := sql.Open("sqlite3", rc.SQL.SqlitePath)
	defer func() {
		rc.Logger.Println("main: SQLite closing database connection")
		if er := db.Close(); er != nil {
			rc.Logger.Printf("main: Cannot close SQLite database: %s", err)
		}
	}()

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

	defer func() {
		rc.Logger.Printf("main: QLDB Stopping database connection : %s", rc.QLDB.Database)
		qldb.Shutdown(context.Background())
	}()


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
		Version: rc.Version,
		Logger: rc.Logger,

		SearchService: search,
		ObitService: obitService,
	}

	return &serverApp{
		RunCommand: rc,

		restSrv:    srv,
		terminated: make(chan struct{}),
	}, nil
}

func (a *serverApp) run(ctx context.Context) error {
	go func() {
		// shutdown on context cancellation
		<-ctx.Done()
		a.Logger.Print("server :: shutdown initiated")
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

