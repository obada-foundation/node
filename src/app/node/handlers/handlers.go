package handlers

import (
	"context"
	"github.com/obada-foundation/node/business/obit"
	"github.com/obada-foundation/node/business/search"
	"github.com/obada-foundation/node/business/web/mid"
	"github.com/obada-foundation/node/foundation/web"

	"github.com/pkg/errors"

	"log"
	"net/http"
	"os"
)

// Options represent optional parameters.
type Options struct {
	corsOrigin string
}

// APIConfig represents API server dependencies
type APIConfig struct {
	Shutdown      chan os.Signal
	Logger        *log.Logger
	ObitService   *obit.Service
	SearchService *search.Service
}

// API register REST api endpoints
func API(cfg APIConfig, options ...func(opts *Options)) http.Handler {
	var opts Options
	for _, option := range options {
		option(&opts)
	}

	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Logger),
		mid.Errors(cfg.Logger),
		mid.Metrics(),
		mid.Panics(cfg.Logger),
	)

	ob := obitGroup{
		obitService:   cfg.ObitService,
		searchService: cfg.SearchService,
	}

	app.Handle(http.MethodPost, "/obit/id", ob.generateID)
	app.Handle(http.MethodPost, "/obit/checksum", ob.checksum)
	app.Handle(http.MethodGet, "/obits", ob.search)
	app.Handle(http.MethodPost, "/obits", ob.save)
	app.Handle(http.MethodGet, "/obits/:obitDID", ob.get)
	app.Handle(http.MethodGet, "/obits/:obitDID/history", ob.history)

	// Accept CORS 'OPTIONS' preflight requests if config has been provided.
	// Don't forget to apply the CORS middleware to the routes that need it.
	// Example Config: `conf:"default:https://MY_DOMAIN.COM"`
	if opts.corsOrigin != "" {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return nil
		}
		app.Handle(http.MethodOptions, "/*", h)
	}

	return app
}

func parseObitIDFromRequest(r *http.Request) (string, error) {
	ID := web.Param(r, "obitDID")

	if ID == "" {
		return "", errors.New("cannot find obitDID in URI")
	}

	return ID, nil
}
