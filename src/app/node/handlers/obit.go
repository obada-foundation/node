package handlers

import (
	"context"
	"github.com/obada-foundation/node/business/obit"
	"net/http"
)

type obitGroup struct {
	service *obit.Service
}

func (og obitGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (og obitGroup) search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}
