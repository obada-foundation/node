// Package obit contains obit related CRUD functionality.
package obit

import (
	"context"
	"github.com/amzn/ion-go/ion"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	"github.com/obada-protocol/server-gateway/services/qldb/foundation/database"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/trace"
	"log"
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
type Service struct {
	qldb *qldbdriver.QLDBDriver
	log *log.Logger
}

// New constructs a Obit for api access.
func New(log *log.Logger, qldb *qldbdriver.QLDBDriver) Service {
	return Service{
		qldb: qldb,
		log: log,
	}
}

// Create adds a Obit to the QLDB. It returns the created Obit with
// fields like ID and DateCreated populated.
func (s Service) Create(ctx context.Context, traceID string, no NewObit) (Obit, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.obit.create")
	defer span.End()

	const q = "INSERT INTO Obits ?"

	s.log.Printf("%s : %s : query : %s", traceID, "obit.Create",
		database.Log(q, no),
	)

	_, err := s.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
		return txn.Execute(q, no)
	})

	var o Obit

	if err != nil {
		return o, errors.Wrap(err, "creating obit")
	}

	return o, nil
}


func (s Service) FindById(ctx context.Context, traceID string, obitDID string) (Obit, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.obit.show")
	defer span.End()

	const q = "SELECT * FROM Obits WHERE ObitDID = ?"

	s.log.Printf("%s : %s : query : %s", traceID, "obit.FindById",
		database.Log(q, obitDID),
	)

	ob, err := s.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
		result, err := txn.Execute(q, obitDID)
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

	if err != nil {
		if err.Error() == "no more values" {
			return obit, ErrNotFound
		}

		return obit, errors.Wrap(err, "selecting single obit")
	}

	obit = ob.(Obit)

	return obit, nil
}

func (s Service) FindBy(ctx context.Context, traceID string) ([]Obit, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.obit.search")
	defer span.End()

	const q = "SELECT * FROM Obits"

	s.log.Printf("%s : %s : query : %s", traceID, "obit.FindById",
		database.Log(q),
	)

	return nil, nil
}

func (s Service) Update(ctx context.Context, traceID string, obitId string) (*Obit, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.obit.update")
	defer span.End()

	return nil, nil
}