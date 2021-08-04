package handlers

import (
	"context"
	"github.com/obada-foundation/node/business/types"
	"log"
	"net/http"

	helperService "github.com/obada-foundation/node/business/helper"
	obitService "github.com/obada-foundation/node/business/obit"
	"github.com/obada-foundation/node/foundation/web"
)

type client struct {
	logger        *log.Logger
	helperService *helperService.Service
	obitService   *obitService.Service
}

type SaveClientObitResponse struct {
	Status int                  `json:"status"`
	Obit   types.QLDBObit `json:"obit"`
}

type GetClientObitResponse struct {
	Status int                  `json:"status"`
	Obit   types.QLDBObit `json:"obit"`
}

type GetClientObitsResponse struct {
	Status int                    `json:"status"`
	Obits  []types.QLDBObit `json:"obits"`
}

type GetServerObitResponse struct {
	Status         int            `json:"status"`
	BlockChainObit BlockChainObit `json:"blockchain_obit"`
}

type BlockChainObit struct {
	RootHash string `json:"root_hash"`
	ObitDID  string `json:"obit_did"`
}

func (c client) getServerObit(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	DID, err := parseObitIDFromRequest(r)

	if err != nil {
		return err
	}

	o, err := c.obitService.Get(ctx, DID[10:])

	if err != nil {
		return err
	}

	resp := GetServerObitResponse{
		Status: 0,
	}

	resp.BlockChainObit.RootHash = o.Checksum
	resp.BlockChainObit.ObitDID = o.ObitDID

	web.Respond(ctx, w, resp, http.StatusOK)

	return nil
}

func (c client) saveClientObit(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var lo helperService.LocalObit

	if err := web.Decode(r, &lo); err != nil {
		return err
	}

	dto, err := c.helperService.ToDto(lo)

	if err != nil {
		return err
	}

	obit, err := c.obitService.Save(ctx, dto)

	if err != nil {
		return err
	}

	resp := SaveClientObitResponse{
		Status: 0,
		Obit:   obit,
	}

	web.Respond(ctx, w, resp, http.StatusOK)

	return nil
}
