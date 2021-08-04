package obit

import (
	"github.com/amzn/ion-go/ion"
	"github.com/obada-foundation/node/business/types"
)

type ID struct {
	DID string `json:"did"`
	ID  string `json:"id"`
}

type QldbMeta struct {
	BlockAddress BlockAddress `ion:"blockAddress" json:"block_address"`
	Hash         interface{}  `ion:"hash" json:"hash"`
	Data         types.QLDBObit     `ion:"data" json:"data"`
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
