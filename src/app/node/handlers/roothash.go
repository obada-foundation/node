package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/obada-foundation/node/foundation/web"
	helperService "github.com/obada-foundation/node/business/helper"
)

type rootHash struct {
	logger *log.Logger
	service *helperService.Service
}

type RootHashResponse struct {
	Status int `json:"status"`
	RootHash string `json:"root_hash"`
}

// generateObit Returns the Obit Definition for a given device_id, part_number and serial_number input.
func (rh rootHash) generateRootHash(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var lo helperService.LocalObit

	if err := web.Decode(r, &lo); err != nil {
		return err
	}

	hash, err := rh.service.GenRootHash(lo)

	if err != nil {
		return err
	}

	resp := RootHashResponse{
		Status:   0,
		RootHash: hash,
	}

	web.Respond(ctx, w, resp, http.StatusOK)

	return nil
}