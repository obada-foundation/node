package handlers

import (
	"github.com/obada-foundation/node/business/helper"
	"github.com/obada-foundation/node/business/obit"
	"github.com/obada-foundation/node/business/web/mid"
	"github.com/obada-foundation/node/foundation/web"

	"github.com/pkg/errors"

	"log"
	"net/http"
	"os"
)

// API register REST api endpoints
func API(build string, shutdown chan os.Signal, logger *log.Logger, os *obit.Service, hs *helper.Service) http.Handler {
	app := web.NewApp(shutdown, mid.Logger(logger), mid.Errors(logger), mid.Metrics(), mid.Panics(logger))

	ob := obitGroup{
		service: os,
	}

	app.Handle(http.MethodGet, "/obits", ob.search)
	app.Handle(http.MethodPost, "/obits", ob.create)
	app.Handle(http.MethodGet, "/obits/:obitDID", ob.show)
	app.Handle(http.MethodPut, "/obits/:obitDID", ob.update)
	app.Handle(http.MethodGet, "/obits/:obitDID/history", ob.history)

	od := obitDefinition{
		service: hs,
	}

	app.Handle(http.MethodGet, "/api/obit/definition", od.generateObit)

	rh := rootHash{
		logger: logger,
		service: hs,
	}

	app.Handle(http.MethodPost, "/api/obit/hash", rh.generateRootHash)

	c := client{
		logger: logger,
		helperService: hs,
		obitService: os,
	}

	app.Handle(http.MethodGet, "/api/client/obit/:obitDID", c.getClientObit)
	app.Handle(http.MethodGet, "/api/client/obits", c.getClientObits)
	app.Handle(http.MethodPost, "/api/client/obit", c.saveClientObit)
	app.Handle(http.MethodGet, "/api/server/obit/:obitDID", c.getServerObit)

	return app
}

func parseObitIDFromRequest(r *http.Request) (string, error) {
	params := web.Params(r)

	ID, ok := params["obitDID"]

	if !ok {
		return "", errors.New("Cannot find obitDID in URI")
	}

	return ID, nil
}
