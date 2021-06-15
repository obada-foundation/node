package handlers

import (
	"github.com/obada-foundation/node/business/web/mid"
	"github.com/obada-foundation/node/foundation/web"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger) http.Handler {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	return app
}
