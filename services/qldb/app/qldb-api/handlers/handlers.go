// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	"log"
	"net/http"
	"os"

	"github.com/obada-protocol/server-gateway/services/qldb/business/mid"
	"github.com/obada-protocol/server-gateway/services/qldb/foundation/web"
	"github.com/obada-protocol/server-gateway/services/qldb/business/data/obit"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, qldb *qldbdriver.QLDBDriver) http.Handler {
	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	// Register health check endpoint. This route is not authenticated.
	cg := checkGroup{
		build: build,
	}
	app.Handle(http.MethodGet, "/v1/readiness", cg.readiness)
	app.Handle(http.MethodGet, "/v1/liveness", cg.liveness)

	og := obitGroup{
		obit: obit.New(log, qldb),
	}

	app.Handle(http.MethodPost, "/v1/obits", og.create)
	app.Handle(http.MethodGet, "/v1/obits/:obitDID", og.show)
	app.Handle(http.MethodDelete, "/v1/obits/:obitDID", og.delete)
	app.Handle(http.MethodGet, "/v1/obits", og.search)
	app.Handle(http.MethodPut, "/v1/obits/:obitDID", og.update)

	return app
}