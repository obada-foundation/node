package properties

import (
	"fmt"
	"github.com/obada-foundation/sdkgo/hash"
	"log"
	"strconv"
)

// KvProperty representing slice of records and their hash
type KvProperty struct {
	records []Record
	hash    hash.Hash
}

// Record of key/value
type Record struct {
	key   StringProperty
	value StringProperty
	hash  hash.Hash
}

type KV struct {
	Key string
	Value string
}

// NewRecord creates a new key/value record
func NewRecord(description, key, value string, logger *log.Logger, debug bool) (Record, error) {
	var r Record

	if debug {
		logger.Printf("\n |%s| => NewRecord(%q, %q)", description, key, value)
	}

	k, err := NewStringProperty("New record key", key, logger, debug)

	if err != nil {
		return r, err
	}

	v, err := NewStringProperty("New record value", value, logger, debug)

	if err != nil {
		return r, err
	}

	kh := k.GetHash()
	vh := k.GetHash()
	kvDec := kh.GetDec() + vh.GetDec()

	if debug {
		log.Printf("(%d + %d) -> %d", kh.GetDec(), vh.GetDec(), kvDec)
	}

	h, err := hash.NewHash(strconv.FormatUint(kvDec, 10), logger, debug)

	if err != nil {
		return r, err
	}

	r.key = k
	r.value = v
	r.hash = h

	return r, nil
}

// GetKey returns a record key
func (r *Record) GetKey() StringProperty {
	return r.key
}

// GetValue returns a record value
func (r *Record) GetValue() StringProperty {
	return r.value
}

// GetHash returns a record hash
func (r *Record) GetHash() hash.Hash {
	return r.hash
}

// NewMapProperty creates map property
func NewMapProperty(description string, kvs []KV, logger *log.Logger, debug bool) (KvProperty, error) {
	var mp KvProperty
	var kvDec uint64

	if debug {
		logger.Printf("\n <|%s|> => NewMapProperty(%v)", description, kvs)
	}

	for _, kv := range kvs {
		description = "\t" + description + " :: creating key/value record"

		r, err := NewRecord(description, kv.Key, kv.Value, logger, debug)

		if err != nil {
			return mp, err
		}

		rh := r.GetHash()
		kvDec += rh.GetDec()

		mp.records = append(mp.records, r)
	}

	h, err := hash.NewHash(strconv.FormatUint(kvDec, 10), logger, debug)

	if err != nil {
		return mp, fmt.Errorf("cannot hash %q: %w", kvDec, err)
	}

	mp.hash = h

	return mp, nil
}

// GetAll returns slice of records
func (mp *KvProperty) GetAll() []Record {
	return mp.records
}

// GetHash returns hash of all records
func (mp *KvProperty) GetHash() hash.Hash {
	return mp.hash
}
