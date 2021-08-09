package rest

import (
	"context"
	"fmt"
	"github.com/obada-foundation/node/business/obit"
	"github.com/obada-foundation/node/business/search"
	"github.com/obada-foundation/node/business/web/mid"
	"github.com/obada-foundation/node/foundation/web"
	"log"
	"net/http"
	"sync"
	"time"
)

// Rest is a rest access server
type Rest struct {
	Logger  *log.Logger
	Version string
	Address string
	Port int

	SSLConfig   SSLConfig

	ObitService   *obit.Service
	SearchService *search.Service

	httpsServer *http.Server
	httpServer  *http.Server
	lock        sync.Mutex

	restUtil util
	restNode node
}

func (r *Rest) Run(address string, port int) {
	if address == "*" {
		address = ""
	}

	switch r.SSLConfig.SSLMode {
	case None:
		r.Logger.Printf("rest :: activate http rest server on %s:%d", address, port)
		r.lock.Lock()
		r.httpServer = r.makeHTTPServer(fmt.Sprintf("%s:%d", address, port), r.router())
		r.lock.Unlock()

		err := r.httpServer.ListenAndServe()
		r.Logger.Printf("rest :: http server terminated, %s", err)
	}
}

func (r *Rest) makeHTTPServer(address string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf("%s", address),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout: 30 * time.Second,
	}
}

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
