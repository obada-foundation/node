package rest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/obada-foundation/node/business/obit"
	"github.com/obada-foundation/node/business/search"
	"github.com/obada-foundation/node/business/sys/validate"
	"github.com/obada-foundation/node/foundation/web"
	"github.com/obada-foundation/sdkgo"
	"github.com/obada-foundation/sdkgo/properties"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type node struct {
	obitService   *obit.Service
	searchService *search.Service
}

// GenerateDIDRequest host request for generateID method
type GenerateDIDRequest struct {
	SerialNumber string `validate:"required" json:"serial_number"`
	Manufacturer string `validate:"required" json:"manufacturer"`
	PartNumber   string `validate:"required" json:"part_number"`
}

type requestObit struct {
	ObitDID          interface{} `json:"obit_did"`
	Usn              string      `json:"usn"`
	SerialNumberHash string      `validate:"required" json:"serial_number_hash"`
	Manufacturer     string      `validate:"required" json:"manufacturer"`
	PartNumber       string      `validate:"required" json:"part_number"`
	OwnerDid         string      `validate:"required" json:"owner_did"`
	ObdDid           string      `json:"obd_did"`
	Metadata         []KV        `json:"metadata"`
	StructuredData   []KV        `json:"structured_data"`
	Documents        []Doc       `json:"documents"`
	ModifiedOn       int64       `json:"modified_on"`
	AlternateIDS     []string    `json:"alternate_ids"`
	Status           string      `json:"status"`
}

// Doc host obit documents
type Doc struct {
	Name     string `json:"name"`
	HashLink string `json:"hash_link"`
}

// KV host key/value pairs
type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func hashStr(str string) (string, error) {
	h := sha256.New()

	if _, err := h.Write([]byte(str)); err != nil {
		return "", fmt.Errorf("cannot wite bytes %v to hasher: %w", []byte(str), err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func (n *node) generateDID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var requestData GenerateDIDRequest

	if err := web.Decode(r, &requestData); err != nil {
		return err
	}

	if err := validate.Check(requestData); err != nil {
		return errors.Wrap(err, "validating data")
	}

	snh, err := hashStr(requestData.SerialNumber)

	if err != nil {
		return err
	}

	DID, err := n.obitService.GenerateDID(snh, requestData.Manufacturer, requestData.PartNumber)

	if err != nil {
		return err
	}

	resp := struct {
		DID string `json:"did"`
	}{
		DID: DID,
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

func (n *node) save(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	dto, err := requestBodyToDto(r)

	if err != nil {
		return err
	}

	o, err := n.obitService.Save(ctx, dto)

	if err != nil {
		return err
	}

	if er := web.Respond(ctx, w, o, http.StatusOK); er != nil {
		return er
	}

	return nil
}

func (n *node) checksum(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	dto, err := requestBodyToDto(r)

	if err != nil {
		return err
	}

	checksum, err := n.obitService.Checksum(dto)

	if err != nil {
		return err
	}

	resp := struct {
		Checksum string `json:"checksum"`
	}{
		Checksum: checksum,
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

func (n *node) search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()

	offset, err := strconv.ParseUint(query.Get("offset"), 10, 32)

	if err != nil {
		offset = 0
	}

	offsetUint := uint(offset)

	obits, err := n.searchService.Search(query.Get("q"), offsetUint)

	if err != nil {
		return err
	}

	return web.Respond(ctx, w, obits, http.StatusOK)
}

func (n *node) get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ID, err := parseDIDFromRequest(r)
	if err != nil {
		return err
	}

	o, err := n.obitService.Get(ctx, ID)

	if err != nil {
		return err
	}

	return web.Respond(ctx, w, o, http.StatusOK)
}

func requestBodyToDto(r *http.Request) (sdkgo.ObitDto, error) {
	var requestData requestObit
	var dto sdkgo.ObitDto

	if err := web.Decode(r, &requestData); err != nil {
		return dto, err
	}

	if err := validate.Check(requestData); err != nil {
		return dto, errors.Wrap(err, "validating data")
	}

	kvs := func(kvs []KV) []properties.KV {
		var newKVs []properties.KV

		for _, kv := range kvs {
			newKVs = append(newKVs, properties.KV(kv))
		}

		return newKVs
	}

	dto.SerialNumberHash = requestData.SerialNumberHash
	dto.Manufacturer = requestData.Manufacturer
	dto.PartNumber = requestData.PartNumber
	dto.OwnerDid = requestData.OwnerDid
	dto.ObdDid = requestData.ObdDid
	dto.Matadata = kvs(requestData.Metadata)
	dto.StructuredData = kvs(requestData.StructuredData)
	dto.ModifiedOn = requestData.ModifiedOn
	dto.AlternateIDS = requestData.AlternateIDS
	dto.Status = requestData.Status

	return dto, nil
}

func (n *node) history(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	DID, err := parseDIDFromRequest(r)
	if err != nil {
		return err
	}

	h, err := n.obitService.History(ctx, DID)

	if err != nil {
		return err
	}

	return web.Respond(ctx, w, h, http.StatusOK)
}

func parseDIDFromRequest(r *http.Request) (string, error) {
	DID := web.Param(r, "obitDID")

	if DID == "" {
		return "", errors.New("cannot find obitDID in URI")
	}

	return DID, nil
}
