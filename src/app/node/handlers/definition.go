package handlers

import (
	"github.com/obada-foundation/node/foundation/web"
	helperService "github.com/obada-foundation/node/business/helper"

	"context"
	"net/http"
)

type obitDefinition struct {
	service *helperService.Service
}

type ObitDefRequest struct {
	SerialNumber string
	Manufacturer string
	PartNumber string
}

type GenerateObitResponse struct {
	Status int `json:"status"`
	Obit helperService.ObitDef `json:"obit"`
}

// generateObit Returns the Obit Definition for a given device_id, part_number and serial_number input.
func (og obitDefinition) generateObit(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()

	defRequest := ObitDefRequest{
		SerialNumber: q.Get("serial_number"),
		Manufacturer: q.Get("manufacturer"),
		PartNumber: q.Get("part_number"),
	}

	ob, err := og.service.GenObitDef(defRequest.SerialNumber, defRequest.Manufacturer, defRequest.PartNumber)

	if err != nil {
		return err
	}

	resp := GenerateObitResponse{
		Status: 0,
		Obit: ob,
	}

	web.Respond(ctx, w, resp, http.StatusOK)

	return nil
}
