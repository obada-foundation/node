package handlers

import (
	"context"
	"github.com/obada-foundation/node/business/obit"
	"github.com/obada-foundation/node/foundation/web"
	"github.com/obada-foundation/sdkgo"
	"net/http"
)

type obitGroup struct {
	service *obit.Service
}

type createObit struct {
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
	var requestData createObit

	err := web.Decode(r, &requestData)

	if err != nil {
		return err
	}

	var dto sdkgo.ObitDto
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

	if err := og.service.Create(ctx, dto); err != nil {
		return err
	}

	return nil
}

func (og obitGroup) search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	obits, err := og.service.Search(ctx, 0)

	if err != nil {
		return err
	}

	return web.Respond(ctx, w, obits, http.StatusOK)
}

func (og obitGroup) show(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (og obitGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (og obitGroup) history(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}
