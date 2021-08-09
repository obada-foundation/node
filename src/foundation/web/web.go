package web

import (
	"context"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux/v5"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values are stored/retrieved.
const KeyValues ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// Handler s type that describe Http handlers
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App web application
type App struct {
	mux      *httptreemux.ContextMux
	mw       []Middleware
}

// NewApp creates a new application
func NewApp(mw ...Middleware) *App {
	app := App{
		mux:      httptreemux.NewContextMux(),
		mw:       mw,
	}

	return &app
}

// Handle handles executed route in router
func (a *App) Handle(method, path string, handler Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(mw, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {

		// Start or expand a distributed trace.
		ctx := r.Context()

		// Set the context with the required values to
		// process the request.
		v := Values{
			TraceID: "zzz",
			Now:     time.Now(),
		}
		ctx = context.WithValue(ctx, KeyValues, &v)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			return
		}
	}

	// Add this handler for the specified verb and route.
	a.mux.Handle(method, path, h)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
