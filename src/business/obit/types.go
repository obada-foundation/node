package obit

import (
	"github.com/amzn/ion-go/ion"
	"github.com/obada-foundation/node/business/types"
)

// ID used by generateID method
type ID struct {
	DID string `json:"did"`
	ID  string `json:"id"`
}

// QldbMeta hosts QLDB meta (https://docs.aws.amazon.com/qldb/latest/developerguide/working.metadata.html)
type QldbMeta struct {
	BlockAddress BlockAddress   `ion:"blockAddress" json:"block_address"`
	Hash         interface{}    `ion:"hash" json:"hash"`
	Data         types.QLDBObit `ion:"data" json:"data"`
	Metadata     Metadata       `ion:"metadata" json:"metadata"`
}

// Metadata part of QldbMeta
type Metadata struct {
	ID      string        `ion:"id" json:"id"`
	Version int           `ion:"version" json:"version"`
	TxTime  ion.Timestamp `ion:"txTime" json:"tx_time"`
	TxID    string        `ion:"txId" json:"tx_id"`
}

// BlockAddress part of QldbMeta
type BlockAddress struct {
	StrandID   interface{} `ion:"strandId" json:"strand_id"`
	SequenceNo interface{} `ion:"sequenceNo" json:"sequence_no"`
}
