package sdkgo

import (
	"github.com/obada-foundation/sdkgo/properties"
	"log"
)

// Status is representing the Obit status
type Status string

// StatusFunctional represent a functional obit
const StatusFunctional Status = "FUNCTIONAL"

// StatusNonFunctional represent a non-functional obit
const StatusNonFunctional Status = "NON_FUNCTIONAL"

// StatusDisposed represent a disposed obit
const StatusDisposed Status = "DISPOSED"

// StatusStolen represent a stolen obit
const StatusStolen Status = "STOLEN"

// DisabledByOwner represent obit which was disabled by owner
const DisabledByOwner Status = "DISABLED_BY_OWNER"

// ObitIDDto todo add description
type ObitIDDto struct {
	SerialNumberHash string `validate:"required"`
	Manufacturer     string `validate:"required"`
	PartNumber       string `validate:"required"`
}

// ObitDto todo add description
type ObitDto struct {
	ObitIDDto
	OwnerDid       string `validate:"required"`
	ObdDid         string
	Matadata       []properties.KV
	StructuredData []properties.KV
	Documents      []properties.Doc
	ModifiedOn     int64 `validate:"min-modified-on"`
	AlternateIDS   []string
	Status         string
}

// Obit represent asset data structure
type Obit struct {
	obitID           properties.ObitID
	serialNumberHash properties.StringProperty
	manufacturer     properties.StringProperty
	partNumber       properties.StringProperty
	ownerDid         properties.StringProperty
	obdDid           properties.StringProperty
	metadata         properties.KvCollection
	structuredData   properties.KvCollection
	documents        properties.Documents
	modifiedOn       properties.IntProperty
	alternateIDS     properties.SliceStrProperty
	status           properties.StatusProperty
	debug            bool
	logger           *log.Logger
}
