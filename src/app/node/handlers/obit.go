package handlers

import (
	"context"
	"github.com/obada-foundation/node/business/obit"
	"github.com/obada-foundation/node/foundation/web"
	"github.com/obada-foundation/sdkgo"
	"github.com/pkg/errors"
	"net/http"
)

type obitGroup struct {
	service *obit.Service
}

type requestObit struct {
	SerialNumberHash string            `validate:"required" json:"serial_number_hash"`
	Manufacturer     string            `validate:"required" json:"manufacturer"`
	PartNumber       string            `validate:"required" json:"part_number"`
	OwnerDid         string            `validate:"required" json:"owner_did"`
	ObdDid           string            `json:"obd_did"`
	Metadata         map[string]string `json:"metadata"`
	StructuredData   map[string]string `json:"structured_data"`
	Documents        map[string]string `json:"documents"`
	ModifiedOn       int64             `json:"modified_on"`
	AlternateIDS     []string          `json:"alternate_ids"`
	Status           string            `json:"status"`
}

func (og obitGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	dto, err := requestBodyToDto(ctx, r)

	if err != nil {
		return err
	}

	if err := og.service.Create(ctx, dto); err != nil {
		return err
	}

	return web.Respond(ctx, w, "", http.StatusCreated)
}

func (og obitGroup) search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	obits, err := og.service.Search(ctx, 0)

	if err != nil {
		return err
	}

	return web.Respond(ctx, w, obits, http.StatusOK)
}

func parseObitIDFromRequest(r *http.Request) (string, error) {
	params := web.Params(r)

	ID, ok := params["obitDID"]

	if !ok {
		return "", errors.New("Cannot find obitDID in URI")
	}

	return ID, nil
}

func (og obitGroup) show(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ID, err := parseObitIDFromRequest(r)
	if err != nil {
		return err
	}

	obit, err := og.service.Show(ctx, ID)

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

	dto.SerialNumberHash = requestData.SerialNumberHash
	dto.Manufacturer = requestData.Manufacturer
	dto.PartNumber = requestData.PartNumber
	dto.OwnerDid = requestData.OwnerDid
	dto.ObdDid = requestData.ObdDid
	dto.Matadata = requestData.Metadata
	dto.StructuredData = requestData.StructuredData
	dto.Documents = requestData.Documents
	dto.ModifiedOn = requestData.ModifiedOn
	dto.AlternateIDS = requestData.AlternateIDS
	dto.Status = requestData.Status

	return dto, nil
}

func (og obitGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ID, err := parseObitIDFromRequest(r)
	if err != nil {
		return err
	}

	dto, err := requestBodyToDto(ctx, r)

	if err != nil {
		return err
	}

	if err := og.service.Update(ctx, ID, dto); err != nil {
		return err
	}

	return web.Respond(ctx, w, "", http.StatusOK)
}

func (og obitGroup) history(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ID, err := parseObitIDFromRequest(r)
	if err != nil {
		return err
	}

	h, err := og.service.History(ctx, ID)

	if err != nil {
		return err
	}

	return web.Respond(ctx, w, h, http.StatusOK)
}
