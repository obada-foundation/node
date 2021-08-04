package search

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/obada-foundation/node/business/types"
	"log"
	"math"
	"strings"
)

const perPage = 50

type Obits struct {
	Obits []types.QLDBObit `json:"data"`
	Meta  ObitsMeta  `json:"meta"`
}

type ObitsMeta struct {
	Total       uint `json:"total"`
	PerPage     uint `json:"per_page"`
	CurrentPage uint `json:"current_page"`
	LastPage    uint `json:"last_page"`
}

// Service provider an API to manage obits
type Service struct {
	logger   *log.Logger
	db       *sql.DB
	isSynced bool
}

// NewService creates new version of Obit service
func NewService(logger *log.Logger, db *sql.DB) *Service {
	return &Service{
		logger:   logger,
		db:       db,
	}
}

// Search looking for obits by given term
func (s Service) Search(ctx context.Context, term string, offset uint) (Obits, error) {
	var obits Obits

	term = strings.ReplaceAll(term, ":", "")

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
		return obits, err
	}

	rows, err := stmt.Query(term, perPage, offset)

	if err != nil {
		return obits, err
	}
	defer rows.Close()

	for rows.Next() {
		var o types.QLDBObit

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
			&o.Checksum,
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

		obits.Obits = append(obits.Obits, o)
	}

	obitsCount, err := s.GetObitsCountByTerm(ctx, term)
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
func (s Service) GetObitsCountByTerm(ctx context.Context, term string) (uint, error) {
	var cnt uint

	term = strings.ReplaceAll(term, ":", "")

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

	row := s.db.QueryRow(q, term)
	row.Scan(&cnt)

	return cnt, nil
}
