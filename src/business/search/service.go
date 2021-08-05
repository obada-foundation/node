package search

import (
	"database/sql"
	"encoding/json"
	"github.com/obada-foundation/node/business/types"
	"log"
	"math"
	"regexp"
)

const perPage = 50

// Obits represents the collection of obits with pagination
type Obits struct {
	Obits []types.QLDBObit `json:"data"`
	Meta  ObitsMeta        `json:"meta"`
}

// ObitsMeta contain pagination option, in future might have other options as well
type ObitsMeta struct {
	Total       uint `json:"total"`
	PerPage     uint `json:"per_page"`
	CurrentPage uint `json:"current_page"`
	LastPage    uint `json:"last_page"`
}

// Service provider an API to manage obits
type Service struct {
	logger *log.Logger
	db     *sql.DB
}

// NewService creates new version of Obit service
func NewService(logger *log.Logger, db *sql.DB) *Service {
	return &Service{
		logger: logger,
		db:     db,
	}
}

// GetAll returns all paginated obits from node without filtering
func (s Service) GetAll(offset uint) (*sql.Rows, error) {
	const q = `
		SELECT 
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
		FROM 
			gateway_view
		ORDER BY 
			modified_on
		LIMIT ? OFFSET ?
	`

	stmt, err := s.db.Prepare(q)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(perPage, offset)

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s Service) searchFullText(term string, offset uint) (*sql.Rows, error) {
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
		    gateway_view_fts AS gvf
		JOIN 
			gateway_view as gv ON gvf.rowid = gv.id
		WHERE 
		      gateway_view_fts MATCH ?
		ORDER BY 
			rank
		LIMIT ? OFFSET ?
	`

	stmt, err := s.db.Prepare(q)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(term, perPage, offset)

	if err != nil {
		return nil, err
	}

	return rows, nil
}

// Search looking for obits by given term
func (s Service) Search(term string, offset uint) (Obits, error) {
	var obits Obits

	var rows *sql.Rows

	regexp, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	newTerm := regexp.ReplaceAllString(term, "")

	if len(newTerm) > 0 {
		ftsRows, err := s.searchFullText(newTerm, offset)

		if err != nil {
			return obits, err
		}

		rows = ftsRows
	} else {
		genRows, err := s.GetAll(offset)

		if err != nil {
			return obits, err
		}

		rows = genRows
	}

	defer rows.Close()

	for rows.Next() {
		var o types.QLDBObit

		var altIDS []byte
		var metadata []byte
		var stctData []byte
		var docs []byte

		er := rows.Scan(
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
			&o.Checksum,
		)

		if er != nil {
			return obits, er
		}

		if er := json.Unmarshal(metadata, &o.Metadata); er != nil {
			return obits, er
		}

		if er := json.Unmarshal(stctData, &o.StructuredData); er != nil {
			return obits, er
		}

		if er := json.Unmarshal(docs, &o.Documents); er != nil {
			return obits, er
		}

		if er := json.Unmarshal(altIDS, &o.AlternateIDS); er != nil {
			return obits, er
		}

		obits.Obits = append(obits.Obits, o)
	}

	obitsCount, err := s.GetObitsCountByTerm(newTerm)
	if err != nil {
		return obits, err
	}

	obits.Meta.Total = obitsCount

	lastPage := uint(math.Ceil(float64(obitsCount) / float64(perPage)))

	if lastPage == 0 {
		lastPage = 1
	}

	obits.Meta.CurrentPage = offset + 1
	obits.Meta.LastPage = lastPage

	return obits, nil
}

// GetObitsCountByTerm returns total number of obits in database by given term
func (s Service) GetObitsCountByTerm(term string) (uint, error) {
	var cnt uint
	var row *sql.Row

	if term == "" {
		const q = `
			SELECT 
				COUNT(*) AS cnt
			FROM 
				gateway_view
		`

		row = s.db.QueryRow(q, term)
	} else {
		const q = `
			SELECT 
				COUNT(*) AS cnt
			FROM 
				gateway_view_fts AS gvf
			JOIN 
				gateway_view as gv ON gvf.rowid = gv.id
			WHERE 
				  gateway_view_fts MATCH ?
		`

		row = s.db.QueryRow(q, term)
	}

	if err := row.Scan(&cnt); err != nil {
		return cnt, err
	}

	return cnt, nil
}
