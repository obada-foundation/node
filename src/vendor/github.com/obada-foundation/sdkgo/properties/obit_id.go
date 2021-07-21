package properties

import (
	"fmt"
	"github.com/obada-foundation/sdkgo/base58"
	"github.com/obada-foundation/sdkgo/hash"
	"log"
)

// ObitID represent obit identifier
type ObitID struct {
	usn  string
	did  string
	hash hash.Hash
}

// NewObitIDProperty creates new ObitID from given arguments
func NewObitIDProperty(serialNumberHash, manufacturer, partNumber StringProperty, logger *log.Logger, debug bool) (ObitID, error) {
	var id ObitID

	snhHash := serialNumberHash.GetHash()
	mh := manufacturer.GetHash()
	pnh := partNumber.GetHash()

	sum := serialNumberHash.GetHash().GetDec() + manufacturer.GetHash().GetDec() + partNumber.GetHash().GetDec()

	if debug {
		logger.Printf(
			"\n <|%s|> => NewObitIDProperty(%v, %v, %v) -> (%d + %d + %d) -> %d",
			"Making ObitID",
			serialNumberHash,
			manufacturer,
			partNumber,
			snhHash.GetDec(),
			mh.GetDec(),
			pnh.GetDec(),
			sum,
		)
	}

	h, err := hash.NewHash(fmt.Sprintf("%x", sum), logger, debug)

	if err != nil {
		return id, fmt.Errorf("cannot create obit id: %w", err)
	}

	hashStr := h.GetHash()

	id.hash = h
	id.did = fmt.Sprintf("did:obada:%s", hashStr)
	id.usn = base58.Encode([]byte(hashStr))[:8]

	if debug {
		log.Printf("Hash: %s", h.GetHash())
		log.Printf("Did: %s", id.did)
		log.Printf("Usn: %s", id.usn)
	}

	return id, nil
}

// GetHash returns ObitID hash
func (id *ObitID) GetHash() hash.Hash {
	return id.hash
}

// GetDid returns obit DID
func (id *ObitID) GetDid() string {
	return id.did
}

// GetUsn returns the universal serial number
func (id *ObitID) GetUsn() string {
	return id.usn
}
