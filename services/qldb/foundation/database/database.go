// Package database provides support for access the database.
package database

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"strings"

	"github.com/amzn/ion-go/ion"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/qldbsession"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	"go.opentelemetry.io/otel/api/trace"
)

// Config is the required properties to use the database.
type Config struct {
	Database string
	Region   string
	Key      string
	Secret   string
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (*qldbdriver.QLDBDriver, error) {
	awsConfig := aws.NewConfig()
	awsConfig.Credentials = credentials.NewStaticCredentials(cfg.Key, cfg.Secret, "")

	awsSession := session.Must(session.NewSession(awsConfig.WithRegion(cfg.Region)))
	qldbSession := qldbsession.New(awsSession)

	qldb := qldbdriver.New(
		cfg.Database,
		qldbSession,
		func(options *qldbdriver.DriverOptions) {
			options.LoggerVerbosity = qldbdriver.LogInfo
		})

	return qldb, nil
}

// StatusCheck returns nil if it can successfully talk to the qldb. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, qldb *qldbdriver.QLDBDriver) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "foundation.database.statuscheck")
	defer span.End()

	var tmp bool

	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// qldb but the database has since gone away. Running this query forces a
	// round trip to the database.
	_, err := qldb.Execute(context.Background(), func(txn qldbdriver.Transaction) (interface{}, error) {
		result, err := txn.Execute(`SELECT true`)
		if err != nil {
			return nil, err
		}

		// Assume the result is not empty
		ionBinary, err := result.Next(txn)
		if err != nil {
			return nil, err
		}

		err = ion.Unmarshal(ionBinary, &tmp)
		if err != nil {
			return nil, err
		}


		return nil, nil
	})

	return err
}

// Log provides a pretty print version of the query and parameters.
func Log(query string, args ...interface{}) string {
	for i, arg := range args {
		n := fmt.Sprintf("$%d", i+1)

		var a string
		switch v := arg.(type) {
		case string:
			a = fmt.Sprintf("%q", v)
		case []byte:
			a = string(v)
		case []string:
			a = strings.Join(v, ",")
		default:
			a = fmt.Sprintf("%v", v)
		}

		query = strings.Replace(query, n, a, 1)
	}

	return query
}