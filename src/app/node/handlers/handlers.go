package handlers

import (
	"github.com/obada-foundation/node/business/obit"
	"github.com/obada-foundation/node/business/web/mid"
	"github.com/obada-foundation/node/foundation/web"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger, obitService *obit.Service) http.Handler {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	ob := obitGroup{}

	app.Handle(http.MethodGet, "/obits", ob.search)
	app.Handle(http.MethodPost, "/obits", ob.create)

	return app
}
