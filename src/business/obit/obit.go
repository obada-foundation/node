package obit

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	"github.com/obada-foundation/sdkgo"
	"github.com/pkg/errors"
	"log"
	"math"
)

// Service provider an API to manage obits
type Service struct {
	logger *log.Logger
	sdk    *sdkgo.Sdk
	db     *sql.DB
	qldb   *qldbdriver.QLDBDriver
}

type QLDBObit struct {
	ObitDID          string            `ion:"ObitDID" db:"obit_id" json:"obit_id"`
	Usn              string            `ion:"Usn" json:"usn"`
	SerialNumberHash string            `ion:"SerialNumberHash" json:"serial_number_hash"`
	Manufacturer     string            `ion:"Manufacturer" json:"manufacturer"`
	PartNumber       string            `ion:"PartNumber" json:"part_number"`
	AlternateIDS     []string          `ion:"AlternateIDS" json:"alternate_ids"`
	OwnerDID         string            `ion:"OwnerDID" json:"owner_did"`
	ObdDID           string            `ion:"ObdDID" json:"obd_did"`
	Matadata         map[string]string `ion:"MetaData" json:"matadata"`
	StructuredData   map[string]string `ion:"StructuredData" json:"structured_data"`
	Documents        map[string]string `ion:"Documents" json:"documents"`
	ModifiedOn       int64             `ion:"ModifiedOn" json:"modified_on"`
	Status           string            `ion:"Status" json:"status"`
	RootHash         string            `ion:"RootHash" json:"root_hash"`
}

type Obits struct {
	Obits       []QLDBObit `json:"data"`
	Total       uint       `json:"total"`
	PerPage     uint       `json:"per_page"`
	CurrentPage uint       `json:"current_page"`
	LastPage    uint       `json:"last_page"`
}

// NewObitService creates new version of Obit service
func NewObitService(sdk *sdkgo.Sdk, logger *log.Logger, db *sql.DB, qldb *qldbdriver.QLDBDriver) *Service {
	return &Service{
		logger: logger,
		sdk:    sdk,
		db:     db,
		qldb:   qldb,
	}
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

	metadata, err := json.Marshal(obit.Matadata)
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

func (os Service) createQLDB(ctx context.Context, obit sdkgo.Obit) error {
	const q = "INSERT INTO Obits ?"

	_, err := os.qldb.Execute(ctx, func(txn qldbdriver.Transaction) (interface{}, error) {
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
		metadata := make(map[string]string)

		for _, record := range mdRecords.GetAll() {
			metadata[record.GetKey().GetValue()] = record.GetValue().GetValue()
		}

		o.Matadata = metadata

		strRecords := obit.GetStructuredData()
		strData := make(map[string]string)

		for _, record := range strRecords.GetAll() {
			strData[record.GetKey().GetValue()] = record.GetValue().GetValue()
		}

		o.StructuredData = strData

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
			return nil, err
		}

		o.RootHash = rootHash.GetHash()

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

		return nil, nil
	})

	return err
}

// Create method creates a new Obit
func (os Service) Create(ctx context.Context, dto sdkgo.ObitDto) error {
	obit, err := os.sdk.NewObit(dto)

	if err != nil {
		return err
	}

	if err = os.createQLDB(ctx, obit); err != nil {
		return errors.Wrap(err, "creating obit")
	}

	return nil
}

// Update method updates a new Obit
func (os Service) Update(ctx context.Context) {

}

// Delete method search Obits by given criteria ??
func (os Service) Delete(ctx context.Context) {

}

func (os Service) GetObitsCount(ctx context.Context) (uint, error) {
	var cnt uint

	const q = `SELECT COUNT(*) as cnt FROM gateway_view`

	row := os.db.QueryRow(q)
	row.Scan(&cnt)

	return cnt, nil
}

// Search method search Obits by given criteria
func (os Service) Search(ctx context.Context, offset uint) (Obits, error) {
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
}

// History the history of Obit changes
func (os Service) History(ctx context.Context) {

}
