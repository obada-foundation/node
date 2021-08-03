package obit

import "github.com/amzn/ion-go/ion"

type ID struct {
	DID string `json:"did"`
	ID  string `json:"id"`
}

type QLDBObit struct {
	ObitDID          string            `ion:"ObitDID" db:"obit_id" json:"obit_did"`
	Usn              string            `ion:"Usn" json:"usn"`
	SerialNumberHash string            `ion:"SerialNumberHash" json:"serial_number_hash"`
	Manufacturer     string            `ion:"Manufacturer" json:"manufacturer"`
	PartNumber       string            `ion:"PartNumber" json:"part_number"`
	AlternateIDS     []string          `ion:"AlternateIDS" json:"alternate_ids"`
	OwnerDID         string            `ion:"OwnerDID" json:"owner_did"`
	ObdDID           string            `ion:"ObdDID" json:"obd_did"`
	Metadata         []KV              `ion:"MetaData" json:"metadata"`
	StructuredData   []KV              `ion:"StructuredData" json:"structured_data"`
	Documents        map[string]string `ion:"Documents" json:"documents"`
	ModifiedOn       int64             `ion:"ModifiedOn" json:"modified_on"`
	Status           string            `ion:"Status" json:"status"`
	Checksum         string            `ion:"Checksum" json:"checksum"`
}

type KV struct {
	Key   string `ion:"Key" json:"key"`
	Value string `ion:"Value" json:"value"`
}

type QldbMeta struct {
	BlockAddress BlockAddress `ion:"blockAddress" json:"block_address"`
	Hash         interface{}  `ion:"hash" json:"hash"`
	Data         QLDBObit     `ion:"data" json:"data"`
	Metadata     Metadata     `ion:"metadata" json:"metadata"`
}

type Metadata struct {
	ID      string        `ion:"id" json:"id"`
	Version int           `ion:"version" json:"version"`
	TxTime  ion.Timestamp `ion:"txTime" json:"tx_time"`
	TxID    string        `ion:"txId" json:"tx_id"`
}

type BlockAddress struct {
	StrandId   interface{} `ion:"strandId" json:"strand_id"`
	SequenceNo interface{} `ion:"sequenceNo" json:"sequence_no"`
}

type Obits struct {
	Obits       []QLDBObit `json:"data"`
	Total       uint       `json:"total"`
	PerPage     uint       `json:"per_page"`
	CurrentPage uint       `json:"current_page"`
	LastPage    uint       `json:"last_page"`
}
