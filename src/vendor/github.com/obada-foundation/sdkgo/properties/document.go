package properties

import (
	"fmt"
	"github.com/obada-foundation/sdkgo/hash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Doc is mapping struct to Obit documents
type Doc struct {
	Name     string
	HashLink string
}

// HashLinkProperty represent a hash of the hash link source
type HashLinkProperty struct {
	hashLink string
	hash     hash.Hash
}

// NewHashLinkProperty creates a new property and compute a hash of file under hash link
func NewHashLinkProperty(description, hashLink string, logger *log.Logger, debug bool) (HashLinkProperty, error) {
	var hlp HashLinkProperty

	if debug {
		logger.Printf("\n <|%s|> => NewStringProperty(%v)", description, hashLink)
	}

	u, err := url.Parse(hashLink)
	if err != nil {
		return hlp, err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return hlp, err
	}
	defer resp.Body.Close()

	fileBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return hlp, err
	}

	h, err := hash.NewHash(fileBytes, logger, debug)

	if err != nil {
		return hlp, err
	}

	hlp.hashLink = hashLink
	hlp.hash = h

	return hlp, nil
}

// GetHash returns hash of the file under hash link
func (hlp HashLinkProperty) GetHash() hash.Hash {
	return hlp.hash
}

// GetHashLink returns a document hash link
func (hlp HashLinkProperty) GetHashLink() string {
	return hlp.hashLink
}

// Document Obit document
type Document struct {
	name     StringProperty
	hashLink HashLinkProperty
	hash     hash.Hash
}

// NewDocument creates a new Obit document
func NewDocument(name, hashLink string, logger *log.Logger, debug bool) (Document, error) {
	var d Document

	if debug {
		logger.Printf("\n |New document| => NewDocument(%q, %q)", name, hashLink)
	}

	n, err := NewStringProperty("New document name", name, logger, debug)

	if err != nil {
		return d, err
	}

	hl, err := NewHashLinkProperty("New document hash link", hashLink, logger, debug)

	if err != nil {
		return d, err
	}

	nh := n.GetHash()
	hlh := hl.GetHash()
	docDec := nh.GetDec() + hlh.GetDec()

	if debug {
		logger.Printf("(%d + %d) -> %d", nh.GetDec(), hlh.GetDec(), docDec)
	}

	h, err := hash.NewHash([]byte(strconv.FormatUint(docDec, 10)), logger, debug)

	if err != nil {
		return d, err
	}

	d.name = n
	d.hashLink = hl
	d.hash = h

	return d, nil
}

// GetName returns a document name
func (d *Document) GetName() StringProperty {
	return d.name
}

// GetHashLink returns a document hash link
func (d *Document) GetHashLink() HashLinkProperty {
	return d.hashLink
}

// GetHash returns a document hash
func (d *Document) GetHash() hash.Hash {
	return d.hash
}

// Documents slice of documents
type Documents struct {
	documents []Document
	hash      hash.Hash
}

// NewDocumentsCollection creates the collection of documents
func NewDocumentsCollection(docs []Doc, logger *log.Logger, debug bool) (Documents, error) {
	var ds Documents
	var docDec uint64
	description := "Making documents hash"

	if debug {
		logger.Printf("\n <|%s|> => NewDocumentCollection(%v)", description, docs)
	}

	for _, doc := range docs {
		description = "\t" + description + " :: creating document"

		d, err := NewDocument(doc.Name, doc.HashLink, logger, debug)

		if err != nil {
			return ds, err
		}

		dh := d.GetHash()
		docDec += dh.GetDec()

		ds.documents = append(ds.documents, d)
	}

	h, err := hash.NewHash([]byte(strconv.FormatUint(docDec, 10)), logger, debug)

	if err != nil {
		return ds, fmt.Errorf("cannot hash %q: %w", docDec, err)
	}

	ds.hash = h

	return ds, nil
}

// GetHash returns a hash of documents collection
func (ds Documents) GetHash() hash.Hash {
	return ds.hash
}

// GetAll returns all Obit documents
func (ds Documents) GetAll() []Document {
	return ds.documents
}
