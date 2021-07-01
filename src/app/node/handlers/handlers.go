package handlers

import (
	"github.com/obada-foundation/node/business/obit"
	"github.com/obada-foundation/node/business/web/mid"
	"github.com/obada-foundation/node/foundation/web"
	"log"
	"net/http"
	"os"
)

// API register REST api endpoints
func API(build string, shutdown chan os.Signal, logger *log.Logger, obitService *obit.Service) http.Handler {
	app := web.NewApp(shutdown, mid.Logger(logger), mid.Errors(logger), mid.Metrics(), mid.Panics(logger))

	ob := obitGroup{
		service: obitService,
	}

	app.Handle(http.MethodGet, "/obits", ob.search)
	app.Handle(http.MethodPost, "/obits", ob.create)
	app.Handle(http.MethodGet, "/obits/:obitDID", ob.show)
	app.Handle(http.MethodPut, "/obits/:obitDID", ob.update)
	app.Handle(http.MethodGet, "/obits/:obitDID/history", ob.history)

	return app
}
