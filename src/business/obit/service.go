package obit

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/amzn/ion-go/ion"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	"github.com/obada-foundation/node/business/pubsub"
	"github.com/obada-foundation/node/business/sys/validate"
	"github.com/obada-foundation/node/business/types"
	"github.com/obada-foundation/sdkgo"
	"github.com/obada-foundation/sdkgo/hash"
	"github.com/obada-foundation/sdkgo/properties"
	"github.com/pkg/errors"
	"log"
)

// Service provider an API to manage obits
type Service struct {
	logger *log.Logger
	sdk    *sdkgo.Sdk
	db     *sql.DB
	qldb   *qldbdriver.QLDBDriver
	pubsub pubsub.Client
}

// NewObitService creates new version of Obit service
func NewObitService(sdk *sdkgo.Sdk, logger *log.Logger, db *sql.DB, qldb *qldbdriver.QLDBDriver, ps pubsub.Client) *Service {
	return &Service{
		logger: logger,
		sdk:    sdk,
		db:     db,
		qldb:   qldb,
		pubsub: ps,
	}
}

func (s Service) updateSQL(obit types.QLDBObit) error {
	const q = `
		UPDATE 
		    gateway_view
		SET 
			alternate_ids = ?, 
			owner_did = ?,
			obd_did = ?,
			status = ?,
			metadata = ?,
			structured_data = ?,
			documents = ?,
			modified_on = ?,
			checksum = ?
		WHERE 
			obit_did = ?
	`

	stmt, err := s.db.Prepare(q)

	if err != nil {
		return err
	}

	altIDS, err := json.Marshal(obit.AlternateIDS)
	if err != nil {
		return err
	}

	metadata, err := json.Marshal(obit.Metadata)
	if err != nil {
		return err
	}

	stctData, err := json.Marshal(obit.StructuredData)
	if err != nil {
		return err
	}

	docs, err := json.Marshal(obit.Documents)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		altIDS,
		obit.OwnerDID,
		obit.ObdDID,
		obit.Status,
		metadata,
		stctData,
		docs,
		obit.ModifiedOn,
		obit.Checksum,
		obit.ObitDID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s Service) createSQL(obit types.QLDBObit) error {
	const q = `
		INSERT INTO 
		    gateway_view(
				obit_did, 
			 	usn, 
		 		serial_number_hash, 
		 		manufacturer, 
			 	part_number, 
			 	alternate_ids, 
			 	owner_did,
		 		obd_did,
			 	status,
		 		metadata,
		 		structured_data,
		 		documents,
				modified_on,
			 	checksum
			) 
		    VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := s.db.Prepare(q)

	if err != nil {
		return err
	}

	altIDS, err := json.Marshal(obit.AlternateIDS)
	if err != nil {
		return err
	}

	metadata, err := json.Marshal(obit.Metadata)
	if err != nil {
		return err
	}

	stctData, err := json.Marshal(obit.StructuredData)
	if err != nil {
		return err
	}

	docs, err := json.Marshal(obit.Documents)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		obit.ObitDID,
		obit.Usn,
		obit.SerialNumberHash,
		obit.Manufacturer,
		obit.PartNumber,
		altIDS,
		obit.OwnerDID,
		obit.ObdDID,
		obit.Status,
		metadata,
		stctData,
		docs,
		obit.ModifiedOn,
		obit.Checksum,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s Service) updateQLDB(ctx context.Context, obit sdkgo.Obit) error {
	_, err := s.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {

		h, _ := hash.NewHash([]byte(""), nil, false)

		o, err := NewQLDBObit(obit, &h)

		if err != nil {
			return nil, err
		}

		const q = `
			UPDATE 
				Obits 
			SET
				AlternateIDS = ?,
				OwnerDID = ?,
				ObdDID = ?,
				MetaData = ?,
				StructuredData = ?,
				Documents = ?,
				ModifiedOn = ?,
				Status = ?,
				Checksum = ?
			WHERE
				ObitDID = ?
		`

		_, err = txn.Execute(
			q,
			o.AlternateIDS,
			o.OwnerDID,
			o.ObdDID,
			o.Metadata,
			o.StructuredData,
			o.Documents,
			o.ModifiedOn,
			o.Status,
			o.Checksum,
			o.ObitDID,
		)

		if err != nil {
			return nil, err
		}

		if err := s.updateSQL(o); err != nil {
			s.logger.Printf("Couldn't update obit to sql db: %v. Trying to abort QLDB transaction", obit)

			if er := txn.Abort(); er != nil {
				return nil, errors.Wrap(err, er.Error())
			}

			return nil, err
		}

		if err := s.notify(ctx, o); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

// NewQLDBObit creates a new QLDB obit
func NewQLDBObit(obit sdkgo.Obit, parentChecksum *hash.Hash) (types.QLDBObit, error) {
	var o types.QLDBObit

	obitID := obit.GetObitID()
	o.ObitDID = obitID.GetDid()
	o.Usn = obitID.GetUsn()
	o.SerialNumberHash = obit.GetSerialNumberHash().GetValue()
	o.Manufacturer = obit.GetManufacturer().GetValue()
	o.PartNumber = obit.GetPartNumber().GetValue()
	o.AlternateIDS = obit.GetAlternateIDS().GetValue()
	o.OwnerDID = obit.GetOwnerDID().GetValue()
	o.ObdDID = obit.GetObdDID().GetValue()

	mdRecords := obit.GetMetadata()
	strRecords := obit.GetStructuredData()

	kvs := func(records []properties.Record) []types.KV {
		var kvs []types.KV

		for _, rec := range records {
			kv := types.KV{
				Key:   rec.GetKey().GetValue(),
				Value: rec.GetValue().GetValue(),
			}

			kvs = append(kvs, kv)
		}

		return kvs
	}

	o.Metadata = kvs(mdRecords.GetAll())
	o.StructuredData = kvs(strRecords.GetAll())

	docRecords := obit.GetDocuments()
	docs := make(map[string]string)

	for _, record := range docRecords.GetAll() {
		docs[record.GetName().GetValue()] = record.GetHashLink().GetHashLink()
	}

	o.Documents = docs
	o.ModifiedOn = obit.GetModifiedOn().GetValue()
	o.Status = obit.GetStatus().GetValue()

	checksum, err := obit.GetChecksum(parentChecksum)

	if err != nil {
		return o, err
	}

	o.Checksum = checksum.GetHash()

	return o, nil
}

//nolint:unused,unparam // Needs to be implemented
func (s Service) getParentObit(ctx context.Context, obitDID string) (*sdkgo.Obit, *types.QLDBObit, error) {
	_, err := s.Get(ctx, obitDID)
	if err != nil {
		return nil, nil, err
	}

	_, err = s.History(ctx, obitDID)

	if err != nil {
		return nil, nil, err
	}

	return nil, nil, err
}

func (s Service) notify(ctx context.Context, obit types.QLDBObit) error {
	id, err := s.pubsub.Publish(ctx, &pubsub.Msg{
		DID:      obit.ObitDID,
		Checksum: obit.Checksum,
	})

	if err != nil {
		return err
	}

	s.logger.Printf("obit :: Published update to the network received corresponding id %q", id)

	return nil
}

//nolint:unused // Required for future use
func (s Service) findByDID(ctx context.Context, did string) (types.QLDBObit, error) {
	var o types.QLDBObit

	_, err := s.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {

		const q = "SELECT * FROM Obits WHERE ObitDID = ?"

		res, err := txn.Execute(q, did)

		if err != nil {
			return nil, err
		}

		hasNext := res.Next(txn)
		if !hasNext && res.Err() != nil {
			return nil, res.Err()
		}

		ionBinary := res.GetCurrentData()

		err = ion.Unmarshal(ionBinary, &o)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return o, err
	}

	return o, nil
}

//nolint:unused // Needs to be implemented
func (s Service) findByChecksum(ctx context.Context, checksum string) (types.QLDBObit, error) {
	var o types.QLDBObit

	_, err := s.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {

		const q = "SELECT * FROM Obits WHERE Checksum = ?"

		res, err := txn.Execute(q, checksum)

		if err != nil {
			return nil, err
		}

		hasNext := res.Next(txn)
		if !hasNext && res.Err() != nil {
			return nil, res.Err()
		}

		ionBinary := res.GetCurrentData()

		err = ion.Unmarshal(ionBinary, &o)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("cannot unmarshal ion data given by checksum %s", checksum))
		}

		return nil, nil
	})

	if err != nil {
		return o, err
	}

	return o, nil
}

func (s Service) createQLDB(ctx context.Context, obit sdkgo.Obit) error {
	_, err := s.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
		o, err := NewQLDBObit(obit, nil)

		if err != nil {
			return nil, err
		}

		const q = "INSERT INTO Obits ?"

		if _, err := txn.Execute(q, o); err != nil {
			return nil, err
		}

		if err := s.createSQL(o); err != nil {
			s.logger.Printf("Couldn't insert obit to sql db: %v. Trying to abort QLDB transaction", obit)

			if er := txn.Abort(); er != nil {
				return nil, errors.Wrap(err, er.Error())
			}

			return nil, err
		}

		if err := s.notify(ctx, o); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

// Save create or update Obit
func (s Service) Save(ctx context.Context, dto sdkgo.ObitDto) (types.QLDBObit, error) {
	var o types.QLDBObit

	obit, err := s.sdk.NewObit(dto)

	if err != nil {
		return o, err
	}

	ID := obit.GetObitID()
	DID := ID.GetDid()

	o, err = s.Get(ctx, DID)

	if err != nil {
		// factor this!
		if err.Error() != "not found" {
			return o, err
		}

		if er := s.createQLDB(ctx, obit); er != nil {
			return o, errors.Wrap(er, "creating obit")
		}
	} else if er := s.updateQLDB(ctx, obit); er != nil {
		return o, errors.Wrap(er, "updating obit")
	}

	o, err = s.Get(ctx, DID)

	if err != nil {
		return o, nil
	}

	return o, nil
}

// Get returns obit by given DID or USN
func (s Service) Get(ctx context.Context, did string) (types.QLDBObit, error) {
	var obit types.QLDBObit
	var altIDS []byte
	var metadata []byte
	var stctData []byte
	var docs []byte

	const q = `
		SELECT 
			gv.obit_did,
		    gv.usn,
		    gv.serial_number_hash,
			gv.manufacturer,
		    gv.part_number,
		    gv.alternate_ids,
		    gv.owner_did,
		    gv.obd_did,
		    gv.status,
		   	gv.metadata,
		    gv.structured_data,
		    gv.documents,
		    gv.modified_on,
		    gv.checksum 
		FROM 
			gateway_view AS gv
		WHERE 
			gv.obit_did = ? OR
			gv.usn = ?
	`

	row := s.db.QueryRow(q, did, did)

	err := row.Scan(
		&obit.ObitDID,
		&obit.Usn,
		&obit.SerialNumberHash,
		&obit.Manufacturer,
		&obit.PartNumber,
		&altIDS,
		&obit.OwnerDID,
		&obit.ObdDID,
		&obit.Status,
		&metadata,
		&stctData,
		&docs,
		&obit.ModifiedOn,
		&obit.Checksum,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return obit, &validate.ErrorNotFoundResponse{}
		}

		return obit, err
	}

	if err := json.Unmarshal(metadata, &obit.Metadata); err != nil {
		return obit, err
	}

	if err := json.Unmarshal(stctData, &obit.StructuredData); err != nil {
		return obit, err
	}

	if err := json.Unmarshal(docs, &obit.Documents); err != nil {
		return obit, err
	}

	if err := json.Unmarshal(altIDS, &obit.AlternateIDS); err != nil {
		return obit, err
	}

	return obit, nil
}

// History the history of Obit changes
func (s Service) History(ctx context.Context, did string) ([]QLDBMeta, error) {
	var history []QLDBMeta

	obit, err := s.Get(ctx, did)

	if err != nil {
		return history, err
	}

	m, err := s.getObitQLDBMeta(ctx, obit.ObitDID)

	if err != nil {
		return history, err
	}

	_, err = s.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
		const q = `
			SELECT 
				*
			FROM 
				history(Obits) as h
			WHERE 
				h.metadata.id = ?
		`

		result, er := txn.Execute(q, m.Metadata.ID)
		if er != nil {
			return nil, er
		}

		for result.Next(txn) {
			ionBinary := result.GetCurrentData()

			var m QLDBMeta

			if er := ion.Unmarshal(ionBinary, &m); er != nil {
				return nil, er
			}

			history = append(history, m)
		}

		return nil, nil
	})

	if err != nil {
		return history, err
	}

	return history, nil
}

func (s Service) getObitQLDBMeta(ctx context.Context, id string) (QLDBMeta, error) {
	m, err := s.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
		var m QLDBMeta

		const q = `
			SELECT 
				*
			FROM 
				_ql_committed_Obits as o 
			WHERE 
				o.data.Usn = ? OR 
				o.data.ObitDID = ?
		`

		result, err := txn.Execute(q, id, id)
		if err != nil {
			return m, err
		}

		// Assume the result is not empty
		hasNext := result.Next(txn)
		if !hasNext && result.Err() != nil {
			return m, result.Err()
		}

		ionBinary := result.GetCurrentData()

		if err := ion.Unmarshal(ionBinary, &m); err != nil {
			return m, err
		}

		return m, nil
	})

	if err != nil {
		return m.(QLDBMeta), err
	}

	return m.(QLDBMeta), nil
}

// Sync should be deprecated once we agree about real distributed protocol
func (s Service) Sync(ctx context.Context) error {
	msg, err := s.pubsub.Subscribe(ctx)

	if err != nil {
		return errors.Wrap(err, "Cannot messages pull from the SQS")
	}

	if msg == nil {
		return nil
	}

	s.logger.Printf("obit :: received message from SQS %v", msg.ID)

	const q = `SELECT COUNT(*) as cnt FROM gateway_view WHERE obit_did = ? LIMIT 1`

	stmt, err := s.db.Prepare(q)

	if err != nil {
		return err
	}

	row := stmt.QueryRow(msg.DID)

	var cnt int

	if er := row.Scan(&cnt); er != nil {
		return er
	}

	s.logger.Printf("incoming msg %+v", msg)

	o, err := s.findByDID(ctx, msg.DID)

	if err != nil {
		return err
	}

	if o.ObitDID != msg.DID {
		return errors.New("integrity problem, IDS during sync are not match")
	}

	switch cnt {
	case 0:
		return s.createSQL(o)
	case 1:
		return s.updateSQL(o)
	default:
		return errors.New("integrity problem, broken data")
	}
}

// GenerateDID generates obit DID
func (s Service) GenerateDID(serialNumberHash, manufacturer, partNumber string) (string, error) {
	dto := sdkgo.ObitIDDto{
		SerialNumberHash: serialNumberHash,
		Manufacturer:     manufacturer,
		PartNumber:       partNumber,
	}

	sdkID, err := s.sdk.NewObitID(dto)

	if err != nil {
		return "", err
	}

	return sdkID.GetDid(), nil
}

// Checksum generates obit checksum
func (s Service) Checksum(dto sdkgo.ObitDto) (string, error) {
	o, err := s.sdk.NewObit(dto)

	if err != nil {
		return "", err
	}

	h, err := o.GetChecksum(nil)

	if err != nil {
		return "", err
	}

	return h.GetHash(), nil
}
