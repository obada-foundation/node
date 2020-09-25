// Package obit contains obit related CRUD functionality.
package obit

import (
	"context"
	"github.com/amzn/ion-go/ion"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/trace"
	"log"
	"time"
)

var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// Obit manages the set of API's for obit access.
type Obit struct {
	qldb *qldbdriver.QLDBDriver
	log *log.Logger
}

// New constructs a Obit for api access.
func New(log *log.Logger, qldb *qldbdriver.QLDBDriver) Obit {
	return Obit{
		qldb: qldb,
		log: log,
	}
}

// Create adds a Obit to the QLDB. It returns the created Obit with
// fields like ID and DateCreated populated.
func (o Obit) Create(ctx context.Context, traceID string, no NewObit, now time.Time) (Info, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.obit.create")
	defer span.End()

	var i Info

	return i, nil
}


func (o Obit) FindById(ctx context.Context, traceID string, obitDID string) (Obit, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.obit.show")
	defer span.End()

	ob, err := o.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
		result, err := txn.Execute("SELECT * FROM Obits WHERE ObitDID = ?", obitDID)
		if err != nil {
			return nil, err
		}

		// Assume the result is not empty
		ionBinary, err := result.Next(txn)
		if err != nil {
			return nil, err
		}

		temp := new(Obit)
		err = ion.Unmarshal(ionBinary, temp)
		if err != nil {
			return nil, err
		}

		return *temp, nil
	})

	var obit Obit
	obit = ob.(Obit)

	if err != nil {
		return obit, err
	}

	return obit, nil
}

func (o Obit) Update(ctx context.Context, traceID string, obitId string) (*Obit, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.obit.update")
	defer span.End()

	return nil, nil
}

/**
func (o Obit) Search(ctx context.Context, traceID string) ([]Obit, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.obit.search")
	defer span.End()
}**/