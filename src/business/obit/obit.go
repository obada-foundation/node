package obit

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/amzn/ion-go/ion"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	"github.com/obada-foundation/node/business/pubsub"
	"github.com/obada-foundation/sdkgo"
	"github.com/obada-foundation/sdkgo/properties"
	"github.com/pkg/errors"
	"log"
)

// Service provider an API to manage obits
type Service struct {
	logger   *log.Logger
	sdk      *sdkgo.Sdk
	db       *sql.DB
	qldb     *qldbdriver.QLDBDriver
	pubsub    pubsub.Client
	isSynced bool
}

// NewObitService creates new version of Obit service
func NewObitService(sdk *sdkgo.Sdk, logger *log.Logger, db *sql.DB, qldb *qldbdriver.QLDBDriver, pubsub pubsub.Client) *Service {
	return &Service{
		logger:   logger,
		sdk:      sdk,
		db:       db,
		qldb:     qldb,
		pubsub:    pubsub,
		isSynced: true,
	}
}

func (os Service) updateSql(ctx context.Context, obit QLDBObit) error {
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
			root_hash = ?
		WHERE 
			obit_did = ?
	`

	stmt, err := os.db.Prepare(q)

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
		obit.RootHash,
		obit.ObitDID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (os Service) createSql(ctx context.Context, obit QLDBObit) error {
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
			 	root_hash
			) 
		    VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := os.db.Prepare(q)

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
		obit.RootHash,
	)

	if err != nil {
		return err
	}

	return nil
}

func (os Service) updateQLDB(ctx context.Context, obit sdkgo.Obit) error {
	_, err := os.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
		o, err := NewQLDBObit(obit)

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
				RootHash = ?
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
			o.RootHash,
			o.ObitDID,
		)

		if err != nil {
			return nil, err
		}

		if err := os.updateSql(ctx, o); err != nil {
			os.logger.Printf("Couldn't update obit to sql db: %v. Trying to abort QLDB transaction", obit)

			if er := txn.Abort(); er != nil {
				return nil, errors.Wrap(err, er.Error())
			}

			return nil, err
		}

		if err := os.notify(ctx, o); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func NewQLDBObit(obit sdkgo.Obit) (QLDBObit, error) {
	var o QLDBObit

	obitId := obit.GetObitID()
	o.ObitDID = obitId.GetHash().GetHash()
	o.Usn = obitId.GetUsn()
	o.SerialNumberHash = obit.GetSerialNumberHash().GetValue()
	o.Manufacturer = obit.GetManufacturer().GetValue()
	o.PartNumber = obit.GetPartNumber().GetValue()
	o.AlternateIDS = obit.GetAlternateIDS().GetValue()
	o.OwnerDID = obit.GetOwnerDID().GetValue()
	o.ObdDID = obit.GetObdDID().GetValue()

	mdRecords := obit.GetMetadata()
	strRecords := obit.GetStructuredData()

	kvs := func(records []properties.Record) []KV {
		var kvs []KV

		for _, rec := range records {
			kv := KV{
				Key: rec.GetKey().GetValue(),
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
		docs[record.GetKey().GetValue()] = record.GetValue().GetValue()
	}

	o.Documents = docs
	o.ModifiedOn = obit.GetModifiedOn().GetValue()

	o.Status = obit.GetStatus().GetValue()
	rootHash, err := obit.GetRootHash()

	if err != nil {
		return o, err
	}

	o.RootHash = rootHash.GetHash()

	return o, nil
}

func (os Service) notify(ctx context.Context, obit QLDBObit) error {
	id, err := os.pubsub.Publish(ctx, &pubsub.Msg{
		DID: obit.ObitDID,
		RootHash: obit.RootHash,
	})

	if err != nil {
		return err
	}

	os.logger.Printf("obit :: Published update to the network received corresponding id %q", id)

	return nil
}

func (os Service) findByRootHash(ctx context.Context, rootHash string) (QLDBObit, error) {
	var o QLDBObit

	_, err := os.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
		os.logger.Printf("%v", rootHash)

		const q = "SELECT * FROM Obits WHERE RootHash = ?"

		res, err := txn.Execute(q, rootHash)

		if err != nil {
			return nil, err
		}

		hasNext := res.Next(txn)
		if !hasNext && res.Err() != nil {
			return nil, res.Err()
		}

		ionBinary := res.GetCurrentData()

		os.logger.Printf("%v", ionBinary)

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

func (os Service) createQLDB(ctx context.Context, obit sdkgo.Obit) error {
	_, err := os.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
		o, err := NewQLDBObit(obit)

		if err != nil {
			return nil, err
		}

		const q = "INSERT INTO Obits ?"

		if _, err := txn.Execute(q, o); err != nil {
			return nil, err
		}

		if err := os.createSql(ctx, o); err != nil {
			os.logger.Printf("Couldn't insert obit to sql db: %v. Trying to abort QLDB transaction", obit)

			if er := txn.Abort(); er != nil {
				return nil, errors.Wrap(err, er.Error())
			}

			return nil, err
		}

		if err := os.notify(ctx, o); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (os Service) Save(ctx context.Context, dto sdkgo.ObitDto) (QLDBObit, error) {
	var o QLDBObit

	obit, err := os.sdk.NewObit(dto)

	if err != nil {
		return o, err
	}

	ID := obit.GetObitID()
	DID := ID.GetHash().GetHash()

	o, err = os.Show(ctx, DID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) == false {
			return o, err
		}

		if _, err := os.Create(ctx, dto); err != nil {
			return o, err
		}
	} else {
		if err := os.Update(ctx, DID, dto); err != nil {
			return o, err
		}
	}

	o, err = os.Show(ctx, DID)

	if err != nil {
		return o, nil
	}

	return o, nil
}

// Create method creates a new Obit
func (os Service) Create(ctx context.Context, dto sdkgo.ObitDto) (properties.ObitID, error) {
	obit, err := os.sdk.NewObit(dto)

	var ID properties.ObitID

	if err != nil {
		return ID, err
	}

	if err = os.createQLDB(ctx, obit); err != nil {
		return ID, errors.Wrap(err, "creating obit")
	}

	ID = obit.GetObitID()

	return ID, nil
}

// Update method updates Obit
func (os Service) Update(ctx context.Context, id string, dto sdkgo.ObitDto) error {
	obit, err := os.sdk.NewObit(dto)

	if err != nil {
		return err
	}

	obitID := obit.GetObitID()
	hash := obitID.GetHash().GetHash()

	if hash != id {
		return errors.New(fmt.Sprintf("Obit integrity issue. Expect to have id: %s, %s given", id, hash))
	}

	if err := os.updateQLDB(ctx, obit); err != nil {
		return err
	}

	return nil
}

// Delete method delete Obit by id (did, usn)
func (os Service) Delete(ctx context.Context, id string) error {
	// This method doesn't make sense ask Rohi for removing it. We could just update status

	return nil
}

// GetObitsCount returns total number of obits in database
func (os Service) GetObitsCount(ctx context.Context) (uint, error) {
	var cnt uint

	const q = `SELECT COUNT(*) as cnt FROM gateway_view`

	row := os.db.QueryRow(q)
	row.Scan(&cnt)

	return cnt, nil
}

func (os Service) Search(ctx context.Context) ([]QLDBObit, error) {
	var obits []QLDBObit

	const q = `SELECT * FROM gateway_view`

	stmt, err := os.db.Prepare(q)

	if err != nil {
		return obits, err
	}

	rows, err := stmt.Query()

	if err != nil {
		return obits, err
	}
	defer rows.Close()

	for rows.Next() {
		var o QLDBObit

		var altIDS []byte
		var metadata []byte
		var stctData []byte
		var docs []byte

		err := rows.Scan(
			&o.ObitDID,
			&o.Usn,
			&o.SerialNumberHash,
			&o.Manufacturer,
			&o.PartNumber,
			&altIDS,
			&o.OwnerDID,
			&o.ObdDID,
			&o.Status,
			&metadata,
			&stctData,
			&docs,
			&o.ModifiedOn,
			&o.RootHash,
		)

		json.Unmarshal(metadata, &o.Metadata)
		if err != nil {
			return obits, err
		}

		json.Unmarshal(stctData, &o.StructuredData)
		if err != nil {
			return obits, err
		}

		json.Unmarshal(docs, &o.Documents)
		if err != nil {
			return obits, err
		}

		json.Unmarshal(altIDS, &o.AlternateIDS)
		if err != nil {
			return obits, err
		}

		obits = append(obits, o)
	}

	return obits, nil
}

// Search method search Obits by given criteria
/**func (os Service) Search(ctx context.Context, offset uint) (Obits, error) {
	var obits Obits

	const perPage = 50

	obits.PerPage = perPage

	if offset == 0 {
		obits.CurrentPage = 1
	} else {
		obits.CurrentPage = offset * perPage
	}

	const q = `SELECT * FROM gateway_view LIMIT ? OFFSET ?`

	stmt, err := os.db.Prepare(q)

	if err != nil {
		return obits, err
	}

	rows, err := stmt.Query(perPage, offset)

	if err != nil {
		return obits, err
	}
	defer rows.Close()

	for rows.Next() {
		var o QLDBObit

		var altIDS []byte
		var metadata []byte
		var stctData []byte
		var docs []byte

		err := rows.Scan(
			&o.ObitDID,
			&o.Usn,
			&o.SerialNumberHash,
			&o.Manufacturer,
			&o.PartNumber,
			&altIDS,
			&o.OwnerDID,
			&o.ObdDID,
			&o.Status,
			&metadata,
			&stctData,
			&docs,
			&o.ModifiedOn,
			&o.RootHash,
		)

		json.Unmarshal(metadata, &o.Matadata)
		if err != nil {
			return obits, err
		}

		json.Unmarshal(stctData, &o.StructuredData)
		if err != nil {
			return obits, err
		}

		json.Unmarshal(docs, &o.Documents)
		if err != nil {
			return obits, err
		}

		json.Unmarshal(altIDS, &o.AlternateIDS)
		if err != nil {
			return obits, err
		}

		obits.Obits = append(obits.Obits, o)
	}

	obitsCount, err := os.GetObitsCount(ctx)
	if err != nil {
		return obits, err
	}

	obits.Total = obitsCount

	lastPage := uint(math.Ceil(float64(obitsCount) / float64(perPage)))

	if lastPage == 0 {
		lastPage = 1
	}

	obits.LastPage = lastPage

	return obits, nil
}**/

// Show returns obit by given id
func (os Service) Show(ctx context.Context, id string) (QLDBObit, error) {
	var obit QLDBObit
	var altIDS []byte
	var metadata []byte
	var stctData []byte
	var docs []byte

	const q = `
		SELECT 
			* 
		FROM 
			gateway_view 
		WHERE 
			obit_did = ? OR
			usn = ?
	`

	row := os.db.QueryRow(q, id, id)

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
		&obit.RootHash,
	)

	if err != nil {
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

	os.logger.Printf("%v", obit)

	return obit, nil
}

// History the history of Obit changes
func (os Service) History(ctx context.Context, id string) ([]QldbMeta, error) {
	var history []QldbMeta

	m, err := os.getObitQLDBMeta(ctx, id)

	if err != nil {
		return history, err
	}

	_, err = os.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
		const q = `
			SELECT 
				*
			FROM 
				history(Obits) as h
			WHERE 
				h.metadata.id = ?
		`

		result, err := txn.Execute(q, m.Metadata.ID)
		if err != nil {
			return nil, err
		}

		for result.Next(txn) {
			ionBinary := result.GetCurrentData()

			var m QldbMeta

			if err = ion.Unmarshal(ionBinary, &m); err != nil {
				return nil, err
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

func (os Service) getObitQLDBMeta(ctx context.Context, id string) (QldbMeta, error) {
	m, err := os.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
		var m QldbMeta

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

		if err = ion.Unmarshal(ionBinary, &m); err != nil {
			return m, err
		}

		return m, nil
	})

	if err != nil {
		return m.(QldbMeta), err
	}

	return m.(QldbMeta), nil
}

func (os Service) Sync(ctx context.Context) error {
	msg, err := os.pubsub.Subscribe(ctx)

	if err != nil {
		return errors.Wrap(err, "Cannot messages pull from the SQS")
	}

	if msg == nil {
		return nil
	}

	os.logger.Printf("obit :: received message from SQS %v", msg.ID)

	const q = `SELECT COUNT(*) as cnt FROM gateway_view WHERE obit_did = ? LIMIT 1`

	stmt, err := os.db.Prepare(q)

	if err != nil {
		return err
	}

	row := stmt.QueryRow(msg.DID)

	var cnt int

	if err := row.Scan(&cnt); err != nil {
		return err
	}

	o, err := os.findByRootHash(ctx, msg.RootHash)

	if err != nil {
		return err
	}

	if o.ObitDID != msg.DID {
		return errors.New("Integrity problem. IDS during sync do not match")
	}

	switch cnt {
	case 0:
		return os.createSql(ctx, o)
	case 1:
		return os.updateSql(ctx, o)
	default:
		return errors.New("Integrity problem. Broken data.")
	}
}
