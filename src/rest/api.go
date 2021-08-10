package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/obada-foundation/node/business/obit"
	"github.com/obada-foundation/node/business/search"
	"github.com/obada-foundation/node/business/web/mid"
	"github.com/obada-foundation/node/foundation/web"
)

// Rest is a rest access server
type Rest struct {
	Logger  *log.Logger
	Version string
	Address string
	Port    int
	NodeURL string

	SSLConfig SSLConfig

	ObitService   *obit.Service
	SearchService *search.Service

	httpsServer *http.Server
	httpServer  *http.Server
	lock        sync.Mutex
}

// Run runs a web server that handles an API
func (r *Rest) Run(address string, port int) {
	if address == "*" {
		address = ""
	}

	switch r.SSLConfig.SSLMode {
	case None:
		r.Logger.Printf("rest :: activate http rest server on %s:%d", address, port)
		r.lock.Lock()
		r.httpServer = r.makeHTTPServer(address, port, r.router())
		r.lock.Unlock()

		err := r.httpServer.ListenAndServe()
		r.Logger.Printf("rest :: http server terminated, %s", err)
	case Static:
		r.Logger.Printf("activate https server in 'static' mode on %s:%d", address, r.SSLConfig.Port)

		r.lock.Lock()
		r.httpsServer = r.makeHTTPSServer(address, r.SSLConfig.Port, r.router())
		r.httpServer = r.makeHTTPServer(address, port, r.httpToHTTPSRouter())
		r.lock.Unlock()

		go func() {
			r.Logger.Printf("activate http redirect server on %s:%d", address, port)
			err := r.httpServer.ListenAndServe()
			r.Logger.Printf("http redirect server terminated, %s", err)
		}()

		err := r.httpsServer.ListenAndServeTLS(r.SSLConfig.Cert, r.SSLConfig.Key)
		r.Logger.Printf("https server terminated, %s", err)
	case Auto:
		r.Logger.Printf("activate https server in 'auto' mode on %s:%d", address, r.SSLConfig.Port)

		m := r.makeAutocertManager()
		r.lock.Lock()
		r.httpsServer = r.makeHTTPSAutocertServer(address, r.SSLConfig.Port, r.router(), m)

		r.httpServer = r.makeHTTPServer(address, port, r.httpChallengeRouter(m))

		r.lock.Unlock()

		go func() {
			r.Logger.Printf("activate http challenge server on port %d", port)

			err := r.httpServer.ListenAndServe()
			r.Logger.Printf("http challenge server terminated, %s", err)
		}()

		err := r.httpsServer.ListenAndServeTLS("", "")
		r.Logger.Printf("https server terminated, %s", err)
	}
}

func (r *Rest) makeHTTPServer(address string, port int, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf("%s:%d", address, port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

// Shutdown handles shutting down of API server
func (r *Rest) Shutdown() {
	r.Logger.Print("rest :: shutdown rest server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r.lock.Lock()
	if r.httpServer != nil {
		if err := r.httpServer.Shutdown(ctx); err != nil {
			log.Printf("rest :: http shutdown error, %s", err)
		}
		r.Logger.Print("rest :: shutdown http server completed")
	}

	if r.httpsServer != nil {
		r.Logger.Print("rest :: shutdown https server")
		if err := r.httpsServer.Shutdown(ctx); err != nil {
			r.Logger.Printf("rest :: https shutdown error, %s", err)
		}
		r.Logger.Print("rest :: shutdown https server completed")
	}
	r.lock.Unlock()
}

func (r *Rest) router() http.Handler {
	app := web.NewApp(
		mid.Logger(r.Logger),
		mid.Errors(r.Logger),
		mid.Metrics(),
		mid.Panics(r.Logger),
	)

	utilGrp := util{}
	app.Handle(http.MethodGet, "/ping", utilGrp.ping)

	nodeGrp := node{
		obitService:   r.ObitService,
		searchService: r.SearchService,
	}

	app.Handle(http.MethodPost, "/obit/did", nodeGrp.generateDID)
	app.Handle(http.MethodPost, "/obit/checksum", nodeGrp.checksum)
	app.Handle(http.MethodGet, "/obits", nodeGrp.search)
	app.Handle(http.MethodPost, "/obits", nodeGrp.save)
	app.Handle(http.MethodGet, "/obits/:obitDID", nodeGrp.get)
	app.Handle(http.MethodGet, "/obits/:obitDID/history", nodeGrp.history)

	return app
}
