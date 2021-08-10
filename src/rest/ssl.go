package rest

import (
	"crypto/tls"
	"github.com/dimfeld/httptreemux/v5"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"net/url"
)

// sslMode defines ssl mode for rest server
type sslMode int8

const (
	// None defines to run http server only
	None sslMode = iota

	// Static defines to run both https and http server. Redirect http to https
	Static

	// Auto defines to run both https and http server. Redirect http to https. Https server with autocert support
	Auto
)

// SSLConfig holds all ssl params for rest server
type SSLConfig struct {
	SSLMode      sslMode
	Cert         string
	Key          string
	Port         int
	ACMELocation string
	ACMEEmail    string
}

// getNodeHost returns hostname for OBADA node server.
// For example for nodeURL https://tradeloop.node.obada.io:443 it should return tradeloop.node.obada.io
func (r *Rest) getNodeHost() string {
	u, err := url.Parse(r.NodeURL)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

// httpChallengeRouter creates new router which performs ACME "http-01" challenge response
// with default middlewares. This part is necessary to obtain certificate from LE.
// If it receives not a acme challenge it performs redirect to https server.
// Used in 'auto' ssl mode.
func (r *Rest) httpChallengeRouter(m *autocert.Manager) *httptreemux.ContextMux {
	r.Logger.Printf("create http-challenge routes")
	router := httptreemux.NewContextMux()

	router.Handler(http.MethodGet, "/*", m.HTTPHandler(r.redirectHandler()))
	return router
}

func (r *Rest) redirectHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		newURL := r.NodeURL + req.URL.Path
		if req.URL.RawQuery != "" {
			newURL += "?" + req.URL.RawQuery
		}
		http.Redirect(w, req, newURL, http.StatusTemporaryRedirect)
	})
}

// makeHTTPSAutoCertServer makes https server with autocert mode (LE support)
func (r *Rest) makeHTTPSAutocertServer(address string, port int, router http.Handler, m *autocert.Manager) *http.Server {
	server := r.makeHTTPServer(address, port, router)
	cfg := r.makeTLSConfig()
	cfg.GetCertificate = m.GetCertificate
	server.TLSConfig = cfg
	return server
}

func (r *Rest) makeAutocertManager() *autocert.Manager {
	return &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(r.SSLConfig.ACMELocation),
		HostPolicy: autocert.HostWhitelist(r.getNodeHost()),
		Email:      r.SSLConfig.ACMEEmail,
	}
}

// makeHTTPSServer makes https server for static mode
func (r *Rest) makeHTTPSServer(address string, port int, router http.Handler) *http.Server {
	server := r.makeHTTPServer(address, port, router)
	server.TLSConfig = r.makeTLSConfig()
	return server
}

func (r *Rest) makeTLSConfig() *tls.Config {
	return &tls.Config{
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			// tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		},
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
			tls.CurveP384,
		},
	}
}

// httpToHTTPSRouter creates new router which does redirect from http to https server
// with default middlewares. Used in 'static' ssl mode.
func (r *Rest) httpToHTTPSRouter() *httptreemux.ContextMux {
	r.Logger.Printf("create https-to-http redirect routes")

	router := httptreemux.NewContextMux()

	router.Handler(http.MethodGet, "/*", r.redirectHandler())
	return router
}
