package mid

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/obada-foundation/node/foundation/web"
)

// Logger writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Logger(logger *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(before web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			logger.Printf("%s : started : %s %s -> %s",
				"zz",
				r.Method, r.URL.Path, r.RemoteAddr,
			)

			err := before(ctx, w, r)

			logger.Printf("%s : completed : %s %s -> %s (%d) (%s)",
				"zz",
				r.Method, r.URL.Path, r.RemoteAddr,
				v.StatusCode, time.Since(v.Now),
			)

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return m
}
