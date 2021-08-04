package handlers

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

type obitGroup struct {
	obitService *obit.Service
	searchService *search.Service
}

type GenerateIDRequest struct {
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

type Doc struct {
	Name     string `json:"name"`
	HashLink string `json:"hash_link"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ObitCreate struct {
	Hash string `json:"hash"`
	DID  string `json:"did"`
	Usn  string `json:"usn"`
}

func hashStr(str string) (string, error) {
	h := sha256.New()

	if _, err := h.Write([]byte(str)); err != nil {
		return "", fmt.Errorf("cannot wite bytes %v to hasher: %w", []byte(str), err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func (og obitGroup) generateID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var requestData GenerateIDRequest

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

	ID, err := og.obitService.GenerateID(snh, requestData.Manufacturer, requestData.PartNumber)

	if err != nil {
		return err
	}

	return web.Respond(ctx, w, ID, http.StatusOK)
}

func (og obitGroup) save(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	dto, err := requestBodyToDto(ctx, r)

	if err != nil {
		return err
	}

	obit, err := og.obitService.Save(ctx, dto)

	if err != nil {
		return err
	}

	web.Respond(ctx, w, obit, http.StatusOK)

	return nil
}

func (og obitGroup) checksum(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	dto, err := requestBodyToDto(ctx, r)

	if err != nil {
		return err
	}

	checksum, err := og.obitService.Checksum(ctx, dto)

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

func (og obitGroup) search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()

	offset, err := strconv.ParseUint(query.Get("offset"), 10, 32)

	if err != nil {
		offset = 0
	}

	offsetUint := uint(offset)

	obits, err := og.searchService.Search(ctx, query.Get("q"), offsetUint)

	if err != nil {
		return err
	}

	return web.Respond(ctx, w, obits, http.StatusOK)
}

func (og obitGroup) get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ID, err := parseObitIDFromRequest(r)
	if err != nil {
		return err
	}

	obit, err := og.obitService.Get(ctx, ID)

	if err != nil {
		return err
	}

	return web.Respond(ctx, w, obit, http.StatusOK)
}

func requestBodyToDto(ctx context.Context, r *http.Request) (sdkgo.ObitDto, error) {
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
	//dto.Documents = requestData.Documents
	dto.ModifiedOn = requestData.ModifiedOn
	dto.AlternateIDS = requestData.AlternateIDS
	dto.Status = requestData.Status

	return dto, nil
}

func (og obitGroup) history(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ID, err := parseObitIDFromRequest(r)
	if err != nil {
		return err
	}

	h, err := og.obitService.History(ctx, ID)

	if err != nil {
		return err
	}

	return web.Respond(ctx, w, h, http.StatusOK)
}
