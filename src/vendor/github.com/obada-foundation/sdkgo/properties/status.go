package properties

import (
	"fmt"
	"github.com/obada-foundation/sdkgo/hash"
	"log"
)

// StatusProperty represent a status and status hash
type StatusProperty struct {
	value string
	hash  hash.Hash
}

// NewStatusProperty creates a new status property for obit
func NewStatusProperty(status string, logger *log.Logger, debug bool) (StatusProperty, error) {
	var sp StatusProperty

	if debug {
		logger.Printf("\n <|Making status property|> => NewStatusProperty(%v)", status)
	}

	h, err := hash.NewHash([]byte(status), logger, debug)

	if err != nil {
		return sp, fmt.Errorf("cannot hash %q: %w", status, err)
	}

	sp.hash = h
	sp.value = status

	return sp, nil
}

// GetValue returns obit status
func (sp StatusProperty) GetValue() string {
	return sp.value
}

// GetHash returns status hash
func (sp StatusProperty) GetHash() hash.Hash {
	return sp.hash
}
